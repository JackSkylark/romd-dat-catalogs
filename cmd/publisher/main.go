package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/MoonlarkStudios/romd-dat-catalogs/internal/definitions"
	"github.com/MoonlarkStudios/romd-dat-catalogs/internal/nointro"
	"github.com/MoonlarkStudios/romd-dat-catalogs/internal/publisher"
	"github.com/MoonlarkStudios/romd-dat-catalogs/internal/redump"
	"io"
	"os"
	"path/filepath"
	"time"
)

type acquirer interface {
	Acquire(context.Context, []redump.Catalog, string) ([]redump.Result, error)
}

func run(args []string, out, errOut io.Writer) error {
	return runWithAcquirer(args, out, errOut, redump.New())
}

type noIntroAcquirer interface {
	Acquire(context.Context, string, definitions.Catalog, string) (nointro.Result, error)
}

func runWithAcquirer(args []string, out, errOut io.Writer, adapter acquirer) error {
	return runWithAdapters(args, out, errOut, adapter, nointro.New())
}

// MaxSelectedCatalogs bounds combined active and paused selections across providers.
const MaxSelectedCatalogs = 32

// catalogFlags preserves operator order without allowing implicit discovery.
type catalogFlags []string

func (v *catalogFlags) String() string { return fmt.Sprint([]string(*v)) }
func (v *catalogFlags) Set(value string) error {
	if value == "" {
		return errors.New("empty catalog selection")
	}
	*v = append(*v, value)
	return nil
}

