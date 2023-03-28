package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Auth(c *gin.Context) {
	v, ok := c.Get(session)
	if !ok {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		c.Abort()
		return
	}
	s, ok := v.(Session)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"rest": "session type mismatch"})
		c.Abort()
		return
	}
	if s.IsExpired() {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		c.Abort()
		return
	}
}
