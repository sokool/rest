package docs

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

type OpenAPI struct {
	title, description string
	endpoints
}

func (o *OpenAPI) Value() (driver.Value, error) {
	//TODO implement me
	panic("implement me")
}

func (o *OpenAPI) Scan(src any) error {
	//TODO implement me
	panic("implement me")
}

func NewOpenAPI(title string) *OpenAPI {
	return &OpenAPI{
		title:     title,
		endpoints: make(endpoints),
	}
}

func (o *OpenAPI) Path(method string, pattern string) *Path {
	e := &Path{method, pattern, "", NewParameters()}
	o.endpoints.add(e)
	return e
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

type endpoints map[string][]*Path

func (c endpoints) add(e *Path) {
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

type Path struct {
	method, path, summary string
	parameters            *Parameters
}

//	func (e Path) param() ([]string, bool) {
//		return e.bar('{', '}')
//	}
func (e *Path) param(start, end rune) []string {
	var z []string
	var k int
	var m bool
	for _, c := range e.path {
		if c == end && m {
			k, m = k+1, false
		}
		if m {
			z[k] = z[k] + string(c)
		}
		if c == start {
			m, z = true, append(z, "")
		}
	}
	return z
}

func (e *Path) Parameters() *Parameters {
	return e.parameters
}

func (e *Path) Response(body any) {
	//fmt.Println(body)
}

func (e *Path) String() string {
	s := fmt.Sprintf("    %s:\n", strings.ToLower(e.method))
	if e.summary != "" {
		s += fmt.Sprintf("      summary: %s\n", e.summary)
	}

	x, err := e.parameters.Render(6)
	if err != nil {
		fmt.Println(err)
	}
	s += x
	//xx := e.param(':', '/')
	//if len(xx) != 0 {
	//	s += fmt.Sprintf("      parameters:\n")
	//}
	//for _, p := range xx {
	//	s += fmt.Sprintf("        - in: path\n")
	//	s += fmt.Sprintf("          name: %s\n", p)
	//	s += fmt.Sprintf("          required: true\n")
	//	s += fmt.Sprintf("          schema:\n")
	//	s += fmt.Sprintf("            type: string\n")
	//}
	s += fmt.Sprintf("      responses:\n")
	s += fmt.Sprintf("        200:\n")
	s += fmt.Sprintf("          description: \"undefined\"\n")
	return s
}

type Paths []*Path
