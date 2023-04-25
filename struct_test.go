package textra_test

import (
	"reflect"
	"testing"

	"github.com/ravsii/textra"
)

func TestGetField(t *testing.T) {
	type Tester struct {
		Tag struct{} `json:"tag"`
	}

	data := textra.Extract(Tester{})

	if _, ok := data.Field("Tag"); !ok {
		t.Errorf("Field: %s should be found", "Tag")
	}

	if _, ok := data.Field("NotExists"); ok {
		t.Errorf("Field: %s should NOT be found", "NotExists")
	}
}

func TestByTagName(t *testing.T) {
	type Tester struct {
		Tag1 struct{} `json:"tag1"`
		Tag2 struct{} `json:"tag2" sql:"tag2, pk"`
	}

	testCases := []struct {
		name    string
		tagName string
		want    textra.Struct
	}{
		{"json", "json", textra.Struct{
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
		}},
		{"sql", "sql", textra.Struct{
			textra.Field{
				Name: "Tag2",
				Type: "struct",
				Tags: textra.Tags{
					{"json", "tag2", nil},
					{"sql", "tag2", []string{"pk"}},
				},
			},
		}},
		{"non-existent", "nonexistent", nil},
	}

	data := textra.Extract(Tester{})

	for _, testCase := range testCases {
		testCase := testCase

		got := data.ByTagName(testCase.tagName)

		// nil / empty checks
		if len(got) == len(testCase.want) && reflect.DeepEqual(got, textra.Struct{}) {
			continue
		}

		if !reflect.DeepEqual(got, testCase.want) {
			t.Errorf("%s: got %v want %v", testCase.name, got, testCase.want)
		}
	}
}

func TestByAnyTagName(t *testing.T) {
	type Tester struct {
		Tag1 struct{} `json:"tag1"`
		Tag2 struct{} `json:"tag2" sql:"tag2, pk"`
	}

	testCases := []struct {
		name     string
		tagNames []string
		want     textra.Struct
	}{
		{"json & sql", []string{"json", "sql"}, textra.Struct{
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
		}},
		{"sql only", []string{"sql"}, textra.Struct{
			textra.Field{
				Name: "Tag2",
				Type: "struct",
				Tags: textra.Tags{
					{"json", "tag2", nil},
					{"sql", "tag2", []string{"pk"}},
				},
			},
		}},
		{"non-existent", []string{}, nil},
	}

	data := textra.Extract(Tester{})

	for _, testCase := range testCases {
		testCase := testCase

		got := data.ByAnyTagName(testCase.tagNames...)

		// nil / empty checks
		if len(got) == len(testCase.want) && reflect.DeepEqual(got, textra.Struct{}) {
			continue
		}

		if !reflect.DeepEqual(got, testCase.want) {
			t.Errorf("%s: got %v want %v", testCase.name, got, testCase.want)
		}
	}
}

func TestByTagNames(t *testing.T) {
	type Tester struct {
		Tag1 struct{} `json:"tag1"`
		Tag2 struct{} `json:"tag2" sql:"tag2, pk"`
	}

	testCases := []struct {
		name     string
		tagNames []string
		want     textra.Struct
	}{
		{"json & sql", []string{"json", "sql"}, textra.Struct{
			textra.Field{
				Name: "Tag2",
				Type: "struct",
				Tags: textra.Tags{
					{"json", "tag2", nil},
					{"sql", "tag2", []string{"pk"}},
				},
			},
		}},
		{"sql only", []string{"sql"}, textra.Struct{
			textra.Field{
				Name: "Tag2",
				Type: "struct",
				Tags: textra.Tags{
					{"json", "tag2", nil},
					{"sql", "tag2", []string{"pk"}},
				},
			},
		}},
		{"non-existent", []string{}, nil},
	}

	data := textra.Extract(Tester{})

	for _, testCase := range testCases {
		testCase := testCase

		got := data.ByTagNames(testCase.tagNames...)

		// nil / empty checks
		if len(got) == len(testCase.want) && reflect.DeepEqual(got, textra.Struct{}) {
			continue
		}

		if !reflect.DeepEqual(got, testCase.want) {
			t.Errorf("%s: got %v want %v", testCase.name, got, testCase.want)
		}
	}
}

func TestFilterFunc(t *testing.T) {
	type Tester struct {
		Tag1 struct{} `json:"tag1"`
		Tag2 struct{} `json:"tag2" sql:"tag2, pk"`
	}

	testCases := []struct {
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

	for _, testCase := range testCases {
		testCase := testCase

		got := data.FilterFunc(testCase.filterFunc)

		// nil / empty checks
		if len(got) == len(testCase.want) && reflect.DeepEqual(got, textra.Struct{}) {
			continue
		}

		if !reflect.DeepEqual(got, testCase.want) {
			t.Errorf("%s: got %v want %v", testCase.name, got, testCase.want)
		}
	}
}

func TestRemoveEmpty(t *testing.T) {
	type Tester struct {
		Tag1 struct{} `json:"tag1"`
	}

	type TesterWithEmpty struct {
		Empty struct{}
		Tag1  struct{} `json:"tag1"`
	}

	testCases := []struct {
		name string
		str  interface{}
		want textra.Struct
	}{
		{"with empty", TesterWithEmpty{}, textra.Struct{
			textra.Field{
				Name: "Tag1",
				Type: "struct",
				Tags: textra.Tags{
					{"json", "tag1", nil},
				},
			},
		}},
		{"without empty", Tester{}, textra.Struct{
			textra.Field{
				Name: "Tag1",
				Type: "struct",
				Tags: textra.Tags{
					{"json", "tag1", nil},
				},
			},
		}},
		{"non-existent", struct{}{}, nil},
	}

	for _, testCase := range testCases {
		testCase := testCase

		got := textra.Extract(testCase.str).RemoveEmpty()

		// nil / empty checks
		if len(got) == len(testCase.want) && reflect.DeepEqual(got, textra.Struct{}) {
			continue
		}

		if !reflect.DeepEqual(got, testCase.want) {
			t.Errorf("%s: got %v want %v", testCase.name, got, testCase.want)
		}
	}
}
