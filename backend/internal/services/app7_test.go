package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
