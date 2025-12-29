package middleware

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRequestLogger(t *testing.T) {
	// Capture log output
	var buf strings.Builder
	oldOutput := log.Writer()
	log.SetOutput(&buf)
	defer log.SetOutput(oldOutput)

	r := gin.New()
	r.Use(RequestLogger())
	r.GET("/test-log", func(c *gin.Context) {
		c.Status(200)
	})

	req, _ := http.NewRequest("GET", "/test-log", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, buf.String(), "method=GET path=/test-log status=200")
}
