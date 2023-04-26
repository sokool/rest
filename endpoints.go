package rest

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Endpoint[S Session] func(*Request[S]) (any, error)

type Endpoints[S Session] struct {
	path       string
	http       *httprouter.Router
	sessions   Sessions[S]
	middleware []Middleware
}

func NewEndpoints[S Session](s Sessions[S], m ...Middleware) *Endpoints[S] {
	return &Endpoints[S]{http: httprouter.New(), sessions: s, middleware: m}
}

func (s *Endpoints[S]) Path(name string, m ...Middleware) *Endpoints[S] {
	return &Endpoints[S]{
		path:       s.path + name,
		http:       s.http,
		sessions:   s.sessions,
		middleware: append(s.middleware, m...),
	}
}

func (s *Endpoints[S]) Handle(name, method string, e Endpoint[S], m ...Middleware) *Endpoints[S] {
	m = append(s.middleware, m...)
	h := NewHandler[S](s.sessions, e).Apply(m...)
	s.http.Handler(method, s.path+name, h)

	return s
}

func (s *Endpoints[S]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.http.ServeHTTP(w, r)
}
