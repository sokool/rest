package rest

import (
	"net/http"
)

type Endpoint[S Session] func(*Request[S]) (any, error)

type Handler[S Session] struct {
	endpoint Endpoint[S]
	sessions Sessions[S]
}

func NewHandler[S Session](s Sessions[S], e Endpoint[S]) *Handler[S] {
	return &Handler[S]{
		endpoint: e,
		sessions: s,
	}
}

func (h *Handler[S]) ServeHTTP(hw http.ResponseWriter, hr *http.Request) {
	var r = &Request[S]{Request: hr, Response: hw}
	var err error
	if r.Session, err = h.sessions.Read(hr); err != nil {
		Write[error](responseErr, err, r.Request)
		return
	}
	p, err := h.endpoint(r)
	if err != nil {
		Write(responseErr, err, r.Request)
	}
	if p != nil {
		Write(responseBody, p, r.Request)
	}
	*hr = *r.Request
}

func (h *Handler[S]) Apply(m ...Middleware) http.Handler {
	var n = http.Handler(h)
	for i := len(m) - 1; i >= 0; i-- {
		n = m[i](n)
	}
	return n
}

const (
	responseErr  = "error"
	responseBody = "response"
)
