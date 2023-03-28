package rest

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stoewer/go-strcase"
	"net/http"
	"reflect"
	"regexp"
)

func JSONResponse(c *gin.Context) {
	c.Next()
	err, _ := c.Get(responseError)
	if _, ok := err.(error); ok {
		return
	}

	p, ok := c.Get(responseBody)
	if !ok {
		return
	}

	if p == nil || p == "" || reflect.ValueOf(p).IsZero() {
		c.Writer.WriteHeader(http.StatusNoContent)
		return
	}

	c.Writer.Header().Set("content-type", "application/json")
	s := c.GetHeader("json-style")
	if s == "" {
		s = "camel"
	}
	if err = json.NewEncoder(c.Writer).Encode(&jsonStyle{s, p}); err != nil {
		c.Set(responseError, err)
	}
}

var words = regexp.MustCompile(`\"(\w+)\":`)

type jsonStyle struct {
	name  string
	value any
}

func (j *jsonStyle) MarshalJSON() ([]byte, error) {
	b, err := json.Marshal(j.value)
	if err != nil {
		return nil, err
	}
	var f func([]byte) []byte
	switch j.name {
	case "camel":
		f = func(b []byte) []byte { return []byte(strcase.LowerCamelCase(string(b))) }
	case "snake":
		f = func(b []byte) []byte { return []byte(strcase.SnakeCase(string(b))) }
	case "kebab":
		f = func(b []byte) []byte { return []byte(strcase.KebabCase(string(b))) }
	default:
		return b, nil
	}
	return words.ReplaceAllFunc(b, f), err
}
