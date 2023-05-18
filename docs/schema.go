package docs

import (
	"encoding/json"
	"strings"

	"github.com/invopop/jsonschema"
	"gopkg.in/yaml.v3"
)

type Schema struct {
	json *jsonschema.Schema
}

func NewSchema(value any) *Schema {
	return &Schema{jsonschema.Reflect(value)}
}

func (s *Schema) Render(n int) (string, error) {
	var v any = s.json
	for n := range s.json.Definitions {
		x := s.json.Definitions[n]
		x.AdditionalProperties = nil
		v = x
		break
	}

	j, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	m := make(map[string]any)
	if err = json.Unmarshal(j, &m); err != nil {
		return "", err
	}
	delete(m, "$schema")
	delete(m, "$id")

	b, err := yaml.Marshal(m)
	if err != nil {
		return "", err
	}
	a := indent(n, string(b))
	return a, nil
}

func indent(n int, s string) string {
	var x string
	for i := 0; i < n; i++ {
		x += " "
	}

	if x == "" || s == "" {
		return s
	}
	lines := strings.SplitAfter(s, "\n")
	if len(lines[len(lines)-1]) == 0 {
		lines = lines[:len(lines)-1]
	}
	return strings.Join(append([]string{""}, lines...), x)
}
