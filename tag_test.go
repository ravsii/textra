package textra_test

// import (
// 	"reflect"
// 	"testing"

// 	"github.com/ravsii/textra"
// )

// func TestByName(t *testing.T) {
// 	type Tester struct {
// 		WithName    bool `json:"with_name"`
// 		WithoutName bool
// 	}

// 	testCases := []struct {
// 		testName  string
// 		fieldName string
// 		tagName   string
// 		found     bool
// 		tag       textra.Tag
// 	}{
// 		{"with name", "WithName", "json", true, textra.Tag{"json", "with_name", nil}},
// 		{"with name not found", "WithName", "sql", false, textra.Tag{}},
// 		{"without name", "WithoutName", "anything", false, textra.Tag{}},
// 	}

// 	data := textra.Extract(Tester{})

// 	for _, testCase := range testCases {
// 		testCase := testCase

// 		tag, found := data[testCase.fieldName].ByName(testCase.tagName)

// 		if found != testCase.found {
// 			t.Errorf("%s: found: got %t want %t", testCase.testName, found, testCase.found)
// 		}

// 		if found == false {
// 			continue
// 		}

// 		if !reflect.DeepEqual(tag, testCase.tag) {
// 			t.Errorf("%s: equal: got %s want %s", testCase.testName, tag, testCase.tag)
// 		}
// 	}
// }

// func TestOmitEmpty(t *testing.T) {
// 	type Str struct {
// 		WithOmit    bool `json:"with_omit,omitempty"`
// 		WithoutOmit bool `json:"without_omit"`
// 	}

// 	tags := textra.Extract(Str{})

// 	// Should be found
// 	if tag, _ := tags["WithOmit"].ByName("json"); tag.OmitEmpty() == false {
// 		t.Errorf("omitempty: got %t want %t", tag.OmitEmpty(), true)
// 	}

// 	// Should not be found
// 	if tag, _ := tags["WithoutOmit"].ByName("json"); tag.OmitEmpty() == true {
// 		t.Errorf("omitempty: got %t want %t", tag.OmitEmpty(), false)
// 	}
// }

// func TestIgnored(t *testing.T) {
// 	type Str struct {
// 		ID bool `json:"id,omitempty" sql:"-, pk"`
// 	}

// 	tags := textra.Extract(Str{})

// 	field, ok := tags["ID"]
// 	if !ok {
// 		t.Errorf("Tag \"ID\" wasn't found in the struct")
// 	}

// 	// Should be ignored
// 	if tag, _ := field.ByName("sql"); tag.Ignored() == false {
// 		t.Errorf("ignored: got %t want %t", tag.Ignored(), true)
// 	}

// 	// Should not be ignored
// 	if tag, _ := field.ByName("json"); tag.Ignored() == true {
// 		t.Errorf("ignored: got %t want %t", tag.Ignored(), false)
// 	}
// }
