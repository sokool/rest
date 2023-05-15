package docs

import (
	"fmt"
	"strings"
)

type OpenAPI struct {
	title, description string
	endpoints
}

func NewOpenAPI(title string) *OpenAPI {
	return &OpenAPI{
		title:     title,
		endpoints: make(endpoints),
	}
}

func (o *OpenAPI) Path(method string, name string) {
	o.endpoints.add(endpoint{method, name, ""})
}

func (o *OpenAPI) String() string {
	s := fmt.Sprintf(`
openapi: 3.0.0
info:
  title: "%s"
  description: "%s"
  version: "%s"
servers:
  - url: http://api.example.com/v1
    description: Optional server description, e.g. Main (production) server
  - url: http://staging-api.example.com
    description: Optional server description, e.g. Internal staging server for testing

`, o.title, o.description, "1.0")
	s += o.endpoints.String()
	return s
}

type endpoints map[string][]endpoint

func (c endpoints) add(e endpoint) {
	c[e.path] = append(c[e.path], e)
}

func (c endpoints) String() string {
	var s string
	s += "paths:\n"
	for path, ee := range c {
		s += fmt.Sprintf("  %s:\n", path)
		for _, e := range ee {
			s += e.String()
		}
	}
	return s
}

type endpoint struct {
	method, path, summary string
}

func (e endpoint) String() string {
	s := fmt.Sprintf("    %s:\n", strings.ToLower(e.method))
	if e.summary != "" {
		s += fmt.Sprintf("      summary: %s\n", e.summary)
	}
	s += fmt.Sprintf("      responses:\n")
	s += fmt.Sprintf("        200:\n")
	s += fmt.Sprintf("          description: hello word\n")
	return s
}
