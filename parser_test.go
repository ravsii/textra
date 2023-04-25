package textra

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestParseTags(t *testing.T) {
	tests := []struct {
		name string
		tag  reflect.StructTag
		want Tags
	}{
		{
			name: "Test with single tag",
			tag:  `json:"name"`,
			want: []Tag{{Tag: "json", Value: "name"}},
		},
		{
			name: "Test with multiple tags",
			tag:  `json:"name,omitempty" xml:"Name,v1,v2" sql:"-"`,
			want: []Tag{
				{Tag: "json", Value: "name", Optional: []string{"omitempty"}},
				{Tag: "xml", Value: "Name", Optional: []string{"v1", "v2"}},
				{Tag: "sql", Value: "-"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseTags(tt.tag); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseTags() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseType(t *testing.T) {
	tests := []struct {
		name string
		typ  reflect.Type
		want string
	}{
		{
			name: "Simple type",
			typ:  reflect.TypeOf(42),
			want: "int",
		},
		{
			name: "Struct type",
			typ:  reflect.TypeOf(struct{}{}),
			want: "struct",
		},
		{
			name: "Ptr type",
			typ:  reflect.TypeOf(&struct{}{}),
			want: "*struct",
		},
		{
			name: "Slice type",
			typ:  reflect.TypeOf([]bool{}),
			want: "[]bool",
		},
		{
			name: "Complex slice type",
			typ:  reflect.TypeOf([]*[]*string{}),
			want: "[]*[]*string",
		},
		{
			name: "Map type",
			typ:  reflect.TypeOf(map[string]int{}),
			want: "map[string]int",
		},
		{
			name: "Complex map type",
			typ:  reflect.TypeOf(map[*string]map[*bool]*reflect.Type{}),
			want: "map[*string]map[*bool]*reflect.Type",
		},
		{
			name: "Func type",
			typ:  reflect.TypeOf(func() {}),
			want: "func()",
		},
		{
			name: "Interface type",
			typ:  reflect.TypeOf((*interface{})(nil)),
			want: "interface",
		},
		{
			name: "Error type",
			typ:  reflect.TypeOf(errors.New("test")),
			want: "error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseType(tt.typ); !strings.Contains(got, tt.want) {
				t.Errorf("parseType() = %v, want %v", got, tt.want)
			}
		})
	}
}
