package services

import (
	"encoding/json"
)


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

func summarizeApp7(sections App7Sections) map[string]int {
	return map[string]int{
		"wos_scopus":  len(sections.WosScopus),
		"kokson":      len(sections.Kokson),
		"conferences": len(sections.Conferences),
		"ip":          len(sections.IP),
	}
}
