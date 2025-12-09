package handlers

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"
)

type app7ValidationError struct {
	msg string
}

func (e *app7ValidationError) Error() string { return e.msg }

type App7Entry struct {
	Title         string   `json:"title"`
	Format        string   `json:"format,omitempty"`
	FormatOther   string   `json:"format_other,omitempty"`
	Journal       string   `json:"journal,omitempty"`
	Year          string   `json:"year,omitempty"`
	VolumeIssue   string   `json:"volume_issue,omitempty"`
	PagesOrSheets string   `json:"pages_or_sheets,omitempty"`
	DOI           string   `json:"doi,omitempty"`
	ISSNPrint     string   `json:"issn_print,omitempty"`
	ISSNOnline    string   `json:"issn_online,omitempty"`
	Coauthors     []string `json:"coauthors,omitempty"`
	Indexing      string   `json:"indexing,omitempty"`
	IndexingOther string   `json:"indexing_other,omitempty"`
	IPType        string   `json:"ip_type,omitempty"`
	IPTypeOther   string   `json:"ip_type_other,omitempty"`
	CertificateNo string   `json:"certificate_no,omitempty"`
	ISBN          string   `json:"isbn,omitempty"`
}

type App7Sections struct {
	WosScopus   []App7Entry `json:"wos_scopus"`
	Kokson      []App7Entry `json:"kokson"`
	Conferences []App7Entry `json:"conferences"`
	IP          []App7Entry `json:"ip"`
}

type App7Form struct {
	Sections     App7Sections   `json:"sections"`
	LegacyCounts map[string]int `json:"legacy_counts,omitempty"`
}

type app7Incoming struct {
	Sections     *App7Sections  `json:"sections"`
	LegacyCounts map[string]int `json:"legacy_counts"`
	WosScopus    []App7Entry    `json:"wos_scopus"`
	Kokson       []App7Entry    `json:"kokson"`
	Conferences  []App7Entry    `json:"conferences"`
	IP           []App7Entry    `json:"ip"`
}

type app7LegacyCounts struct {
	CountWosScopus   *int `json:"count_wos_scopus"`
	CountKokson      *int `json:"count_kokson"`
	CountConferences *int `json:"count_conferences"`
	CountIP          *int `json:"count_ip"`
}

var (
	doiPattern   = regexp.MustCompile(`(?i)^10\.[0-9]{4,9}/[-._;()/:A-Z0-9]+$`)
	issnPattern  = regexp.MustCompile(`^[0-9]{4}-[0-9]{3}[0-9Xx]$`)
	yearPattern  = regexp.MustCompile(`^(18|19|20)\d{2}([/\\-](18|19|20)\d{2})?$`)
	stringFields = []string{
		"title", "format", "format_other", "journal", "year", "volume_issue",
		"pages_or_sheets", "doi", "issn_print", "issn_online", "indexing",
		"indexing_other", "ip_type", "ip_type_other", "certificate_no", "isbn",
	}
)

func normalizeApp7Payload(raw json.RawMessage) ([]byte, *App7Form, error) {
	form, err := buildApp7Form(raw)
	if err != nil {
		return nil, nil, err
	}
	if err := sanitizeApp7Form(form); err != nil {
		return nil, nil, err
	}
	sanitized, err := json.Marshal(form)
	if err != nil {
		return nil, nil, err
	}
	return sanitized, form, nil
}

func buildApp7Form(raw json.RawMessage) (*App7Form, error) {
	var incoming app7Incoming
	_ = json.Unmarshal(raw, &incoming) // Ignore error, check content

	// Check if incoming is effectively empty (no sections, no legacy counts map, no flat arrays)
	isEmpty := incoming.Sections == nil && len(incoming.LegacyCounts) == 0 &&
		len(incoming.WosScopus) == 0 && len(incoming.Kokson) == 0 &&
		len(incoming.Conferences) == 0 && len(incoming.IP) == 0

	if isEmpty {
		// Try legacy counts-only structure
		var legacy app7LegacyCounts
		if errLegacy := json.Unmarshal(raw, &legacy); errLegacy == nil {
			// Check if any legacy field is present
			if legacy.CountWosScopus != nil || legacy.CountKokson != nil || legacy.CountConferences != nil || legacy.CountIP != nil {
				counts := map[string]int{}
				if legacy.CountWosScopus != nil {
					counts["wos_scopus"] = *legacy.CountWosScopus
				}
				if legacy.CountKokson != nil {
					counts["kokson"] = *legacy.CountKokson
				}
				if legacy.CountConferences != nil {
					counts["conferences"] = *legacy.CountConferences
				}
				if legacy.CountIP != nil {
					counts["ip"] = *legacy.CountIP
				}
				return &App7Form{Sections: App7Sections{}, LegacyCounts: counts}, nil
			}
		}
	}

	form := &App7Form{Sections: App7Sections{}}
	if incoming.LegacyCounts != nil {
		form.LegacyCounts = incoming.LegacyCounts
	} else {
		form.LegacyCounts = map[string]int{}
	}

	if incoming.Sections != nil {
		form.Sections = *incoming.Sections
	} else {
		form.Sections = App7Sections{
			WosScopus:   incoming.WosScopus,
			Kokson:      incoming.Kokson,
			Conferences: incoming.Conferences,
			IP:          incoming.IP,
		}
	}
	return form, nil
}

