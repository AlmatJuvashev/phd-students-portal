package models

import "github.com/lib/pq"

// ProgramNodeConfig represents the strict structure of `JourneyNodeDefinition.Config` (JSONB).
// This is used by the Builder to validate and type-check the node configuration.
type ProgramNodeConfig struct {
	// For "formEntry" and "checklist"
	Fields []ProgramFieldDefinition `json:"fields,omitempty"`

	// For "checklist" specific notes
	Notes string `json:"notes,omitempty"`
	
	// For "cards"
	Slides []ProgramCardSlide `json:"slides,omitempty"`
}

// ProgramFieldDefinition defines a single input or check item in a node.
type ProgramFieldDefinition struct {
	Key         string            `json:"key"`
	Type        string            `json:"type"` // "text", "textarea", "boolean", "date", "file", "note", "select"
	Label       map[string]string `json:"label"`
	Required    bool              `json:"required"`
	Placeholder map[string]string `json:"placeholder,omitempty"`
	VisibleWhen string            `json:"visible_when,omitempty"` // e.g. "form.other_field == true"
	
	// Type-specific configs
	Options       []ProgramFieldOption `json:"options,omitempty"` // For select/radio
	MimeWhitelist pq.StringArray       `json:"mime_whitelist,omitempty"` // For file
	Validations   []ProgramValidation  `json:"validations,omitempty"`
}

type ProgramFieldOption struct {
	Value string            `json:"value"`
	Label map[string]string `json:"label"`
}

type ProgramValidation struct {
	Type    string `json:"type"` // "regex", "min", "max"
	Rule    string `json:"rule"`
	Message string `json:"message"`
}

type ProgramCardSlide struct {
	Key     string            `json:"key"`
	Title   map[string]string `json:"title"`
	Content map[string]string `json:"content"` // Markdown
	Image   string            `json:"image,omitempty"`
}
