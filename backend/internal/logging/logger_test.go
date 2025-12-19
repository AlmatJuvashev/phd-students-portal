package logging

import (
	"bytes"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogging(t *testing.T) {
	// Capture output
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)

	Info("info msg", "key", "val")
	assert.Contains(t, buf.String(), "INFO")
	assert.Contains(t, buf.String(), "info msg")
	assert.Contains(t, buf.String(), "val")
	
	buf.Reset()
	Warn("warn msg")
	assert.Contains(t, buf.String(), "WARN")
	assert.Contains(t, buf.String(), "warn msg")

	buf.Reset()
	Error("error msg")
	assert.Contains(t, buf.String(), "ERROR")
	assert.Contains(t, buf.String(), "error msg")
}
