package nointro

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestPreparationRetryRecoversAndDiagnosticsExcludeResponseSecrets(t *testing.T) {
	prepares := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Query().Get("page") == "manager" && r.Method == "GET":
			io.WriteString(w, manager)
		case r.URL.Query().Get("page") == "manager":
			w.Write(document())
		case r.Method == "GET":
			io.WriteString(w, selection())
		default:
			prepares++
			if prepares == 1 {
				w.Header().Set("Set-Cookie", "secret-cookie")
				io.WriteString(w, "private-response-token")
				return
			}
			w.Header().Set("Location", "index.php?page=manager&s=49&download=7")
			w.WriteHeader(302)
		}
	}))
	defer server.Close()
	a := newAdapter(server.URL, server.Client().Transport, 0)
	result, err := a.Acquire(context.Background(), "no-intro/snes/standard", catalog(), filepath.Join(t.TempDir(), "stage"))
	if err != nil || result.Attempt.Path == nil || prepares != 2 {
		t.Fatal(result, err, prepares)
	}
	if len(result.Diagnostics) != 6 || result.Diagnostics[1].Outcome != "prepare_failed" || result.Diagnostics[1].Status != 200 || result.Diagnostics[2].Attempt != 2 {
		t.Fatal(result.Diagnostics)
	}
	raw, _ := json.Marshal(result.Diagnostics)
	if strings.Contains(string(raw), "private-response-token") || strings.Contains(string(raw), "secret-cookie") || strings.Contains(string(raw), server.URL) {
		t.Fatal("diagnostics leaked response data")
	}
}

func TestPreparationRetriesAreBounded(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if r.Method == "GET" {
			io.WriteString(w, selection())
		} else {
			io.WriteString(w, "unexpected page")
		}
	}))
	defer server.Close()
	result, err := newAdapter(server.URL, server.Client().Transport, 0).Acquire(context.Background(), "no-intro/snes/standard", catalog(), filepath.Join(t.TempDir(), "stage"))
	if err != nil || result.Code != "prepare_failed" || requests != 6 || len(result.Diagnostics) != 6 {
		t.Fatal(result, err, requests)
	}
}

func TestN64NodumpControlIsOptionalButIncludedWhenOffered(t *testing.T) {
	c := catalog()
	c.SystemID = "n64"
	c.ProviderSystemID = "24"
	c.ExpectedName = "Nintendo - Nintendo 64 (BigEndian)"
	raw := strings.ReplaceAll(strings.ReplaceAll(selection(), "49", "24"), catalog().ExpectedName, "Nintendo - Nintendo 64")
	for _, name := range []string{"collection", "inc_adult", "inc_nodump"} {
		raw = regexp.MustCompile(`<input[^>]*name="`+name+`"[^>]*>`).ReplaceAllString(raw, "")
	}
	for _, present := range []bool{false, true} {
		t.Run(fmt.Sprint(present), func(t *testing.T) {
			form := raw
			if present {
				form = strings.Replace(form, "</form>", `<input type="radio" name="inc_nodump" value="1"><input type="radio" name="inc_nodump" value="0"></form>`, 1)
			}
			values, err := prepareForm([]byte(form), c)
			if err != nil || values.Has("inc_nodump") != present || (present && values.Get("inc_nodump") != "1") {
				t.Fatal(values, err)
			}
			// Optional does not mean accepting an offered exclusion-only control.
			broken := strings.Replace(raw, "</form>", `<input type="radio" name="inc_nodump" value="0"></form>`, 1)
			if _, err := prepareForm([]byte(broken), c); err == nil {
				t.Fatal("exclusion-only control accepted")
			}
		})
	}
}
