package nointro

import (
	"errors"
	"html"
	"net/url"
	"regexp"
	"strings"

	"github.com/MoonlarkStudios/romd-dat-catalogs/internal/definitions"
)

// This intentionally recognizes only the reviewed form shape. It is not an HTML
// browser: changed controls fail closed instead of guessing new export defaults.
var forms = regexp.MustCompile(`(?is)<form\b([^>]*)>(.*?)</form\s*>`)
var inputs = regexp.MustCompile(`(?is)<input\b([^>]*)>`)
var attributes = regexp.MustCompile(`([^\s=/>]+)(?:\s*=\s*(?:"([^"]*)"|'([^']*)'|([^\s>]+)))?`)
var prepareName = regexp.MustCompile(`^download_dat_[0-9]{4}-[0-9]{2}-[0-9]{2}$`)
var downloadName = regexp.MustCompile(`^[0-9a-f]{32}$`)
var selects = regexp.MustCompile(`(?is)<select\b([^>]*)>`)
var options = regexp.MustCompile(`(?is)<option\b([^>]*)>([^<]*)</option\s*>`)

func attrs(raw string) (map[string]string, error) {
	a := map[string]string{}
	for _, m := range attributes.FindAllStringSubmatch(raw, -1) {
		key := strings.ToLower(m[1])
		if _, exists := a[key]; exists {
			return nil, errors.New("duplicate attribute")
		}
		value := m[2]
		if m[3] != "" {
			value = m[3]
		}
		if m[4] != "" {
			value = m[4]
		}
		a[key] = html.UnescapeString(value)
	}
	return a, nil
}

