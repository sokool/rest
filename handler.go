package rest

import (
	"net/http"
)

type Handler[S Session] func(*Request[S]) (any, error)

func (h Handler[S]) Do(r *Request[S]) (any, error) {
	return h(r)
}

func NewHTTPHandler[S Session](s NewSession[S], h Handler[S]) http.Handler {
	return http.HandlerFunc(func(hw http.ResponseWriter, hr *http.Request) {
		r, err := NewRequest[S](s, hw, hr)
		if err != nil {
			Write(responseErr, err, r.Request)
			return
		}
		p, err := h.Do(r)
		if err != nil {
			Write(responseErr, err, r.Request)
		}
		if p != nil {
			Write(responseBody, p, r.Request)
		}
		*hr = *r.Request
	})
}

const (
	responseErr  = "error"
	responseBody = "response"
)
