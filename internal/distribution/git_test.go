package distribution

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/MoonlarkStudios/romd-dat-catalogs/internal/publisher"
)

func git(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	cmd.Env = append(os.Environ(), "GIT_AUTHOR_NAME=Test", "GIT_AUTHOR_EMAIL=test@example.invalid", "GIT_COMMITTER_NAME=Test", "GIT_COMMITTER_EMAIL=test@example.invalid")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v %s", args, err, out)
	}
	return strings.TrimSpace(string(out))
}

func dataClient(dir string) *http.Client {
	return &http.Client{Transport: transportFunc(func(r *http.Request) (*http.Response, error) {
		if !gitAssetURL.MatchString(r.URL.String()) {
			return nil, fmt.Errorf("unexpected URL %s", r.URL)
		}
		parts := strings.SplitN(strings.TrimPrefix(r.URL.Path, "/"), "/", 4)
		b, err := exec.Command("git", "-C", dir, "show", parts[2]+":"+parts[3]).Output()
		if err != nil {
			return nil, err
		}
		return &http.Response{StatusCode: 200, Header: make(http.Header), Body: io.NopCloser(bytes.NewReader(b)), Request: r}, nil
	})}
}

func TestGitPublicationLifecycle(t *testing.T) {
	for _, id := range []string{"synthetic/console/standard", "no-intro/snes/standard"} {
		t.Run(id, func(t *testing.T) { testGitPublicationLifecycle(t, id) })
	}
}
func testGitPublicationLifecycle(t *testing.T, id string) {
	f := makeFixture(t)
	if id != "synthetic/console/standard" {
		f.state = filepath.Join(f.dir, "nointro-state")
		path := "../../fixtures/example.dat"
		if _, err := publisher.Publish(f.state, []publisher.Attempt{{CatalogID: id, ExpectedName: "ROMD Synthetic Console", SourceURL: "https://datomatic.no-intro.org/index.php?page=download&op=dat&s=49", Path: &path}}, "https://example.invalid/", publisher.Options{Now: f.now}); err != nil {
			t.Fatal(err)
		}
	}
	repo := filepath.Join(f.dir, "data")
	if err := os.Mkdir(repo, 0755); err != nil {
		t.Fatal(err)
	}
	git(t, repo, "init", "-b", "main")
	export := func() {
		t.Helper()
		if err := ExportData(f.state, repo); err != nil {
			t.Fatal(err)
		}
		git(t, repo, "add", "--", "*.dat")
	}
	export()
	git(t, repo, "commit", "-m", "initial")
	firstCommit := git(t, repo, "rev-parse", "HEAD")
	stage := func(version int64) (Index, string) {
		t.Helper()
		out := filepath.Join(f.dir, fmt.Sprintf("git-stage-%d", version))
		idx, err := StageGit(StageOptions{State: f.state, TrustDir: f.trust, Output: out, Keys: f.online, Version: version, Now: f.now.Add(time.Duration(version) * time.Second)}, "example/romd-dat-data", git(t, repo, "rev-parse", "HEAD"))
		if err != nil {
			t.Fatal(err)
		}
		return idx, out
	}
	// Migrate through the existing trust root and cache, not a bootstrap/reset.
	_, legacy := f.stage(t, 1)
	f.location.Store(legacy)
	cache := filepath.Join(f.dir, "cache")
	if _, err := f.refresh(cache); err != nil {
		t.Fatal(err)
	}
	idx, out := stage(2)
	if _, err := os.Stat(filepath.Join(out, "assets")); !os.IsNotExist(err) {
		t.Fatal("Git stage emitted DAT assets/archive")
	}
	if idx.State.URL != "" || len(idx.Downloads) != 1 {
		t.Fatal("unexpected recovery/history assets")
	}
	f.location.Store(out)
	idx, err := f.refresh(cache)
	if err != nil {
		t.Fatal(err)
	}
	if idx.Format != GitFormat {
		t.Fatal(idx.Format)
	}
	serialized, err := json.Marshal(idx)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(serialized, []byte(`"state":`)) || bytes.Contains(serialized, []byte("/releases/download/")) {
		t.Fatal("legacy storage leaked into Git index")
	}
	client := dataClient(repo)
	if err := VerifyDownloads(idx, client); err != nil {
		t.Fatal(err)
	}
	candidate, original, err := ReadCandidate(idx, id, "ROMD Synthetic Console", client)
	if err != nil {
		t.Fatal(err)
	}
	tampered := &http.Client{Transport: transportFunc(func(r *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: 200, Header: make(http.Header), Body: io.NopCloser(strings.NewReader("tampered")), Request: r}, nil
	})}
	if err := VerifyDownloads(idx, tampered); err == nil {
		t.Fatal("tampered public data accepted before deployment")
	}
	if _, _, err := ReadCandidate(idx, candidate.CatalogID, candidate.Name, tampered); err == nil {
		t.Fatal("tampered candidate accepted")
	}
	broken := idx
	broken.Downloads = map[string]Asset{}
	if err := RestorePublication(broken, filepath.Join(f.dir, "missing-binding"), client); err == nil {
		t.Fatal("missing binding accepted")
	}
	restored := filepath.Join(f.dir, "restored")
	if err := RestorePublication(idx, restored, client); err != nil {
		t.Fatal(err)
	}
	f.state = restored
	restoredSnapshot, err := publisher.LoadSnapshot(restored)
	if err != nil {
		t.Fatal(err)
	}
	oldCatalog := restoredSnapshot.Catalogs[candidate.CatalogID]
	oldEvent := restoredSnapshot.Events[0].ID
	// Different ZIP packaging preserves the exact document, Git tree, and RSS GUID.
	var zipped bytes.Buffer
	zw := zip.NewWriter(&zipped)
	w, err := zw.Create("repacked.dat")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := w.Write(original); err != nil {
		t.Fatal(err)
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(f.dir, "repacked.zip")
	if err := os.WriteFile(path, zipped.Bytes(), 0600); err != nil {
		t.Fatal(err)
	}
	attempt := publisher.Attempt{CatalogID: candidate.CatalogID, ExpectedName: candidate.Name, SourceURL: "http://redump.org/datfile/psx/", Path: &path}
	if _, err := publisher.Publish(f.state, []publisher.Attempt{attempt}, "https://catalogs.example.invalid/", publisher.Options{Now: f.now.Add(3 * time.Second)}); err != nil {
		t.Fatal(err)
	}
	export()
	if git(t, repo, "diff", "--cached", "--name-only") != "" || git(t, repo, "rev-parse", "HEAD") != firstCommit {
		t.Fatal("repack created Git changes")
	}
	unchanged, _ := stage(3)
	if len(unchanged.Snapshot.Events) != 1 || unchanged.Snapshot.Events[0].ID != oldEvent {
		t.Fatal("unchanged RSS event changed")
	}
	// A failed check retains the document, timestamps and commit; restore still works.
	failure := "upstream unavailable"
	attempt.Path = nil
	attempt.Failure = &failure
	if _, err := publisher.Publish(f.state, []publisher.Attempt{attempt}, "https://catalogs.example.invalid/", publisher.Options{Now: f.now.Add(4 * time.Second)}); err != nil {
		t.Fatal(err)
	}
	export()
	if git(t, repo, "diff", "--cached", "--name-only") != "" {
		t.Fatal("failure changed DAT")
	}
	failed, _ := stage(4)
	current := failed.Snapshot.Catalogs[candidate.CatalogID]
	if current.Health != "failed" || *current.Artifact != *oldCatalog.Artifact || *current.LastChanged != *oldCatalog.LastChanged {
		t.Fatal("failure lost last working document")
	}
	if err := RestorePublication(failed, filepath.Join(f.dir, "failed-restored"), client); err != nil {
		t.Fatal(err)
	}
	if _, raw, err := ReadCandidate(failed, candidate.CatalogID, candidate.Name, client); err != nil || !bytes.Equal(raw, original) {
		t.Fatal("retained verified catalog unavailable", err)
	}
	// A real byte change creates a new Git commit and one latest RSS event.
	changed := bytes.Replace(original, []byte("</datafile>"), []byte("<!-- document update -->\n</datafile>"), 1)
	if bytes.Equal(changed, original) {
		t.Fatal("fixture did not change")
	}
	if err := os.WriteFile(path, changed, 0600); err != nil {
		t.Fatal(err)
	}
	attempt.Path = &path
	attempt.Failure = nil
	if _, err := publisher.Publish(f.state, []publisher.Attempt{attempt}, "https://catalogs.example.invalid/", publisher.Options{Now: f.now.Add(5 * time.Second)}); err != nil {
		t.Fatal(err)
	}
	export()
	if git(t, repo, "diff", "--cached", "--name-only") != id+".dat" {
		t.Fatal("wrong stable path")
	}
	git(t, repo, "commit", "-m", "update")
	updated, _ := stage(5)
	if len(updated.Snapshot.Events) != 1 || updated.Snapshot.Events[0].ID == oldEvent || len(updated.Downloads) != 1 {
		t.Fatal("not latest-only")
	}
	if _, raw, err := ReadCandidate(updated, candidate.CatalogID, candidate.Name, client); err != nil || !bytes.Equal(raw, changed) {
		t.Fatalf("changed candidate: %v", err)
	}
	// Git history remains accessible by the old signed URL after the branch moves.
	if _, raw, err := ReadCandidate(idx, candidate.CatalogID, candidate.Name, client); err != nil || !bytes.Equal(raw, original) {
		t.Fatalf("old commit broken: %v", err)
	}
	if err := RestorePublication(updated, filepath.Join(f.dir, "updated-restored"), client); err != nil {
		t.Fatal(err)
	}
	if _, err := StageGit(StageOptions{State: f.state, TrustDir: f.trust, Output: filepath.Join(f.dir, "regressed"), Keys: f.online, Version: 2, Now: f.now}, "example/data", firstCommit); err == nil {
		t.Fatal("version rollback accepted")
	}
}

func TestGitURLValidation(t *testing.T) {
	sha := strings.Repeat("a", 40)
	good := "https://raw.githubusercontent.com/example/data/" + sha + "/redump/psx/discs.dat"
	for _, bad := range []string{strings.Replace(good, sha, "main", 1), strings.Replace(good, sha, sha[:7], 1), good + "?x=y", good + "#fragment", strings.Replace(good, "https:", "http:", 1), strings.Replace(good, "/redump/", "/../", 1)} {
		if _, err := FetchAsset(Asset{URL: bad, Bytes: 1, SHA256: "unused"}, nil); err == nil {
			t.Fatalf("accepted %s", bad)
		}
	}
	if !gitAssetURL.MatchString(good) {
		t.Fatal("full commit URL rejected")
	}
}
