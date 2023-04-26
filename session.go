package rest

import (
	"net/http"
)

type Sessions[S Session] interface {
	Read(r *http.Request) (S, error)
}
type Session interface {
	IsExpired() bool
}

type Tokens struct{}

func NewTokens() *Tokens {
	return &Tokens{}
}

func (t *Tokens) Read(r *http.Request) (Token, error) {
	s := r.Header.Get("Authorization")
	return NewToken(s)
}

type Token string

func NewToken(s string) (Token, error) {
	t := Token(s)
	return t, nil
}

func (s Token) IsExpired() bool {
	return false
}
