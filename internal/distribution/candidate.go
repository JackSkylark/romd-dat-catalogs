package distribution

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/MoonlarkStudios/romd-dat-catalogs/internal/publisher"
)

// Candidate is an authenticated document ready for ROMD's separate semantic review.
type Candidate struct {
	SystemID           string `json:"systemId,omitempty"`
	CatalogID          string `json:"catalogId"`
	Name               string `json:"name"`
	SHA256             string `json:"sha256"`
	Bytes              int    `json:"bytes"`
	PublicationVersion int64  `json:"publicationVersion"`
	SourceURL          string `json:"sourceUrl"`
}

// ReadCandidate accepts only an index returned by RefreshIndex. Catalog selection
// binds the operator's stable identity and expected name to the signed artifact.
func ReadCandidate(index Index, id, name string, client *http.Client) (Candidate, []byte, error) {
	var result Candidate
	catalog, ok := index.Snapshot.Catalogs[id]
	if !ok || catalog.Name != name || name == "" {
		return result, nil, errors.New("catalog identity is not present in the signed index")
	}
	var systemID string
	if index.Definitions != nil {
		if err := index.Definitions.Validate(); err != nil {
			return result, nil, err
		}
		definition, ok := index.Definitions.Catalogs[id]
		if ok && definition.ExpectedName == name {
			systemID = definition.SystemID
		} else if id != "synthetic/console/standard" {
			return result, nil, errors.New("catalog system binding is invalid")
		}
	}
	if catalog.Artifact == nil {
		return result, nil, errors.New("no retained catalog artifact is available; keep the installed version and retry later")
	}
	ref := catalog.Artifact
	asset, ok := index.Downloads[ref.Path]
	if !ok || !strings.HasSuffix(ref.Path, ".dat") || asset.Bytes != ref.Bytes || asset.SHA256 != ref.SHA256 || ref.Bytes < 1 || ref.Bytes > publisher.MaxDocument {
		return result, nil, errors.New("catalog artifact binding is invalid")
	}
	if index.Format == GitFormat && !gitAssetURL.MatchString(asset.URL) {
		return result, nil, errors.New("expected immutable Git document URL")
	}
	raw, err := FetchAsset(asset, client)
	if err != nil {
		return result, nil, err
	}
	document, _, err := publisher.Document(raw, name)
	if err != nil {
		return result, nil, err
	}
	if len(document) != ref.Bytes || publisher.Hash(document) != ref.SHA256 {
		return result, nil, errors.New("expected an exact extracted DAT document")
	}
	result = Candidate{SystemID: systemID, CatalogID: id, Name: name, SHA256: ref.SHA256, Bytes: ref.Bytes, PublicationVersion: index.Version}
	if catalog.Provenance != nil {
		result.SourceURL = catalog.Provenance.SourceURL
	}
	return result, document, nil
}

// BundleClient verifies an operator-supplied local Stage bundle using the same TUF
// updater as HTTPS publication. File presence conveys no trust: callers still
// supply an independent root and retain their metadata cache. No network is used.
func BundleClient(bundle string) *http.Client {
	return &http.Client{Timeout: 30 * time.Second, Transport: bundleTransport{bundle}}
}

const BundleSite = "https://bundle.invalid"

type bundleTransport struct{ root string }

func (b bundleTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	var relative string
	if r.Method != http.MethodGet {
		return nil, errors.New("bundle supports GET only")
	}
	if r.URL.Scheme == "https" && r.URL.Host == "bundle.invalid" && (strings.HasPrefix(r.URL.Path, "/metadata/") || strings.HasPrefix(r.URL.Path, "/targets/")) {
		relative = "site/" + strings.TrimPrefix(r.URL.Path, "/")
	} else if assetURL.MatchString(r.URL.String()) {
		relative = "assets/" + filepath.Base(r.URL.Path)
	} else {
		return nil, errors.New("unexpected bundle resource")
	}
	if strings.Contains(relative, "..") || strings.Contains(relative, "\\") || r.URL.RawQuery != "" {
		return nil, errors.New("invalid bundle resource")
	}
	file, err := os.Open(filepath.Join(b.root, filepath.FromSlash(relative)))
	if os.IsNotExist(err) {
		return &http.Response{StatusCode: 404, Header: make(http.Header), Body: io.NopCloser(strings.NewReader("")), Request: r}, nil
	}
	if err != nil {
		return nil, err
	}
	stat, err := file.Stat()
	if err != nil || !stat.Mode().IsRegular() {
		file.Close()
		return nil, fmt.Errorf("invalid bundle file")
	}
	return &http.Response{StatusCode: 200, Header: make(http.Header), Body: file, ContentLength: stat.Size(), Request: r}, nil
}