// Policy v2: complete public standard, default representation, non-merged,
// canonical names, every inclusion category, all regions/languages/specials.
// The live Aftermarket checkbox's value is literally 0; presence selects it.
func completeSelection() url.Values {
	return url.Values{
		"format": {"0"}, "collection": {"0"}, "naming": {"0"},
		"inc_bios": {"1"}, "release_1": {"1"}, "release_2": {"1"},
		"license_1": {"1"}, "license_2": {"1"}, "license_0": {"1"},
		"lifespan_1": {"1"}, "lifespan_2": {"0"}, "inc_adult": {"1"},
		"storage_1": {"1"}, "storage_2": {"1"}, "inc_nodump": {"1"},
		"special1_filter": {"all_specials1"}, "special2_filter": {"all_specials2"},
		"language_filter": {"all_languages"}, "region_filter": {"all_regions"},
	}
}
func prepareForm(raw []byte, c definitions.Catalog) (url.Values, error) {
	bad := errors.New("unrecognized standard DAT form")
	var content string
	count := 0
	for _, f := range forms.FindAllStringSubmatch(string(raw), -1) {
		a, err := attrs(f[1])
		if err != nil {
			return nil, bad
		}
		if a["name"] == "main_form" {
			if a["method"] != "post" || a["action"] != "" {
				return nil, bad
			}
			content = f[2]
			count++
		}
	}
	if count != 1 {
		return nil, bad
	}
	for _, selectTag := range selects.FindAllStringSubmatch(content, -1) {
		a, err := attrs(selectTag[1])
		if err != nil || (a["name"] != "system_selection" && a["name"] != "sys_list_order") {
			return nil, bad
		}
	}
	selectorName := c.ExpectedName
	// Official selector labels omit the representation suffix in DAT headers.
	if c.ProviderSystemID == "24" && c.ExpectedName == "Nintendo - Nintendo 64 (BigEndian)" {
		selectorName = "Nintendo - Nintendo 64"
	}
	if c.ProviderSystemID == "45" && c.ExpectedName == "Nintendo - Nintendo Entertainment System (Headered)" {
		selectorName = "Nintendo - Nintendo Entertainment System"
	}
	if (c.ProviderSystemID == "28" && c.ExpectedName == "Nintendo - Nintendo DS (Decrypted)") || (c.ProviderSystemID == "64" && c.ExpectedName == "Nintendo - Nintendo 3DS (Decrypted)") {
		selectorName = strings.TrimSuffix(c.ExpectedName, " (Decrypted)")
	}
	identity := false
	for _, o := range options.FindAllStringSubmatch(content, -1) {
		a, err := attrs(o[1])
		if err != nil {
			return nil, bad
		}
		if _, ok := a["selected"]; ok && a["value"] == c.ProviderSystemID && strings.TrimSpace(html.UnescapeString(o[2])) == selectorName {
			identity = true
		}
	}
	if !identity {
		return nil, bad
	}
	values := completeSelection()
	// Reviewed per-system controls: GBC has no collection selector; GBA
	// exposes numbered, x-ROM and z-ROM categories instead of adult filtering.
	if c.ProviderSystemID == "47" || c.ProviderSystemID == "45" || c.ProviderSystemID == "32" || c.ProviderSystemID == "24" {
		delete(values, "collection")
	}
	if c.ProviderSystemID == "23" {
		delete(values, "inc_adult")
		values.Set("numbered", "0")
		values.Set("inc_xroms", "1")
		values.Set("inc_zroms", "1")
	}
	// NES format 0 is Headered; omit the header-removal plugin. N64 format 0
	// is BigEndian. These reviewed systems expose no collection selector.
	if c.ProviderSystemID == "45" {
		values.Set("header_plugin", "0")
	}
	if c.ProviderSystemID == "24" {
		delete(values, "inc_adult")
	}
	// Reviewed batch-3 forms have no collection/adult selectors. Only SMS
	// exposes nodump; SMS includes x/z-ROM groups and GG includes x-ROMs.
	switch c.ProviderSystemID {
	case "26", "25", "17", "12":
		delete(values, "collection")
		delete(values, "inc_adult")
		if c.ProviderSystemID != "26" {
			delete(values, "inc_nodump")
		}
		if c.ProviderSystemID == "26" || c.ProviderSystemID == "25" {
			values.Set("inc_xroms", "1")
		}
		if c.ProviderSystemID == "26" {
			values.Set("inc_zroms", "1")
		}
	}
	// DS and 3DS use explicit decrypted format 0, canonical unnumbered names,
	// and all x/z-ROM groups. 3DS exposes no license or lifespan filters.
	if c.ProviderSystemID == "28" || c.ProviderSystemID == "64" {
		delete(values, "collection")
		delete(values, "inc_adult")
		values.Set("numbered", "0")
		values.Set("inc_xroms", "1")
		values.Set("inc_zroms", "1")
		if c.ProviderSystemID == "64" {
			for _, name := range []string{"license_0", "license_1", "license_2", "lifespan_1", "lifespan_2"} {
				delete(values, name)
			}
		}
	}
	// N64 currently omits nodump when no such records are offered. Include the
	// category whenever its control exists; never submit an invented control.
	if c.ProviderSystemID == "24" {
		delete(values, "inc_nodump")
	}
	found := map[string]int{}
	submit := ""
	for _, i := range inputs.FindAllStringSubmatch(content, -1) {
		a, err := attrs(i[1])
		if err != nil {
			return nil, bad
		}
		if _, disabled := a["disabled"]; disabled {
			return nil, bad
		}
		if _, hidden := a["hidden"]; hidden {
			return nil, bad
		}
		name := a["name"]
		// Some systems expose a missing-in-action filter. Include these records
		// with canonical names; legacy forms without this control remain valid.
		if name == "inc_mia" || (c.ProviderSystemID == "24" && name == "inc_nodump") {
			values.Set(name, "1")
		}
		if prepareName.MatchString(name) && a["type"] == "submit" && a["value"] == "Prepare" {
			if submit != "" {
				return nil, bad
			}
			submit = name
			continue
		}
		if want, ok := values[name]; ok {
			if a["type"] != "radio" && a["type"] != "checkbox" {
				return nil, bad
			}
			if a["value"] == want[0] {
				found[name]++
			}
			continue
		}
		// The custom option arrays are deliberately unused; their all-* radio wins.
		if a["type"] == "checkbox" && (strings.HasPrefix(name, "specials1[") || strings.HasPrefix(name, "specials2[") || strings.HasPrefix(name, "languages[") || strings.HasPrefix(name, "regions[")) {
			continue
		}
		return nil, bad
	}
	if submit == "" {
		return nil, bad
	}
	for name := range values {
		if found[name] != 1 {
			return nil, bad
		}
	}
	values.Set(submit, "Prepare")
	return values, nil
}
func downloadForm(raw []byte) (url.Values, error) {
	bad := errors.New("unrecognized download form")
	var selected string
	for _, f := range forms.FindAllStringSubmatch(string(raw), -1) {
		a, err := attrs(f[1])
		if err != nil {
			return nil, bad
		}
		for _, i := range inputs.FindAllStringSubmatch(f[2], -1) {
			control, err := attrs(i[1])
			if err != nil {
				return nil, bad
			}
			if control["value"] != "Download!!" {
				continue
			}
			if control["style"] == "display:none;" {
				continue
			}
			if _, disabled := control["disabled"]; disabled {
				return nil, bad
			}
			if _, hidden := control["hidden"]; hidden {
				return nil, bad
			}
			if control["style"] != "" || control["type"] != "submit" || !downloadName.MatchString(control["name"]) || a["method"] != "post" || a["action"] != "" || len(inputs.FindAllStringSubmatch(f[2], -1)) != 1 || selected != "" {
				return nil, bad
			}
			selected = control["name"]
		}
	}
	if selected == "" {
		return nil, bad
	}
	return url.Values{selected: {"Download!!"}}, nil
}
