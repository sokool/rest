package rest

import "net/http"

type Session interface {
	IsExpired() bool
}

type NewSession[S Session] func(http.ResponseWriter, *http.Request) (S, error)

type Token string

func NewToken(w http.ResponseWriter, r *http.Request) (Token, error) {
	return Token(r.Header.Get("Authorization")), nil
}

func (s Token) IsExpired() bool { return false }
