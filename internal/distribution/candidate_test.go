package distribution

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLocalSignedCandidate(t *testing.T) {
	f := makeFixture(t)
	_, stage := f.stage(t, 1)
	client := BundleClient(stage)
	cache := filepath.Join(f.dir, "candidate-cache")
	index, err := RefreshIndex(f.root, BundleSite, cache, client)
	if err != nil {
		t.Fatal(err)
	}
	info, raw, err := ReadCandidate(index, "synthetic/console/standard", "ROMD Synthetic Console", client)
	if err != nil || len(raw) != info.Bytes || info.PublicationVersion != 1 {
		t.Fatalf("candidate: %+v %v", info, err)
	}
	for _, id := range []string{"redump/psx/discs", "missing"} {
		if _, _, err = ReadCandidate(index, id, "ROMD Synthetic Console", client); err == nil {
			t.Fatal("wrong ID accepted")
		}
	}
	if _, _, err = ReadCandidate(index, info.CatalogID, "Sony - PlayStation", client); err == nil {
		t.Fatal("wrong platform header accepted")
	}
	c := index.Snapshot.Catalogs[info.CatalogID]
	asset := index.Downloads[c.Artifact.Path]
	asset.Bytes++
	index.Downloads[c.Artifact.Path] = asset
	if _, _, err = ReadCandidate(index, info.CatalogID, info.Name, client); err == nil {
		t.Fatal("inconsistent artifact accepted")
	}
	asset.Bytes--
	index.Downloads[c.Artifact.Path] = asset
	if err = os.WriteFile(filepath.Join(stage, "assets", filepath.Base(c.Artifact.Path)), []byte("tampered"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, _, err = ReadCandidate(index, info.CatalogID, info.Name, client); err == nil {
		t.Fatal("tampered bytes accepted")
	}
	c.Health = "failed"
	index.Snapshot.Catalogs[info.CatalogID] = c
	if _, _, err = ReadCandidate(index, info.CatalogID, info.Name, client); err == nil {
		t.Fatal("tampered retained artifact accepted")
	}
	c.Artifact = nil
	index.Snapshot.Catalogs[info.CatalogID] = c
	if _, _, err = ReadCandidate(index, info.CatalogID, info.Name, client); err == nil {
		t.Fatal("missing retained artifact accepted")
	}

	// A different independent trust root must reject an otherwise intact signed bundle.
	other := makeFixture(t)
	if _, err = RefreshIndex(other.root, BundleSite, filepath.Join(f.dir, "wrong-cache"), client); err == nil {
		t.Fatal("untrusted bundle accepted")
	}
}
