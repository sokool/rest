package rest

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Router[S Session] struct {
	http       *httprouter.Router
	sessions   NewSession[S]
	middleware []Middleware
	path       string
}

func NewRouter(m ...Middleware) *Router[Token] {
	return NewSessionRouter(NewToken, m...)
}

func NewSessionRouter[S Session](ns NewSession[S], m ...Middleware) *Router[S] {
	return &Router[S]{http: httprouter.New(), sessions: ns, middleware: m}
}

func (r *Router[S]) Path(name string, m ...Middleware) *Router[S] {
	return &Router[S]{
		http:       r.http,
		sessions:   r.sessions,
		path:       r.path + name,
		middleware: append(r.middleware, m...),
	}
}

func (r *Router[S]) Handle(name, method string, h Handler[S], m ...Middleware) *Router[S] {
	n := NewHTTPHandler(r.sessions, h)
	m = append(r.middleware, m...)
	for i := len(m) - 1; i >= 0; i-- {
		n = m[i](n)
	}

	r.http.Handler(method, r.path+name, n)

	return r
}

func (r *Router[S]) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	r.http.ServeHTTP(res, req)
}
