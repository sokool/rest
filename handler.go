package rest

import (
	"net/http"

	"github.com/sokool/rest/docs"
)

type Handler[S Session] struct {
	endpoint Endpoint[S]
	sessions Sessions[S]
	docs     *docs.Path
}

func NewHandler[S Session](s Sessions[S], e Endpoint[S]) *Handler[S] {
	return &Handler[S]{
		endpoint: e,
		sessions: s,
	}
}

func (h *Handler[S]) ServeHTTP(hw http.ResponseWriter, hr *http.Request) {
	var r *Request[S]
	var p any
	var err error

	if r, err = NewRequest(h.sessions, hw, hr); err != nil {
		Write[error](responseErr, err, r.Request)
		return
	}
	r.docs = h.docs
	if p, err = h.endpoint(r); err != nil {
		Write(responseErr, err, r.Request)
	}
	if p != nil {
		Write(responseBody, p, r.Request)
	}
	if h.docs != nil {
		h.docs.Response(p)
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

func (h *Handler[S]) Doc(p *docs.Path) *Handler[S] {
	h.docs = p
	return h
}

const (
	responseErr  = "error"
	responseBody = "response"
)