func sanitizeApp7Form(form *App7Form) error {
	if form.LegacyCounts == nil {
		form.LegacyCounts = map[string]int{}
	}

	var errs []string

	sanitize := func(entries []App7Entry, section string, requireIPType bool) []App7Entry {
		cleaned := make([]App7Entry, 0, len(entries))
		for idx, entry := range entries {
			normalized := normalizeEntry(entry)
			if err := validateEntry(normalized, requireIPType); err != nil {
				errs = append(errs, fmt.Sprintf("%s[%d]: %s", section, idx+1, err.Error()))
			}
			cleaned = append(cleaned, normalized)
		}
		if cleaned == nil {
			return []App7Entry{}
		}
		return cleaned
	}

	form.Sections.WosScopus = sanitize(form.Sections.WosScopus, "wos_scopus", false)
	form.Sections.Kokson = sanitize(form.Sections.Kokson, "kokson", false)
	form.Sections.Conferences = sanitize(form.Sections.Conferences, "conferences", false)
	form.Sections.IP = sanitize(form.Sections.IP, "ip", true)

	if len(errs) > 0 {
		sort.Strings(errs)
		return &app7ValidationError{msg: strings.Join(errs, "; ")}
	}

	return nil
}

func normalizeEntry(entry App7Entry) App7Entry {
	entry.Title = strings.TrimSpace(entry.Title)
	entry.Format = strings.TrimSpace(entry.Format)
	entry.FormatOther = strings.TrimSpace(entry.FormatOther)
	entry.Journal = strings.TrimSpace(entry.Journal)
	entry.Year = strings.TrimSpace(entry.Year)
	entry.VolumeIssue = strings.TrimSpace(entry.VolumeIssue)
	entry.PagesOrSheets = strings.TrimSpace(entry.PagesOrSheets)
	entry.DOI = strings.TrimSpace(entry.DOI)
	entry.ISSNPrint = strings.TrimSpace(entry.ISSNPrint)
	entry.ISSNOnline = strings.TrimSpace(entry.ISSNOnline)
	entry.Indexing = strings.TrimSpace(entry.Indexing)
	entry.IndexingOther = strings.TrimSpace(entry.IndexingOther)
	entry.IPType = strings.TrimSpace(entry.IPType)
	entry.IPTypeOther = strings.TrimSpace(entry.IPTypeOther)
	entry.CertificateNo = strings.TrimSpace(entry.CertificateNo)
	entry.ISBN = strings.TrimSpace(entry.ISBN)

	cleanedAuthors := make([]string, 0, len(entry.Coauthors))
	for _, name := range entry.Coauthors {
		name = strings.TrimSpace(name)
		if name != "" {
			cleanedAuthors = append(cleanedAuthors, name)
		}
	}
	entry.Coauthors = cleanedAuthors

	return entry
}

func validateEntry(entry App7Entry, requireIPType bool) error {
	if entry.Title == "" {
		return fmt.Errorf("title is required")
	}
	if len(entry.Title) > 1000 {
		return fmt.Errorf("title exceeds 1000 characters")
	}
	if entry.Year != "" && !yearPattern.MatchString(entry.Year) {
		return fmt.Errorf("year must be in format YYYY or YYYY/YYYY")
	}
	if entry.DOI != "" && !doiPattern.MatchString(entry.DOI) {
		return fmt.Errorf("invalid DOI")
	}
	if entry.ISSNPrint != "" && !issnPattern.MatchString(entry.ISSNPrint) {
		return fmt.Errorf("invalid ISSN print")
	}
	if entry.ISSNOnline != "" && !issnPattern.MatchString(entry.ISSNOnline) {
		return fmt.Errorf("invalid ISSN online")
	}
	if requireIPType && entry.IPType == "" && entry.IPTypeOther == "" {
		return fmt.Errorf("type is required for intellectual property entries")
	}
	return nil
}

func summarizeApp7(sections App7Sections) map[string]int {
	return map[string]int{
		"wos_scopus":  len(sections.WosScopus),
		"kokson":      len(sections.Kokson),
		"conferences": len(sections.Conferences),
		"ip":          len(sections.IP),
	}
}
