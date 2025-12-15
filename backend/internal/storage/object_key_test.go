package storage

import (
	"strings"
	"testing"
)

func TestBuildNodeObjectKey(t *testing.T) {
	tests := []struct {
		name     string
		userID   string
		nodeID   string
		slotKey  string
		filename string
		want     string // partial match
	}{
		{
			name:     "Standard",
			userID:   "user-1",
			nodeID:   "node-A",
			slotKey:  "submission",
			filename: "My Thesis.pdf",
			want:     "nodes/user-1/node-a/submission/",
		},
		{
			name:     "Sanitization",
			userID:   "User@One",
			nodeID:   "Node#2",
			slotKey:  "Slot Key!",
			filename: "../../Evil.exe",
			want:     "nodes/userone/node2/slotkey/",
		},
		{
			name:     "Empty Segments",
			userID:   "",
			nodeID:   "",
			slotKey:  "",
			filename: "",
			want:     "nodes/student/node/slot/",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := BuildNodeObjectKey(tc.userID, tc.nodeID, tc.slotKey, tc.filename)
			if !strings.HasPrefix(got, tc.want) {
				t.Errorf("BuildNodeObjectKey() = %v, want prefix %v", got, tc.want)
			}
			// Check filename at end
			if tc.filename == "My Thesis.pdf" && !strings.HasSuffix(got, "-my-thesis.pdf") {
				t.Errorf("Filename sanitization failed: %v", got)
			}
			if tc.filename == "" && !strings.HasSuffix(got, "-file") {
				t.Errorf("Empty filename fallback failed: %v", got)
			}
		})
	}
}

func TestBuildDocumentObjectKey(t *testing.T) {
	got := BuildDocumentObjectKey("Doc-123", "Specs v1.docx")
	
	if !strings.HasPrefix(got, "documents/doc-123/") {
		t.Errorf("BuildDocumentObjectKey prefix wrong: %v", got)
	}
	if !strings.HasSuffix(got, "-specs-v1.docx") {
		t.Errorf("BuildDocumentObjectKey suffix wrong: %v", got)
	}
}

func TestSanitizeFilename(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"Hello World.txt", "hello-world.txt"},
		{"../../../etc/passwd", "passwd"}, // filepath.Base takes last part
		{"foo/bar/baz.png", "baz.png"},
		{"", "file"},
		{"   ", "file"},
		{"a...b", "a...b"}, // dots are allowed in middle
		{"_Start", "start"}, // trim leading
		{"End_", "end"}, // trim trailing
	}

	for _, tc := range tests {
		got := sanitizeFilename(tc.input)
		if got != tc.want {
			t.Errorf("sanitizeFilename(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}
