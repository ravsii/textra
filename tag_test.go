package textra_test

import (
	"reflect"
	"testing"

	"github.com/ravsii/textra"
)

func TestByName(t *testing.T) {
	type Tester struct {
		Named       struct{} `json:"named"`
		NamedSpaced struct{} `json:"named spaced"` //nolint: tagliatelle
		Empty       struct{}
	}

	testCases := []struct {
		testName  string
		fieldName string
		tagName   string
		tagFound  bool
		tag       textra.Tag
	}{
		{"named", "Named", "json", true, textra.Tag{"json", "named", nil}},
		{"named with space", "NamedSpaced", "json", true, textra.Tag{"json", "named spaced", nil}},
		{"empty", "Empty", "json", false, textra.Tag{}},
	}

	data := textra.Extract(Tester{})

	for _, testCase := range testCases {
		testCase := testCase

		field, found := data.Field(testCase.fieldName)

		if !found {
			t.Errorf("%s: field %s not found", testCase.testName, testCase.fieldName)
		}

		tag, found := field.Tags.ByName(testCase.tagName)

		if found != testCase.tagFound {
			t.Errorf("%s: found: got %t want %t", testCase.testName, found, testCase.tagFound)
		}

		if found == false {
			continue
		}

		if !reflect.DeepEqual(tag, testCase.tag) {
			t.Errorf("%s: equal: got %s want %s", testCase.testName, tag, testCase.tag)
		}
	}
}

func TestOmitEmpty(t *testing.T) {
	type Str struct {
		WithOmit    bool `json:"with_omit,omitempty"`
		WithoutOmit bool `json:"without_omit"`
	}

	tags := textra.Extract(Str{})

	omitField, _ := tags.Field("WithOmit")

	// Should be found, omitempty == true is the expected result
	if tag, _ := omitField.Tags.ByName("json"); !tag.OmitEmpty() {
		t.Errorf("omitempty: got %t want %t", tag.OmitEmpty(), true)
	}

	nonOmitField, _ := tags.Field("WithoutOmit")

	// Should not be found, omitempty == false is the expected result
	if tag, _ := nonOmitField.Tags.ByName("json"); tag.OmitEmpty() {
		t.Errorf("omitempty: got %t want %t", tag.OmitEmpty(), false)
	}
}

func TestIgnored(t *testing.T) {
	type Str struct {
		ID struct{} `json:"id,omitempty" sql:"-, pk"`
	}

	str := textra.Extract(Str{})

	field, ok := str.Field("ID")
	if !ok {
		t.Errorf("Tag \"ID\" wasn't found in the struct")
	}

	// Should be ignored
	if tag, _ := field.Tags.ByName("sql"); !tag.Ignored() {
		t.Errorf("ignored: got %t want %t", tag.Ignored(), true)
	}

	// Should not be ignored
	if tag, _ := field.Tags.ByName("json"); tag.Ignored() {
		t.Errorf("ignored: got %t want %t", tag.Ignored(), false)
	}
}

func TestTagString(t *testing.T) {
	type Tester struct {
		TagBig   struct{} `json:"tag_big,omitempty"`
		TagSmall struct{} `sql:"-,pk" json:"tag_small"`
	}

	testCases := []struct {
		testName  string
		fieldName string
		tagName   string
		want      string
	}{
		{"big", "TagBig", "json", "json:\"tag_big,omitempty\""},
		{"small sql", "TagSmall", "sql", "sql:\"-,pk\""},
		{"small json", "TagSmall", "json", "json:\"tag_small\""},
	}

	data := textra.Extract(Tester{})

	for _, testCase := range testCases {
		testCase := testCase

		field, found := data.Field(testCase.fieldName)

		if !found {
			t.Errorf("%s: field %s not found", testCase.testName, testCase.fieldName)
		}

		tag, found := field.Tags.ByName(testCase.tagName)

		if !found {
			t.Errorf("%s: tag %s not found", testCase.testName, testCase.tagName)
		}

		got := tag.String()
		if got != testCase.want {
			t.Errorf("%s: got %s want %s", testCase.testName, got, testCase.want)
		}
	}
}

func TestTagsString(t *testing.T) {
	type Tester struct {
		TagBig   struct{} `json:"tag_big,omitempty"`
		TagSmall struct{} `sql:"-,pk" json:"tag_small"`
	}

	testCases := []struct {
		testName  string
		fieldName string
		want      string
	}{
		{"big", "TagBig", "[json:\"tag_big,omitempty\"]"},
		{"small sql", "TagSmall", "[sql:\"-,pk\" json:\"tag_small\"]"},
	}

	data := textra.Extract(Tester{})

	for _, testCase := range testCases {
		testCase := testCase

		field, found := data.Field(testCase.fieldName)

		if !found {
			t.Errorf("%s: field %s not found", testCase.testName, testCase.fieldName)
		}

		got := field.Tags.String()
		if got != testCase.want {
			t.Errorf("%s: got %s want %s", testCase.testName, got, testCase.want)
		}
	}
}