func runWithAdapters(args []string, out, errOut io.Writer, adapter acquirer, noIntro noIntroAcquirer) error {
	f := flag.NewFlagSet("publisher", flag.ContinueOnError)
	f.SetOutput(errOut)
	output := f.String("output", "", "publication directory (required)")
	paused := f.Bool("paused", false, "pause all selected catalogs")
	registryPath := f.String("definitions", "definitions", "reviewed definitions directory")
	var active, pauses catalogFlags
	f.Var(&active, "catalog", "exact catalog ID to acquire (repeatable); exclusive with manifest")
	f.Var(&pauses, "pause-catalog", "exact catalog ID to pause (repeatable)")
	base := f.String("base-url", "https://catalogs.example.invalid/", "public HTTPS base URL")
	if err := f.Parse(args); err != nil {
		return err
	}
	selected := append(append(catalogFlags{}, active...), pauses...)
	if *output == "" || (len(selected) == 0 && (f.NArg() != 1 || *paused)) || (len(selected) > 0 && f.NArg() != 0) {
		return errors.New("usage: publisher --output DIR MANIFEST | --catalog ID [--catalog ID ...] [--pause-catalog ID] [--definitions DIR] [--paused]")
	}
	if len(selected) > MaxSelectedCatalogs {
		return fmt.Errorf("at most %d catalog selections are allowed", MaxSelectedCatalogs)
	}
	if len(selected) == 0 {
		attempts, err := publisher.ReadManifest(f.Arg(0))
		if err != nil {
			return err
		}
		snapshot, err := publisher.Publish(*output, attempts, *base, publisher.Options{})
		if err != nil {
			return err
		}
		return writeSummary(out, snapshot)
	}
	registry, err := definitions.LoadSource(*registryPath)
	if err != nil {
		return err
	}
	oldDefinitions, err := definitions.Load(filepath.Join(*output, ".definitions.json"))
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	if err := registry.Compatible(oldDefinitions); err != nil {
		return err
	}
	seen := map[string]bool{}
	for _, id := range selected {
		if _, ok := registry.Catalogs[id]; !ok {
			return errors.New("unknown catalog selection")
		}
		if seen[id] {
			return errors.New("duplicate catalog selection")
		}
		seen[id] = true
	}
	prior, err := publisher.LoadSnapshot(*output)
	if err != nil {
		if _, stat := os.Stat(filepath.Join(*output, "current.json")); !errors.Is(stat, os.ErrNotExist) {
			return err
		}
		prior = publisher.Snapshot{}
	}
	// Provider cooldown survives new processes, even if the throttled catalog
	// is not selected in this invocation. Validate all deadlines before fetching.
	retryByProvider := map[string]time.Time{}
	for id, c := range prior.Catalogs {
		definition, known := registry.Catalogs[id]
		if !known {
			continue
		}
		if c.Name != definition.ExpectedName {
			return errors.New("catalog identity changed; explicit migration required")
		}
		if c.RetryAt == nil {
			continue
		}
		retry, err := time.Parse(time.RFC3339Nano, *c.RetryAt)
		if err != nil {
			return errors.New("invalid published retry time")
		}
		if retry.After(retryByProvider[definition.Provider]) {
			retryByProvider[definition.Provider] = retry
		}
	}
	stageRoot, err := os.MkdirTemp("", "romd-acquisition-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(stageRoot)
	var attempts []publisher.Attempt
	for i, id := range selected {
		c := registry.Catalogs[id]
		old, exists := prior.Catalogs[id]
		if *paused || i >= len(active) {
			if exists {
				failure := "publication_paused"
				attempts = append(attempts, publisher.Attempt{CatalogID: id, ExpectedName: c.ExpectedName, SourceURL: sourceURL(c), Failure: &failure, RetryAt: old.RetryAt})
			}
			continue
		}
		if retry := retryByProvider[c.Provider]; time.Now().Before(retry) {
			stamp := retry.UTC().Format(time.RFC3339Nano)
			if exists && old.RetryAt != nil && *old.RetryAt == stamp {
				continue
			}
			failure := "provider_backoff"
			attempts = append(attempts, publisher.Attempt{CatalogID: id, ExpectedName: c.ExpectedName, SourceURL: sourceURL(c), Failure: &failure, RetryAt: &stamp})
			continue
		}
		stage := filepath.Join(stageRoot, fmt.Sprint(i))
		var attempt publisher.Attempt
		switch c.Provider {
		case "no-intro":
			result, err := noIntro.Acquire(context.Background(), id, c, stage)
			if err != nil {
				return err
			}
			attempt = result.Attempt
			for _, diagnostic := range result.Diagnostics {
				encoded, err := json.Marshal(diagnostic)
				if err != nil {
					return err
				}
				fmt.Fprintf(errOut, "catalog %s response: %s\n", id, encoded)
			}
		case "redump":
			results, err := adapter.Acquire(context.Background(), []redump.Catalog{{ID: id, System: c.ProviderSystemID, ExpectedName: c.ExpectedName, Platform: c.SystemID, Representation: c.Representation, PolicyVersion: "1", MinGames: c.Validation.MinimumGames, MinROMs: c.Validation.MinimumROMs}}, stage)
			if err != nil {
				return err
			}
			if len(results) != 1 {
				return errors.New("unexpected acquisition result count")
			}
			attempt = results[0].Attempt
		default:
			return errors.New("unsupported provider")
		}
		if attempt.Failure != nil {
			fmt.Fprintf(errOut, "catalog %s: %s\n", id, *attempt.Failure)
		}
		if attempt.RetryAt != nil {
			retry, err := time.Parse(time.RFC3339Nano, *attempt.RetryAt)
			if err != nil {
				return errors.New("invalid acquisition retry time")
			}
			if retry.After(retryByProvider[c.Provider]) {
				retryByProvider[c.Provider] = retry
			}
		}
		attempts = append(attempts, attempt)
	}
	if len(attempts) == 0 {
		return writeSummary(out, prior)
	}
	snapshot, err := publisher.Publish(*output, attempts, *base, publisher.Options{})
	if err != nil {
		return err
	}
	return writeSummary(out, snapshot)
}

func writeSummary(out io.Writer, s publisher.Snapshot) error {
	failed := 0
	for _, c := range s.Catalogs {
		if c.Health == "failed" {
			failed++
		}
	}
	return json.NewEncoder(out).Encode(map[string]any{"sequence": s.Sequence, "catalogs": len(s.Catalogs), "failed": failed, "trust": publisher.Trust, "details": s.Catalogs})
}
func main() {
	if e := run(os.Args[1:], os.Stdout, os.Stderr); e != nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(1)
	}
}

func sourceURL(c definitions.Catalog) string {
	if c.Provider == "no-intro" {
		return nointro.SourceURL(c.ProviderSystemID)
	}
	return redump.SourceURL(c.ProviderSystemID)
}
