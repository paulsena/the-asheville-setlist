package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// BodyLimit limits request body size in megabytes
func BodyLimit(megabytes int64) gin.HandlerFunc {
	maxBytes := megabytes << 20
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)
		c.Next()
	}
}
