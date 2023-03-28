package rest

import (
	"github.com/gin-gonic/gin"
)

type Path[S Session] struct {
	router  *gin.RouterGroup
	session NewSession[S]
}

func NewPath[S Session](r *Router[S], name string, m ...Middleware) *Path[S] {
	return &Path[S]{
		router:  r.gin.Group(name).Use(m...).(*gin.RouterGroup),
		session: r.sessions,
	}
}

func (p *Path[S]) Path(name string, m ...Middleware) *Path[S] {
	return &Path[S]{
		router:  p.router.Group(name).Use(m...).(*gin.RouterGroup),
		session: p.session,
	}
}

func (p *Path[S]) Handle(name, method string, h Handler[S], m ...Middleware) *Path[S] {
	p.router.Handle(method, name, func(c *gin.Context) {
		var v any
		var ok bool
		if v, ok = c.Get(request); !ok {
			return
		}
		r, ok := v.(*Request)
		if !ok {
			return
		}
		var s S
		if v, ok = c.Get(session); ok {
			s, _ = v.(S)
		}

		data, err := h.Do(s, r)
		c.Set(responseBody, data)
		c.Set(responseError, err)
	})

	return p
}

const (
	responseBody  = "responseBody"
	responseError = "responseError"
	session       = "session"
	request       = "request"
)
