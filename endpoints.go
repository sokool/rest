package rest

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/sokool/rest/docs"
)

type Endpoint[S Session] func(*Request[S]) (any, error)

type Endpoints[S Session] struct {
	path       string
	http       *httprouter.Router
	sessions   Sessions[S]
	middleware []Middleware

	docs *docs.OpenAPI
}

func NewEndpoints[S Session](s Sessions[S], m ...Middleware) *Endpoints[S] {
	return &Endpoints[S]{
		http:       httprouter.New(),
		sessions:   s,
		middleware: m,
		docs:       docs.NewOpenAPI("example"),
	}
}

func (s *Endpoints[S]) Path(name string, m ...Middleware) *Endpoints[S] {
	return &Endpoints[S]{
		path:       s.path + name,
		http:       s.http,
		sessions:   s.sessions,
		docs:       s.docs,
		middleware: append(s.middleware, m...),
	}
}

func (s *Endpoints[S]) Handle(n, method string, e Endpoint[S], m ...Middleware) *Endpoints[S] {
	m = append(s.middleware, m...)
	n = s.path + n
	d := s.docs.Path(method, n)
	h := NewHandler[S](s.sessions, e).
		Doc(d).
		Apply(m...)

	s.http.Handler(method, n, h)
	return s
}

func (s *Endpoints[S]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.http.ServeHTTP(w, r)
}

func (s *Endpoints[S]) String() string {
	return s.docs.String()
}
