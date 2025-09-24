package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestLogger prints structured logs per request.
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		lat := time.Since(start)
		log.Printf(`method=%s path=%s status=%d latency=%s`, c.Request.Method, c.Request.URL.Path, c.Writer.Status(), lat)
	}
}
