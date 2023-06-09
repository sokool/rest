package docs_test

import (
	"hash/fnv"
	"testing"

	"github.com/sokool/rest/docs"
)

func TestParameter_String(t *testing.T) {
	p := docs.NewParameters()
	p.Query("q", "some tricky text").Description("oh some tricky query param")
	p.Cookie("age", 41)
	p.Path("filter", true)
	p.Header("X-Key", "Bearer abc")
	s, err := p.Render(0)
	if err != nil {
		t.Fatal(err)
	}

	if hash(s) != 1772388703 {
		t.Fatalf("invalid string output")
	}

}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
