package rest

import (
	"context"
	"encoding"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/sokool/errors"
)

type Request[S Session] struct {
	Session  S
	Request  *http.Request
	Response http.ResponseWriter
}

func NewRequest[S Session](s Sessions[S], w http.ResponseWriter, r *http.Request) (*Request[S], error) {
	var x = &Request[S]{Request: r, Response: w}
	var err error
	if x.Session, err = s.Read(r); err != nil {
		return nil, err
	}

	return x, nil
}

func (r *Request[S]) Body(to any) error {
	defer r.Request.Body.Close()
	switch err := json.NewDecoder(r.Request.Body).Decode(to); {
	case err == io.EOF:
		if j, ok := to.(json.Unmarshaler); ok {
			return j.UnmarshalJSON([]byte{})
		}
	case err != nil:
		return err
	}
	return nil
}

func (r *Request[S]) Param(name string, to ...any) (s string, err error) {
	if s = httprouter.ParamsFromContext(r.Request.Context()).ByName(name); s == "" {
		return s, ErrValueNotFound
	}
	if len(to) != 0 {
		if err = r.decode(s, to[0]); err != nil {
			return s, Err("decode failed %w", err)
		}
	}
	return
}

func (r *Request[S]) Query(name string, to ...any) (s string, err error) {
	if s = r.Request.URL.Query().Get(name); s == "" {
		return s, ErrValueNotFound
	}
	if len(to) != 0 {
		if err = r.decode(s, to[0]); err != nil {
			return s, Err("decode failed %w", err)
		}
	}
	return
}

func (r *Request[S]) Header(name string, to ...any) (s string, err error) {
	if s = r.Request.Header.Get(name); s == "" {
		return s, ErrValueNotFound
	}
	if len(to) != 0 {
		if err = r.decode(s, to[0]); err != nil {
			return s, Err("decode failed %w", err)
		}
	}
	return s, nil
}

func (r *Request[S]) Method(names ...string) bool {
	for i := range names {
		if names[i] == r.Request.Method {
			return true
		}
	}
	return false
}

func (r *Request[S]) decode(from string, to any) error {
	var err error
	switch v := to.(type) {
	case *string:
		*v = from
	case *[]string:
		*v = strings.Split(from, "|")
	case *int:
		var n int64
		n, err = strconv.ParseInt(from, 2, 64)
		*v = int(n)
	case *float64:
		var n float64
		n, err = strconv.ParseFloat(from, 64)
		*v = n
	case encoding.TextUnmarshaler:
		err = v.UnmarshalText([]byte(from))
	case json.Unmarshaler:
		err = v.UnmarshalJSON([]byte(from))
	default:
		err = json.Unmarshal([]byte(from), to)
	}

	return err
}

var Err = errors.Errorf

var ErrValueNotFound = Err("rest: value not found")

func Read[T any](from *http.Request, key string) (T, bool) {
	var ok bool
	var t T
	var v = from.Context().Value(key)
	if v == nil || reflect.ValueOf(v).IsZero() {
		return t, false
	}
	if t, ok = v.(T); ok {
		return t, true
	}
	return t, false
}

func Write[T any](key string, value T, to *http.Request) {
	t := to.WithContext(context.WithValue(to.Context(), key, value))
	*to = *t
}
