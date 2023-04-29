package textra_test

import (
	"testing"

	"github.com/ravsii/textra"
)

func TestField_String(t *testing.T) {
	field := textra.Field{
		Name: "field",
		Type: "string",
		Tags: textra.Tags{
			{Tag: "json", Value: "name"},
			{Tag: "xml", Optional: []string{"pk"}},
		},
	}
	want := `field(string):[json:"name" xml:",pk"]`
	if got := field.String(); got != want {
		t.Errorf("Field.String() = %q, expected %q", got, want)
	}
}

func TestFieldTag_String(t *testing.T) {
	fieldTag := textra.FieldTag{
		Name: "field",
		Type: "*[]int",
		Tag:  textra.Tag{Tag: "sql", Value: "-", Optional: []string{"omitempty"}},
	}
	want := `field(*[]int):sql:"-,omitempty"`
	if got := fieldTag.String(); got != want {
		t.Errorf("FieldTag.String() = %q, expected %q", got, want)
	}
}
