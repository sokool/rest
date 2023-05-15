package rest_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/sokool/rest"
)

func TestRouter_MiddlewareOrder(t *testing.T) {
	type _case struct {
		description string
		register    *Endpoints[Token]
		method      string
		path        string
		headers     string
	}

	m := func(s string) Middleware {
		return func(next http.Handler) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("string", w.Header().Get("string")+" "+s+":before")
				next.ServeHTTP(w, r)
				w.Header().Set("string", w.Header().Get("string")+" "+s+":after")
			}
		}
	}
	h := func(s string) Endpoint[Token] {
		return func(r *Request[Token]) (any, error) {
			r.Response.Header().Set("string", r.Response.Header().Get("string")+" "+s)
			return s, nil
		}
	}

	r := NewEndpoints[Token](&Tokens{}, JSON)

	cases := []_case{
		{
			"without",
			r.Path("/a").Handle("", "GET", h("hello")),
			"GET", "/a?test=abc",
			"hello",
		},
		{
			"one",
			r.Path("/b", m("a#1")).Handle("", "GET", h("orange")),
			"GET", "/b",
			"a#1:before orange a#1:after",
		},
		{
			"one on one path",
			r.Path("/c", m("a#1")).Handle("", "GET", h("orange"), m("b#1")),
			"GET", "/c",
			"a#1:before b#1:before orange b#1:after a#1:after",
		},
		{
			"multiple on different paths",
			r.Path("/apps", m("apps#1")).Path("/tools", m("tool#1"), m("tool#2")).Handle("", "GET", h("handler"), m("get#1")),
			"GET", "/apps/tools",
			"apps#1:before tool#1:before tool#2:before get#1:before handler get#1:after tool#2:after tool#1:after apps#1:after",
		},
	}

	s := httptest.NewServer(r)
	defer s.Close()
	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			res, _ := http.Get(s.URL + c.path)
			if h := res.Header.Get("string"); h != c.headers {
				t.Fatalf("\nexpected `%s`\n     got `%s`", c.headers, h)
			}
		})
	}

	fmt.Println(r)
}
