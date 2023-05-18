package docs

import (
	"fmt"
)

type Parameters []*Parameter

func NewParameters() *Parameters {
	return new(Parameters)
}

func (p *Parameters) Query(name string, value any) *Parameters {
	n := NewParameter("query", name, value)
	*p = append(*p, n)
	return p
}

func (p *Parameters) Path(name string, value any) *Parameters {
	n := NewParameter("path", name, value)
	*p = append(*p, n)
	return p
}

func (p *Parameters) Header(name string, value any) *Parameters {
	n := NewParameter("header", name, value)
	*p = append(*p, n)
	return p
}

func (p *Parameters) Cookie(name string, value any) *Parameters {
	n := NewParameter("cookie", name, value)
	*p = append(*p, n)
	return p
}

func (p *Parameters) Render(n int) (string, error) {
	var s string
	if len(*p) != 0 {
		s = "parameters:"
	}
	for i := range *p {
		ps, err := (*p)[i].Render(2)
		if err != nil {
			return "", err
		}
		s += ps
	}
	return indent(n, s), nil
}

type Parameter struct {
	in          string
	name        string
	description string
	schema      *Schema
	required    bool
}

func NewParameter(in, name string, value any) *Parameter {
	return &Parameter{
		in:       in,
		name:     name,
		required: true,
		schema:   NewSchema(value),
	}
}

func (p *Parameter) Render(n int) (string, error) {
	s, err := p.schema.Render(4)
	if err != nil {
		return "", err
	}

	x := fmt.Sprintf(`
- in: %s
  name: %s
  required: %v
  schema:
%s`,
		p.in,
		p.name,
		p.required,
		s,
	)

	return indent(n, x), nil
}
