package rest

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

type Router[S Session] struct {
	gin      *gin.Engine
	sessions NewSession[S]
}

func NewRouter[S Session](e *gin.Engine, fn NewSession[S]) *Router[S] {
	e.Use(func(c *gin.Context) {
		r := &Request{gin: c, Request: c.Request, Response: c.Writer}
		s, err := fn(r)
		if err != nil {
			c.Set(responseError, err)
			return
		}
		if v := reflect.ValueOf(s); v.Kind() == reflect.Ptr && !v.IsNil() {
			c.Set(session, s)
		}
		c.Set(request, r)
	})
	return &Router[S]{e, fn}
}

func (s *Router[S]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.gin.ServeHTTP(w, r)
}

func (s *Router[S]) Path(name string, m ...Middleware) *Path[S] {
	return NewPath(s, name, m...)
}

type Middleware = gin.HandlerFunc

type Error struct{ error }

func Err(space, message string, args ...any) Error {
	message = fmt.Sprintf("rest %s: %s", space, message)
	return Error{fmt.Errorf(message, args...)}
}
