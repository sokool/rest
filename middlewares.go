package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/sokool/log"
	"github.com/stoewer/go-strcase"
)

type Middleware func(http.Handler) http.HandlerFunc

func JSON(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
		w.Header().Set("Content-Type", "application/json")

		var bdy any
		if err, ok := Read[error](r, responseErr); ok {
			bdy = err
			w.WriteHeader(http.StatusBadRequest)
		} else if res, ok := Read[any](r, responseBody); ok {
			bdy = res
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		s := r.Header.Get("json-style")
		if s == "" {
			s = "camel"
		}

		if err := json.NewEncoder(w).Encode(&jsonStyle{s, bdy}); err != nil {
			fmt.Println("damn, failed", err)
		}
	}
}

func Errors(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)

		err, ok := Read[error](r, responseErr)
		if !ok || err == nil {
			return
		}

		var m = err.Error()
		var t string
		var s = http.StatusBadRequest

		for _, e := range []string{"access denied", "access forbidden"} {
			if strings.Contains(m, e) {
				s = http.StatusForbidden
				break
			}
		}

		switch p := strings.Split(m, ":"); {
		case len(p) >= 2:
			t, m = p[0], strings.Join(p[1:], ":")
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(s)

		fmt.Fprintf(w, `{"message":"%s", "type": "%s"}`, strings.TrimSpace(m), strings.TrimSpace(t))
	}
}

func Logger(verbose bool, w ...io.Writer) Middleware {
	log := log.
		New(os.Stdout).
		Tag("rest").
		Options(log.Date | log.Time | log.Type | log.Type | log.Colors).
		Verbose(verbose)

	return func(h http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			n := time.Now()
			h.ServeHTTP(w, r)
			log.Printf("%s %s in %s", r.Method, r.URL, time.Since(n))
		}
	}
}

//func Renderer(h http.Handler) http.Handler {
//	type renderer interface {
//		Render(types ...string) (any, error)
//	}
//	return http.HandlerFunc(func(hw http.ResponseWriter, hr *http.Request) {
//		h.ServeHTTP(hw, hr)
//
//		var r renderer
//		var ok bool
//
//		if r, ok = Read[renderer](hr, responseBody); !ok {
//			return
//		}
//		//if p == nil || p == "" || reflect.ValueOf(p).IsZero() {
//		//	return
//		//}
//		//if r, ok = p.(view.Renderer); !ok {
//		//	return
//		//}
//
//		r.Render(hr.Header.Values("render")...)
//		b, err := view.JSON(r, c.Request.Header.Values("render")...)
//		if err != nil {
//			c.Set(responseError, err)
//			return
//		}
//
//		Write(responseBody)
//		c.Set(responseBody, b.JSON())
//	})
//}

//func Auth(c *gin.Context) {
//	v, ok := c.Get(session)
//	if !ok {
//		c.Writer.WriteHeader(http.StatusUnauthorized)
//		c.Abort()
//		return
//	}
//	s, ok := v.(Session)
//	if !ok {
//		c.JSON(http.StatusInternalServerError, gin.H{"rest": "session type mismatch"})
//		c.Abort()
//		return
//	}
//	if s.IsExpired() {
//		c.Writer.WriteHeader(http.StatusUnauthorized)
//		c.Abort()
//		return
//	}
//}

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
