package textra_test

import (
	"reflect"
	"testing"

	"github.com/ravsii/textra"
)

func TestGetField(t *testing.T) {
	type TestOne struct {
		Tag struct{} `json:"tag"`
	}

	data := textra.Extract((*TestOne)(nil))
	if _, ok := data.Field("Tag"); !ok {
		t.Errorf("Field() =  %t, want: found", ok)
	}
	if _, ok := data.Field("NotExists"); ok {
		t.Errorf("Field() =  %t, want: not found", ok)
	}
}

func TestByTagName(t *testing.T) {
	type TestMultiple struct {
		Tag1 struct{} `json:"tag1"`
		Tag2 struct{} `json:"tag2" pg:"tag2" sql:"tag2, pk"`
		Tag3 struct{} `json:"tag3" sql:"tag3, pk"`
		Tag4 struct{} `json:"tag4" gorm:",pk" sql:"tag4, pk"`
	}

	tests := []struct {
		tagName string
		want    textra.Struct
	}{
		{"pg", textra.Struct{
			{
				Name: "Tag2",
				Type: "struct",
				Tags: textra.Tags{
					{"json", "tag2", nil},
					{"pg", "tag2", nil},
					{"sql", "tag2", []string{"pk"}},
				},
			},
		}},
		{"sql", textra.Struct{
			{
				Name: "Tag2",
				Type: "struct",
				Tags: textra.Tags{
					{"json", "tag2", nil},
					{"pg", "tag2", nil},
					{"sql", "tag2", []string{"pk"}},
				},
			},
			{
				Name: "Tag3",
				Type: "struct",
				Tags: textra.Tags{
					{"json", "tag3", nil},
					{"sql", "tag3", []string{"pk"}},
				},
			},
			{
				Name: "Tag4",
				Type: "struct",
				Tags: textra.Tags{
					{"json", "tag4", nil},
					{"gorm", "", []string{"pk"}},
					{"sql", "tag4", []string{"pk"}},
				},
			},
		}},
		{"nonexistent", nil},
	}

	data := textra.Extract((*TestMultiple)(nil))

	for _, tt := range tests {
		tt := tt
		t.Run(tt.tagName, func(t *testing.T) {
			t.Parallel()
			got := data.ByTagName(tt.tagName)
			if !checkEqual(t, got, tt.want) {
				t.Errorf("TestByTagName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestByTagNameAny(t *testing.T) {
	type TestMultiple struct {
		Tag1 struct{} `json:"tag1"`
		Tag2 struct{} `json:"tag2" pg:"tag2" sql:"tag2, pk"`
		Tag3 struct{} `json:"tag3" sql:"tag3, pk"`
		Tag4 struct{} `json:"tag4" gorm:",pk" sql:"tag4, pk"`
	}

	tests := []struct {
		name     string
		tagNames []string
		want     textra.Struct
	}{
		{"pg & gorm", []string{"pg", "gorm"}, textra.Struct{
			{
				Name: "Tag2",
				Type: "struct",
				Tags: textra.Tags{
					{"json", "tag2", nil},
					{"pg", "tag2", nil},
					{"sql", "tag2", []string{"pk"}},
				},
			},
			{
				Name: "Tag4",
				Type: "struct",
				Tags: textra.Tags{
					{"json", "tag4", nil},
					{"gorm", "", []string{"pk"}},
					{"sql", "tag4", []string{"pk"}},
				},
			},
		}},
		{"sql only", []string{"sql"}, textra.Struct{
			{
				Name: "Tag2",
				Type: "struct",
				Tags: textra.Tags{
					{"json", "tag2", nil},
					{"pg", "tag2", nil},
					{"sql", "tag2", []string{"pk"}},
				},
			},
			{
				Name: "Tag3",
				Type: "struct",
				Tags: textra.Tags{
					{"json", "tag3", nil},
					{"sql", "tag3", []string{"pk"}},
				},
			},
			{
				Name: "Tag4",
				Type: "struct",
				Tags: textra.Tags{
					{"json", "tag4", nil},
					{"gorm", "", []string{"pk"}},
					{"sql", "tag4", []string{"pk"}},
				},
			},
		}},
		{"non-existent", []string{}, nil},
	}

	data := textra.Extract((*TestMultiple)(nil))

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := data.ByTagNameAny(tt.tagNames...)
			if !checkEqual(t, got, tt.want) {
				t.Errorf("TestByTagNameAny() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestByTagNameAll(t *testing.T) {
	type TestMultiple struct {
		Tag1 struct{} `json:"tag1"`
		Tag2 struct{} `json:"tag2" pg:"tag2" sql:"tag2, pk"`
		Tag3 struct{} `json:"tag3" sql:"tag3, pk"`
		Tag4 struct{} `json:"tag4" gorm:",pk" sql:"tag4, pk"`
	}

	tests := []struct {
		name     string
		tagNames []string
		want     textra.Struct
	}{
		{"json & pg & sql", []string{"json", "pg", "sql"},
			textra.Struct{
				{
					Name: "Tag2",
					Type: "struct",
					Tags: textra.Tags{
						{"json", "tag2", nil},
						{"pg", "tag2", nil},
						{"sql", "tag2", []string{"pk"}},
					},
				},
			},
		},
		{"gorm only", []string{"gorm"},
			textra.Struct{
				{
					Name: "Tag4",
					Type: "struct",
					Tags: textra.Tags{
						{"json", "tag4", nil},
						{"gorm", "", []string{"pk"}},
						{"sql", "tag4", []string{"pk"}},
					},
				},
			}},
		{"non-existent", []string{}, nil},
	}

	data := textra.Extract((*TestMultiple)(nil))

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := data.ByTagNameAll(tt.tagNames...)
			if !checkEqual(t, got, tt.want) {
				t.Errorf("TestByTagNameAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilterFunc(t *testing.T) {
	type Tester struct {
		Tag1 struct{} `json:"tag1"`
		Tag2 struct{} `json:"tag2" sql:"tag2, pk"`
	}

	tests := []struct {
		name       string
		filterFunc func(textra.Field) bool
		want       textra.Struct
	}{
		{
			"all fields", func(f textra.Field) bool { return true }, textra.Struct{
				textra.Field{
					Name: "Tag1",
					Type: "struct",
					Tags: textra.Tags{
						{"json", "tag1", nil},
					},
				},
				textra.Field{
					Name: "Tag2",
					Type: "struct",
					Tags: textra.Tags{
						{"json", "tag2", nil},
						{"sql", "tag2", []string{"pk"}},
					},
				},
			},
		},
		{"non-existent", func(f textra.Field) bool { return false }, nil},
	}

	data := textra.Extract(Tester{})

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := data.FilterFunc(tt.filterFunc)
			if !checkEqual(t, got, tt.want) {
				t.Errorf("TestFilterFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemoveEmpty(t *testing.T) {
	type Tester struct {
		Tag1 struct{} `json:"tag1"`
	}

	type TesterWithEmpty struct {
		Empty  struct{}
		Empty2 struct{}
		Empty3 struct{}
		Tag1   struct{} `json:"tag1"`
		Empty4 struct{}
		Empty5 struct{}
		Tag2   struct{} `json:"tag2"`
	}

	tests := []struct {
		name string
		str  interface{}
		want textra.Struct
	}{
		{"with empty", (*TesterWithEmpty)(nil), textra.Struct{
			{
				Name: "Tag1",
				Type: "struct",
				Tags: textra.Tags{
					{"json", "tag1", nil},
				},
			},
			{
				Name: "Tag2",
				Type: "struct",
				Tags: textra.Tags{
					{"json", "tag2", nil},
				},
			},
		}},
		{"without empty", (*Tester)(nil), textra.Struct{
			{
				Name: "Tag1",
				Type: "struct",
				Tags: textra.Tags{
					{"json", "tag1", nil},
				},
			},
		}},
		{"non-existent", struct{}{}, nil},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := textra.Extract(tt.str).RemoveEmpty()
			if !checkEqual(t, got, tt.want) {
				t.Errorf("TestRemoveEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOnlyTag(t *testing.T) {
	type empty struct{}
	type onlyOne struct {
		Example string `json:"ex1"`
	}
	type many struct {
		Str1 string `json:"str1" pg:"str1" sql:"str1"`
		Str2 string `json:"str2" pg:"str2"`
	}

	tests := []struct {
		name string
		str  interface{}
		only string
		want []textra.FieldTag
	}{
		{"empty", (*empty)(nil), "test", []textra.FieldTag{}},
		{"only one", (*onlyOne)(nil), "json",
			[]textra.FieldTag{
				{
					Name: "Example",
					Type: "string",
					Tag:  textra.Tag{Tag: "json", Value: "ex1"},
				},
			},
		},
		{"many json", (*many)(nil), "json",
			[]textra.FieldTag{
				{
					Name: "Str1",
					Type: "string",
					Tag:  textra.Tag{Tag: "json", Value: "str1"},
				},
				{
					Name: "Str2",
					Type: "string",
					Tag:  textra.Tag{Tag: "json", Value: "str2"},
				},
			},
		},
		{"many sql", (*many)(nil), "pg",
			[]textra.FieldTag{
				{
					Name: "Str1",
					Type: "string",
					Tag:  textra.Tag{Tag: "pg", Value: "str1"},
				},
				{
					Name: "Str2",
					Type: "string",
					Tag:  textra.Tag{Tag: "pg", Value: "str2"},
				},
			},
		},
		{"many pg", (*many)(nil), "sql",
			[]textra.FieldTag{
				{
					Name: "Str1",
					Type: "string",
					Tag:  textra.Tag{Tag: "sql", Value: "str1"},
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			extracted := textra.Extract(tt.str)
			got := extracted.OnlyTag(tt.only)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OnlyTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemoveFields(t *testing.T) {
	type TestMultiple struct {
		Tag1 struct{} `json:"tag1"`
		Tag2 struct{} `json:"tag2" pg:"tag2" sql:"tag2, pk"`
		Tag3 struct{} `json:"tag3" sql:"tag3, pk"`
		Tag4 struct{} `json:"tag4" gorm:",pk" sql:"tag4, pk"`
	}

	tests := []struct {
		name   string
		fields []string
		want   textra.Struct
	}{
		{"empty", []string{},
			textra.Struct{
				{
					Name: "Tag1",
					Type: "struct",
					Tags: textra.Tags{
						{"json", "tag1", nil},
					},
				},
				{
					Name: "Tag2",
					Type: "struct",
					Tags: textra.Tags{
						{"json", "tag2", nil},
						{"pg", "tag2", nil},
						{"sql", "tag2", []string{"pk"}},
					},
				},
				{
					Name: "Tag3",
					Type: "struct",
					Tags: textra.Tags{
						{"json", "tag3", nil},
						{"sql", "tag3", []string{"pk"}},
					},
				},
				{
					Name: "Tag4",
					Type: "struct",
					Tags: textra.Tags{
						{"json", "tag4", nil},
						{"gorm", "", []string{"pk"}},
						{"sql", "tag4", []string{"pk"}},
					},
				},
			},
		},
		{"tag1", []string{"Tag1"},
			textra.Struct{
				{
					Name: "Tag2",
					Type: "struct",
					Tags: textra.Tags{
						{"json", "tag2", nil},
						{"pg", "tag2", nil},
						{"sql", "tag2", []string{"pk"}},
					},
				},
				{
					Name: "Tag3",
					Type: "struct",
					Tags: textra.Tags{
						{"json", "tag3", nil},
						{"sql", "tag3", []string{"pk"}},
					},
				},
				{
					Name: "Tag4",
					Type: "struct",
					Tags: textra.Tags{
						{"json", "tag4", nil},
						{"gorm", "", []string{"pk"}},
						{"sql", "tag4", []string{"pk"}},
					},
				},
			},
		},
		{"tag2 tag4", []string{"Tag2", "Tag4"},
			textra.Struct{
				{
					Name: "Tag1",
					Type: "struct",
					Tags: textra.Tags{
						{"json", "tag1", nil},
					},
				},
				{
					Name: "Tag3",
					Type: "struct",
					Tags: textra.Tags{
						{"json", "tag3", nil},
						{"sql", "tag3", []string{"pk"}},
					},
				},
			},
		},
		{"all but 1", []string{"Tag2", "Tag3", "Tag4"},
			textra.Struct{
				{
					Name: "Tag1",
					Type: "struct",
					Tags: textra.Tags{
						{"json", "tag1", nil},
					},
				},
			},
		},
	}

	data := textra.Extract((*TestMultiple)(nil))

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := data.RemoveFields(tt.fields...)
			if !checkEqual(t, got, tt.want) {
				t.Errorf("TestRemoveFields() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestString(t *testing.T) {
	type (
		TestEmpty struct{}
		TestOne   struct {
			Tag struct{} `json:"tag"`
		}
		TestMultiple struct {
			Tag1 struct{} `json:"tag1"`
			Tag2 struct{} `json:"tag2" pg:"tag2" sql:"tag2, pk"`
			Tag3 struct{} `json:"tag3" sql:"tag3, pk"`
			Tag4 struct{} `json:"tag4" gorm:",pk" sql:"tag4, pk"`
		}
	)

	tests := []struct {
		name string
		str  interface{}
		want string
	}{
		{"TestEmpty", (*TestEmpty)(nil), ``},
		{"TestOne", (*TestOne)(nil), `Tag(struct):[json:"tag"]`},
		{
			"TestMultiple",
			(*TestMultiple)(nil),
			`Tag1(struct):[json:"tag1"]Tag2(struct):[json:"tag2" pg:"tag2" sql:"tag2,pk"]Tag3(struct):[json:"tag3" sql:"tag3,pk"]Tag4(struct):[json:"tag4" gorm:",pk" sql:"tag4,pk"]`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := textra.Extract(tt.str).String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func checkEqual(t *testing.T, got, want textra.Struct) bool {
	// nil / empty checks
	if len(got) == len(want) && reflect.DeepEqual(got, textra.Struct{}) {
		return true
	}

	if !reflect.DeepEqual(got, want) {
		return false
	}

	return true
}
