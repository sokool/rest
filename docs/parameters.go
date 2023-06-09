package docs

import (
	"fmt"
)

type Parameters []*Parameter

func NewParameters() *Parameters {
	return new(Parameters)
}

func (p *Parameters) Query(name, description string, value ...any) *Parameter {
	n := NewParameter("query", name, description, value)
	*p = append(*p, n)
	return n
}

func (p *Parameters) Path(name string, value any) *Parameter {
	n := NewParameter("path", name, "", value)
	*p = append(*p, n)
	return n
}

func (p *Parameters) Header(name string, value any) *Parameter {
	n := NewParameter("header", name, "", value)
	*p = append(*p, n)
	return n
}

func (p *Parameters) Cookie(name string, value any) *Parameter {
	n := NewParameter("cookie", name, "", value)
	*p = append(*p, n)
	return n
}

func (p *Parameters) Render(n int) (string, error) {
	if len(*p) == 0 {
		return "", nil
	}
	s := "parameters:"
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

func NewParameter(in, name, description string, value any) *Parameter {
	return &Parameter{
		in:          in,
		name:        name,
		description: description,
		required:    true,
		schema:      NewSchema(value),
	}
}

func (p *Parameter) Description(s string) *Parameter {
	p.description = s
	return p
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
  description: %s
  schema:
%s`,
		p.in,
		p.name,
		p.required,
		p.description,
		s,
	)

	return indent(n, x), nil
}
