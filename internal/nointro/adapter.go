// Package nointro acquires one reviewed public Standard DAT through Datomatic's
// anonymous form. It does not discover platforms, log in, or execute JavaScript.
package nointro

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/MoonlarkStudios/romd-dat-catalogs/internal/definitions"
	"github.com/MoonlarkStudios/romd-dat-catalogs/internal/publisher"
)

const formLimit = 1 << 20
const requestGap = 5 * time.Second

// Datomatic serves this notice while an official export is being generated.
// It is an upstream availability failure, not evidence of a changed form.
func exportPending(raw []byte) bool {
	for _, f := range forms.FindAllStringSubmatch(string(raw), -1) {
		a, err := attrs(f[1])
		if err == nil && a["name"] == "main_form" &&
			strings.Contains(f[2], "The requested file is temporarily not available.") &&
			strings.Contains(f[2], "it's added to the queue.") {
			return true
		}
	}
	return false
}

var number = regexp.MustCompile(`^[1-9][0-9]{0,5}$`)

// ResponseDiagnostic contains only bounded structural metadata, never cookies, URLs or response text.
type ResponseDiagnostic struct {
	Attempt    int    `json:"attempt"`
	Stage      string `json:"stage"`
	Status     int    `json:"status"`
	Bytes      int    `json:"bytes"`
	BodySHA256 string `json:"bodySha256"`
	PageKind   string `json:"pageKind,omitempty"`
	Outcome    string `json:"outcome,omitempty"`
}

type Result struct {
	Diagnostics []ResponseDiagnostic
	Attempt     publisher.Attempt
	SHA256      string
	Counts      publisher.Counts
	Code        string
}

type Adapter struct {
	origin string
	client *http.Client
	gate   chan struct{}
	next   time.Time
	gap    time.Duration
}

func SourceURL(system string) string {
	return "https://datomatic.no-intro.org/index.php?page=download&op=dat&s=" + system
}
func New() *Adapter {
	return newAdapter("https://datomatic.no-intro.org", http.DefaultTransport, requestGap)
}
func newAdapter(origin string, transport http.RoundTripper, gap time.Duration) *Adapter {
	return &Adapter{origin: origin, gap: gap, gate: make(chan struct{}, 1), client: &http.Client{Transport: transport, Timeout: 30 * time.Second,
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error { return http.ErrUseLastResponse }}}
}

// Acquire makes at most three attempts within one two-minute deadline. Only
// transport failures and unexpected preparation responses are retried. Queued
// exports, provider cooldowns, form drift and invalid documents fail immediately.
func (a *Adapter) Acquire(ctx context.Context, id string, c definitions.Catalog, stage string) (Result, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()
	if err := os.Mkdir(stage, 0700); err != nil {
		return Result{}, err
	}
	var diagnostics []ResponseDiagnostic
	for attempt := 1; ; attempt++ {
		result, err := a.acquire(ctx, id, c, stage)
		for i := range result.Diagnostics {
			result.Diagnostics[i].Attempt = attempt
		}
		diagnostics = append(diagnostics, result.Diagnostics...)
		result.Diagnostics = diagnostics
		if err != nil || attempt == 3 || ctx.Err() != nil ||
			(result.Code != "prepare_failed" && result.Code != "request_failed" && result.Code != "incomplete_response") {
			return result, err
		}
		// Backoff supplements the shared provider admission gap; cancellation stops it.
		timer := time.NewTimer(a.gap * time.Duration(1<<(attempt-1)))
		select {
		case <-timer.C:
		case <-ctx.Done():
			timer.Stop()
			return result, ctx.Err()
		}
	}
}

