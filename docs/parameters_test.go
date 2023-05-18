package docs_test

import (
	"fmt"
	"testing"

	"github.com/sokool/rest/docs"
)

func TestParameter_String(t *testing.T) {
	type Filter struct {
		ID      int    `json:"id"`
		Name    string `json:"name" jsonschema:"title=the name,description=The name of a friend,example=joe,example=lucy,default=alex"`
		Friends []int  `json:"friends,omitempty" jsonschema_description:"The list of IDs, omitted when empty"`
		//BirthDate   time.Time   `json:"birth_date,omitempty" jsonschema:"oneof_required=date"`
		//YearOfBirth string      `json:"year_of_birth,omitempty" jsonschema:"oneof_required=year"`
		//Metadata    interface{} `json:"metadata,omitempty" jsonschema:"oneof_type=string;array"`
		//FavColor    string      `json:"fav_color,omitempty" jsonschema:"enum=red,enum=green,enum=blue"`
	}

	p := docs.
		NewParameters().
		Query("q", "some tricky text").
		Cookie("age", 41).
		Path("filter", true).
		Header("X-Key", "Bearer abc")
	fmt.Println(p.Render(6))
	//fmt.Println(docs.NewPathParameter("filter", Filter{}).Render())
	//fmt.Println(docs.NewPathParameter("product", "adf39875fa").Render())
	//fmt.Println(docs.NewHeaderParameter("age", 38).Render())

}
