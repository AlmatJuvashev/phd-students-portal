package storage

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// BuildNodeObjectKey generates a deterministic prefix for student node uploads.
func BuildNodeObjectKey(userID, nodeID, slotKey, filename string) string {
	cleanName := sanitizeFilename(filename)
	segmentUser := sanitizeSegment(userID)
	segmentNode := sanitizeSegment(nodeID)
	segmentSlot := sanitizeSegment(slotKey)
	if segmentUser == "" {
		segmentUser = "student"
	}
	if segmentNode == "" {
		segmentNode = "node"
	}
	if segmentSlot == "" {
		segmentSlot = "slot"
	}
	stamp := time.Now().UTC().Format("20060102")
	return fmt.Sprintf("nodes/%s/%s/%s/%s-%s-%s", segmentUser, segmentNode, segmentSlot, stamp, uuid.NewString(), cleanName)
}

// BuildDocumentObjectKey generates an object key for generic documents.
func BuildDocumentObjectKey(documentID, filename string) string {
	cleanName := sanitizeFilename(filename)
	segmentDoc := sanitizeSegment(documentID)
	if segmentDoc == "" {
		segmentDoc = "doc"
	}
	stamp := time.Now().UTC().Format("20060102")
	return fmt.Sprintf("documents/%s/%s-%s", segmentDoc, stamp, cleanName)
}

func sanitizeFilename(name string) string {
	base := filepath.Base(strings.TrimSpace(name))
	if base == "" {
		base = "file"
	}
	base = strings.ToLower(base)
	var b strings.Builder
	for _, r := range base {
		switch {
		case r >= 'a' && r <= 'z':
			b.WriteRune(r)
		case r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == '.' || r == '-' || r == '_':
			b.WriteRune(r)
		default:
			b.WriteRune('-')
		}
	}
	out := strings.Trim(b.String(), "-_.")
	if out == "" {
		out = "file"
	}
	if len(out) > 80 {
		out = out[len(out)-80:]
	}
	return out
}

func sanitizeSegment(seg string) string {
	seg = strings.TrimSpace(seg)
	seg = strings.ToLower(seg)
	var b strings.Builder
	for _, r := range seg {
		switch {
		case r >= 'a' && r <= 'z':
			b.WriteRune(r)
		case r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == '-' || r == '_':
			b.WriteRune(r)
		default:
		}
	}
	out := strings.Trim(b.String(), "-_")
	if len(out) > 48 {
		out = out[:48]
	}
	return out
}
