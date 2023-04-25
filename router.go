package rest

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Router[S Session] struct {
	path       string
	http       *httprouter.Router
	sessions   Sessions[S]
	middleware []Middleware
}

func NewRouter[S Session](s Sessions[S], m ...Middleware) *Router[S] {
	return &Router[S]{http: httprouter.New(), sessions: s, middleware: m}
}

func (r *Router[S]) Path(name string, m ...Middleware) *Router[S] {
	return &Router[S]{
		path:       r.path + name,
		http:       r.http,
		sessions:   r.sessions,
		middleware: append(r.middleware, m...),
	}
}

func (r *Router[S]) Handle(name, method string, e Endpoint[S], m ...Middleware) *Router[S] {
	m = append(r.middleware, m...)
	h := NewHandler[S](r.sessions, e).Apply(m...)
	r.http.Handler(method, r.path+name, h)

	return r
}

func (r *Router[S]) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	r.http.ServeHTTP(res, req)
}