func (a *Adapter) acquire(ctx context.Context, id string, c definitions.Catalog, stage string) (Result, error) {
	r := Result{Attempt: publisher.Attempt{CatalogID: id, ExpectedName: c.ExpectedName, SourceURL: SourceURL(c.ProviderSystemID)}}
	if a.client == nil || a.gate == nil || c.Provider != "no-intro" || c.Representation != "standard" || !number.MatchString(c.ProviderSystemID) || id != "no-intro/"+c.SystemID+"/standard" || c.ExpectedName == "" || c.Validation.MinimumGames < 1 || c.Validation.MinimumROMs < 1 {
		return r, errors.New("invalid No-Intro adapter or reviewed catalog")
	}
	select {
	case a.gate <- struct{}{}:
		defer func() { <-a.gate }()
	case <-ctx.Done():
		return r, ctx.Err()
	}
	// Anonymous session cookies live only for this acquisition and are never logged.
	jar, _ := cookiejar.New(nil)
	client := *a.client
	client.Jar = jar
	fail := func(code string) (Result, error) {
		r.Code = code
		r.Attempt.Failure = &r.Code
		if len(r.Diagnostics) > 0 {
			r.Diagnostics[len(r.Diagnostics)-1].Outcome = code
		}
		return r, nil
	}
	if time.Until(a.next) > a.gap {
		stamp := a.next.UTC().Format(time.RFC3339Nano)
		r.Attempt.RetryAt = &stamp
		return fail("provider_backoff")
	}
	source := a.origin + "/index.php?page=download&op=dat&s=" + c.ProviderSystemID
	request := func(step, target string, values url.Values, limit int) ([]byte, *http.Response, string) {
		if delay := time.Until(a.next); delay > 0 {
			timer := time.NewTimer(delay)
			defer timer.Stop()
			select {
			case <-timer.C:
			case <-ctx.Done():
				return nil, nil, "cancelled"
			}
		}
		method := http.MethodGet
		var body io.Reader
		if values != nil {
			method = http.MethodPost
			body = strings.NewReader(values.Encode())
		}
		req, err := http.NewRequestWithContext(ctx, method, target, body)
		if err != nil {
			return nil, nil, "request_failed"
		}
		req.Header.Set("User-Agent", "ROMD-DAT-Catalogs/experimental (+https://github.com/MoonlarkStudios/romd-dat-catalogs)")
		req.Header.Set("Accept-Encoding", "identity")
		if values != nil {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
		a.next = time.Now().Add(a.gap)
		diagnostic := ResponseDiagnostic{Stage: step}
		defer func() { r.Diagnostics = append(r.Diagnostics, diagnostic) }()
		resp, err := client.Do(req)
		if err != nil {
			return nil, nil, "request_failed"
		}
		defer resp.Body.Close()
		diagnostic.Status = resp.StatusCode
		if resp.StatusCode == 429 || resp.StatusCode == 503 {
			a.next = retryAt(resp.Header.Get("Retry-After"), time.Now())
			if a.next.Before(time.Now().Add(a.gap)) {
				a.next = time.Now().Add(a.gap)
			}
			stamp := a.next.UTC().Format(time.RFC3339Nano)
			r.Attempt.RetryAt = &stamp
			return nil, resp, "provider_backoff"
		}
		if resp.StatusCode != 200 && resp.StatusCode != 302 {
			return nil, resp, "http_status"
		}
		if encoding := resp.Header.Get("Content-Encoding"); encoding != "" && encoding != "identity" {
			return nil, resp, "unexpected_encoding"
		}
		if resp.ContentLength > int64(limit) {
			return nil, resp, "size_limit"
		}
		b, err := io.ReadAll(io.LimitReader(resp.Body, int64(limit)+1))
		diagnostic.Bytes = len(b)
		diagnostic.BodySHA256 = fmt.Sprintf("%x", sha256.Sum256(b))
		if step != "download" {
			diagnostic.PageKind = "unrecognized"
			if exportPending(b) {
				diagnostic.PageKind = "queued_export"
			} else if _, e := prepareForm(b, c); e == nil {
				diagnostic.PageKind = "selection_form"
			} else if _, e := downloadForm(b); e == nil {
				diagnostic.PageKind = "download_form"
			}
		}
		if err != nil {
			return nil, resp, "incomplete_response"
		}
		if len(b) > limit {
			return nil, resp, "size_limit"
		}
		return b, resp, ""
	}
	form, resp, code := request("selection", source, nil, formLimit)
	if code != "" {
		return fail(code)
	}
	if resp.StatusCode != 200 {
		return fail("unexpected_redirect")
	}
	if exportPending(form) {
		return fail("upstream_pending")
	}
	values, err := prepareForm(form, c)
	if err != nil {
		return fail("form_changed")
	}
	form, resp, code = request("prepare", source, values, formLimit)
	if code != "" {
		return fail(code)
	}
	if resp.StatusCode == 200 && exportPending(form) {
		return fail("upstream_pending")
	}
	if resp.StatusCode != 302 {
		return fail("prepare_failed")
	}
	manager, err := managerURL(a.origin, resp.Header.Get("Location"), c.ProviderSystemID)
	if err != nil {
		return fail("unexpected_redirect")
	}
	form, resp, code = request("manager", manager, nil, formLimit)
	if code != "" {
		return fail(code)
	}
	if resp.StatusCode != 200 {
		return fail("unexpected_redirect")
	}
	if exportPending(form) {
		return fail("upstream_pending")
	}
	values, err = downloadForm(form)
	if err != nil {
		return fail("form_changed")
	}
	raw, resp, code := request("download", manager, values, publisher.MaxInput)
	if code != "" {
		return fail(code)
	}
	if resp.StatusCode != 200 {
		return fail("unexpected_redirect")
	}
	document, counts, err := Validate(raw, c)
	if err != nil {
		return fail("invalid_document")
	}
	r.Counts = counts
	r.SHA256 = publisher.Hash(document)
	path, err := filepath.Abs(filepath.Join(stage, "catalog.input"))
	if err != nil {
		return r, err
	}
	if err = os.WriteFile(path, raw, 0600); err != nil {
		return r, err
	}
	r.Attempt.Path = &path
	return r, nil
}

func managerURL(origin, location, system string) (string, error) {
	base, err := url.Parse(origin)
	if err != nil {
		return "", err
	}
	u, err := url.Parse(location)
	if err != nil {
		return "", err
	}
	u = base.ResolveReference(u)
	q, err := url.ParseQuery(u.RawQuery)
	if err != nil || u.Scheme != base.Scheme || u.Host != base.Host || u.User != nil || u.Fragment != "" || u.Path != "/index.php" || len(q) != 3 || len(q["page"]) != 1 || q.Get("page") != "manager" || len(q["s"]) != 1 || q.Get("s") != system || len(q["download"]) != 1 || !number.MatchString(q.Get("download")) {
		return "", errors.New("unexpected manager URL")
	}
	return u.String(), nil
}
func retryAt(value string, now time.Time) time.Time {
	if value != "" && strings.Trim(value, "0123456789") == "" {
		seconds, err := strconv.ParseUint(value, 10, 64)
		ceiling := time.Date(9999, 1, 1, 0, 0, 0, 0, time.UTC)
		if err != nil || seconds > uint64(ceiling.Unix()-now.Unix()) {
			return ceiling
		}
		return time.Unix(now.Unix()+int64(seconds), int64(now.Nanosecond()))
	}
	if date, err := http.ParseTime(value); err == nil && date.After(now) {
		return date
	}
	return now.Add(time.Minute)
}
