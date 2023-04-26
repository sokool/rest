package rest_test

import (
	"net/http/httptest"
	"strings"
	"testing"

	. "github.com/sokool/rest"
)

func TestJSON(t *testing.T) {
	type scenario struct {
		description string
		endpoint    Endpoint[Token]

		code int
		body string
	}

	cases := []scenario{
		{
			"without payload and errors",
			func(r *Request[Token]) (any, error) { return nil, nil },
			204,
			``,
		},
		{
			"with payload and no errors",
			func(r *Request[Token]) (any, error) { return "hello", nil },
			200,
			`"hello"`,
		},
		{
			"with error and no payload",
			func(r *Request[Token]) (any, error) { return nil, Err("endpoint#test: %s failed", "request") },
			400,
			`{"message":"request failed","name":"endpoint","code":"test"}`,
		},
	}

	tt := NewTokens()
	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			r := httptest.NewRequest("GET", "http://none", nil)
			w := httptest.NewRecorder()
			NewHandler[Token](tt, c.endpoint).Apply(JSON).ServeHTTP(w, r)

			if s := w.Header().Get("Content-Type"); s != "application/json" {
				t.Fatalf("expected application/json content type in response header, got %s", s)
			}
			if n := w.Code; n != c.code {
				t.Fatalf("expected %d response code, got %d", c.code, n)
			}
			if s := w.Body.String(); strings.ReplaceAll(s, "\n", "") != c.body {
				t.Fatalf("expected %s response body, got %s", c.body, s)
			}
		})
	}
}
