package services

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApp7_NormalizePayload(t *testing.T) {
	// Case 1: Valid payload with sections
	raw := json.RawMessage(`{
		"sections": {
			"wos_scopus": [
				{"title": " Paper 1 ", "year": "2023", "doi": "10.1234/5678"}
			],
			"ip": [
				{"title": "Patent 1", "ip_type": "patent"}
			]
		}
	}`)
	sanitized, form, err := normalizeApp7Payload(raw)
	require.NoError(t, err)
	assert.NotNil(t, form)
	assert.Len(t, form.Sections.WosScopus, 1)
	assert.Equal(t, "Paper 1", form.Sections.WosScopus[0].Title)
	assert.Len(t, form.Sections.IP, 1)
	assert.NotNil(t, sanitized)

	// Case 2: Legacy counts payload
	rawLegacy := json.RawMessage(`{
		"count_wos_scopus": 5,
		"count_kokson": 3
	}`)
	sanitizedLegacy, formLegacy, err := normalizeApp7Payload(rawLegacy)
	require.NoError(t, err)
	assert.NotNil(t, formLegacy)
	assert.Equal(t, 5, formLegacy.LegacyCounts["wos_scopus"])
	assert.Equal(t, 3, formLegacy.LegacyCounts["kokson"])
	assert.NotNil(t, sanitizedLegacy)

	// Case 3: Invalid entry (missing title)
	rawInvalid := json.RawMessage(`{
		"sections": {
			"wos_scopus": [
				{"year": "2023"}
			]
		}
	}`)
	_, _, err = normalizeApp7Payload(rawInvalid)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "title is required")

	// Case 4: Invalid DOI
	rawInvalidDOI := json.RawMessage(`{
		"sections": {
			"wos_scopus": [
				{"title": "Paper", "doi": "invalid-doi"}
			]
		}
	}`)
	_, _, err = normalizeApp7Payload(rawInvalidDOI)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid DOI")
}

func TestApp7_Summarize(t *testing.T) {
	sections := App7Sections{
		WosScopus: []App7Entry{{Title: "1"}, {Title: "2"}},
		Kokson:    []App7Entry{{Title: "1"}},
	}
	summary := summarizeApp7(sections)
	assert.Equal(t, 2, summary["wos_scopus"])
	assert.Equal(t, 1, summary["kokson"])
	assert.Equal(t, 0, summary["conferences"])
}

func TestApp7_ValidateEntry(t *testing.T) {
	tests := []struct {
		name    string
		entry   App7Entry
		reqIP   bool
		wantErr string
	}{
		{
			name:  "Valid Entry",
			entry: App7Entry{Title: "Title", Year: "2023"},
			reqIP: false,
		},
		{
			name:    "Missing Title",
			entry:   App7Entry{Year: "2023"},
			reqIP:   false,
			wantErr: "title is required",
		},
		{
			name:    "Title Too Long",
			entry:   App7Entry{Title: string(make([]byte, 1001))},
			reqIP:   false,
			wantErr: "title exceeds 1000 characters",
		},
		{
			name:    "Invalid Year",
			entry:   App7Entry{Title: "Title", Year: "abcd"},
			reqIP:   false,
			wantErr: "year must be in format YYYY",
		},
		{
			name:    "Invalid ISSN",
			entry:   App7Entry{Title: "Title", ISSNPrint: "123"},
			reqIP:   false,
			wantErr: "invalid ISSN print",
		},
		{
			name:    "Missing IP Type",
			entry:   App7Entry{Title: "Patent"},
			reqIP:   true,
			wantErr: "type is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateEntry(tt.entry, tt.reqIP)
			if tt.wantErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
