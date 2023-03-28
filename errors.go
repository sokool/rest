package rest

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Errors(c *gin.Context) {
	c.Next()

	v, ok := c.Get(responseError)
	if !ok {
		return
	}
	err, ok := v.(error)
	if !ok {
		return
	}

	var m = err.Error()
	var t string
	var s = http.StatusBadRequest

	if str.Contains(m, "access denied", "access forbidden") {
		s = http.StatusForbidden
	}

	switch p := strings.Split(m, ":"); {
	case len(p) >= 2:
		t, m = p[0], strings.Join(p[1:], ":")
	}

	c.Writer.Header().Set("content-type", "application/json")
	c.Writer.WriteHeader(s)

	fmt.Fprintf(c.Writer, `{"message":"%s", "type": "%s"}`, strings.TrimSpace(m), strings.TrimSpace(t))
}
