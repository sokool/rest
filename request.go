package rest

import (
	"encoding"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Request struct {
	*http.Request
	Response http.ResponseWriter
	gin      *gin.Context
	errors   []error
}

func (r *Request) Body(to any) error {
	defer r.Request.Body.Close()
	if err := json.NewDecoder(r.Request.Body).Decode(to); err != nil && err != io.EOF {
		return err
	}

	return nil
}

func (r *Request) Value(name string, to any) error {
	from, ok := r.gin.Get(name)
	if !ok {
		return nil
	}

	tv := reflect.ValueOf(to)
	if k := tv.Kind(); k != reflect.Pointer {
		return Err("request", "%s value need to be read by reference, %s given", name, k)
	}
	tv = tv.Elem()

	fv := reflect.ValueOf(from)
	if tv.Type() != fv.Type() {
		return Err("request", "%s value type mismatch - expected %s but %s has been given", name, fv.Type(), tv.Type())
	}

	tv.Set(fv)

	return nil
}

func (r *Request) Param(name string, to any) (err error) {
	s := r.gin.Param(name)
	if s == "" {
		return Err("request", "%s param not found", name)
	}
	switch v := to.(type) {
	case *string:
		if v == nil {
			return Err("request", "reading %s string parameter failed", name)
		}
		*v = s
	case encoding.TextUnmarshaler:
		err = v.UnmarshalText([]byte(s))
	default:
		err = json.Unmarshal([]byte(s), to)
	}
	if err != nil {
		err = Err("request", "%s param decode failed %w", name, err)
	}
	return
}

func (r *Request) Query(name string, to any) (err error) {
	var s string
	if s = r.Request.URL.Query().Get(name); s == "" {
		return
	}

	switch v := to.(type) {
	case *string:
		*v = s
	case *[]string:
		*v = strings.Split(s, "|")
	case *int:
		var n float64
		n, err = strconv.ParseFloat(s, 64)
		*v = int(n)
	case *float64:
		var n float64
		n, err = strconv.ParseFloat(s, 64)
		*v = n
	case *time.Time:
		*v, err = time.Parse(time.RFC3339, s)
	case encoding.TextUnmarshaler:
		err = v.UnmarshalText([]byte(s))
	default:
		if err = r.gin.BindQuery(to); err != nil {
			err = Err("request", "param decode failed %w", err)
		}
	}
	return
}

func (r *Request) Header(name string, to ...*string) bool {
	s := r.Request.Header.Get(name)
	if len(to) != 0 {
		*to[0] = s
	}
	return s != ""
}

func (r *Request) Method(names ...string) bool {
	for i := range names {
		if names[i] == r.Request.Method {
			return true
		}
	}
	return false
}
