package models

type SearchResult struct {
	Type        string `json:"type"`                // "student", "document", "message"
	ID          string `json:"id"`                  // ID to navigate to
	Title       string `json:"title"`               // Display title (Name, Filename, etc.)
	Subtitle    string `json:"subtitle"`            // Secondary info (Email, Node Name, etc.)
	Description string `json:"description"`         // Context (Message snippet, etc.)
	Link        string `json:"link"`                // Frontend route
	Metadata    any    `json:"metadata,omitempty"`
}
