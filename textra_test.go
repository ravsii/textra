package textra_test

import (
	"reflect"
	"testing"

	"github.com/ravsii/textra"
)

func TestExtract(t *testing.T) {
	type Empty struct{}

	type TesterFilled struct {
		Tag1 struct{} `json:"tag1"`
		Tag2 struct{} `json:"tag2" sql:"tag2, pk"`
	}

	expectedFilled := textra.Struct{
		textra.Field{
			Name: "Tag1",
			Tags: textra.Tags{
				{"json", "tag1", nil},
			},
		},
		textra.Field{
			Name: "Tag2",
			Tags: textra.Tags{
				{"json", "tag2", nil},
				{"sql", "tag2", []string{"pk"}},
			},
		},
	}

	testCases := []struct {
		name  string
		input any
		want  textra.Struct
	}{
		{"Empty, non-pointer", Empty{}, nil},
		{"Empty, pointer", &Empty{}, nil},
		{"Empty, nil-check", (*Empty)(nil), nil},
		{"Populated, non-pointer", TesterFilled{}, expectedFilled},
		{"Populated, pointer", &TesterFilled{}, expectedFilled},
		{"Populated, nil-check", (*TesterFilled)(nil), expectedFilled},
	}

	for _, testCase := range testCases {
		testCase := testCase

		got := textra.Extract(testCase.input)

		// nil / empty checks
		if len(got) == len(testCase.want) && reflect.DeepEqual(got, textra.Struct{}) {
			continue
		}

		if !reflect.DeepEqual(got, testCase.want) {
			t.Errorf("%s: got %v want %v", testCase.name, got, testCase.want)
		}
	}
}

func TestExtractNonSlice(t *testing.T) {
	boolptr := true
	testCases := []interface{}{
		4,
		uintptr(1),
		"",
		"test",
		true,
		false,
		&boolptr,
	}

	for _, testCase := range testCases {
		got := textra.Extract(testCase)

		// nil / empty checks
		if got != nil {
			t.Errorf("%T: result should be nil", testCase)
		}
	}
}

// func TestFilter(t *testing.T) {
// 	testCases := []struct {
// 		name    string
// 		str     any
// 		lookFor string
// 		want    textra.StructTags
// 	}{
// 		{"look for empty", (*Populated)(nil), "", nil},
// 		{
// 			"look for test", (*Populated)(nil), "test",
// 			textra.StructTags{
// 				"OneTag": []textra.Tag{{Tag: "test", Value: "notag"}},
// 				"TwoTag": []textra.Tag{
// 					{"pg", "pgvalue", nil},
// 					{"json", "two_tag", []string{"omitempty"}},
// 					{"test", "two_tag", []string{"opt1", "opt2"}},
// 				},
// 			},
// 		},
// 		{
// 			"look for json", (*Populated)(nil), "json",
// 			textra.StructTags{
// 				"TwoTag": []textra.Tag{
// 					{"pg", "pgvalue", nil},
// 					{"json", "two_tag", []string{"omitempty"}},
// 					{"test", "two_tag", []string{"opt1", "opt2"}},
// 				},
// 				"ThreeTag": []textra.Tag{
// 					{"pg", "three_tag", nil},
// 					{"json", "three_tag", nil},
// 				},
// 			},
// 		},
// 		{
// 			"look for pg", (*Populated)(nil), "pg",
// 			textra.StructTags{
// 				"TwoTag": []textra.Tag{
// 					{"pg", "pgvalue", nil},
// 					{"json", "two_tag", []string{"omitempty"}},
// 					{"test", "two_tag", []string{"opt1", "opt2"}},
// 				},
// 				"ThreeTag": []textra.Tag{
// 					{"pg", "three_tag", nil},
// 					{"json", "three_tag", nil},
// 				},
// 			},
// 		},
// 	}

// 	for _, testCase := range testCases {
// 		testCase := testCase

// 		got := textra.Extract(testCase.str).Filter(testCase.lookFor)

// 		if !reflect.DeepEqual(got, testCase.want) && (len(got) != 0 && len(testCase.want) != 0) {
// 			t.Errorf("%s: got %v want %v", testCase.name, got, testCase.want)
// 		}
// 	}
// }

// func TestFilterMany(t *testing.T) {
// 	testCases := []struct {
// 		name    string
// 		str     any
// 		lookFor []string
// 		want    textra.StructTags
// 	}{
// 		{"look for nil", (*Populated)(nil), nil, nil},
// 		{"look for empty", (*Populated)(nil), []string{}, nil},
// 		{
// 			"look for json", (*Populated)(nil),
// 			[]string{"json"},
// 			textra.StructTags{
// 				"TwoTag": []textra.Tag{
// 					{"pg", "pgvalue", nil},
// 					{"json", "two_tag", []string{"omitempty"}},
// 					{"test", "two_tag", []string{"opt1", "opt2"}},
// 				},
// 				"ThreeTag": []textra.Tag{
// 					{"pg", "three_tag", nil},
// 					{"json", "three_tag", nil},
// 				},
// 			},
// 		},
// 		{
// 			"look for json and test", (*Populated)(nil),
// 			[]string{"json", "test"},
// 			textra.StructTags{
// 				"OneTag": []textra.Tag{{Tag: "test", Value: "notag"}},
// 				"TwoTag": []textra.Tag{
// 					{"pg", "pgvalue", nil},
// 					{"json", "two_tag", []string{"omitempty"}},
// 					{"test", "two_tag", []string{"opt1", "opt2"}},
// 				},
// 				"ThreeTag": []textra.Tag{
// 					{"pg", "three_tag", nil},
// 					{"json", "three_tag", nil},
// 				},
// 			},
// 		},
// 	}

// 	for _, testCase := range testCases {
// 		_ = testCase

// 		got := textra.Extract(testCase.str).FilterAny(testCase.lookFor...)

// 		if !reflect.DeepEqual(got, testCase.want) && (len(got) != 0 && len(testCase.want) != 0) {
// 			t.Errorf("%s: got %v want %v", testCase.name, got, testCase.want)
// 		}
// 	}
// }

// func TestFilterFunc(t *testing.T) {
// 	testCases := []struct {
// 		name string
// 		str  any
// 		fn   func(string, textra.Tags) bool
// 		want textra.StructTags
// 	}{
// 		{
// 			"look for threetag", (*Populated)(nil),
// 			func(field string, _ textra.Tags) bool {
// 				return field == "ThreeTag"
// 			},
// 			textra.StructTags{
// 				"ThreeTag": []textra.Tag{
// 					{"pg", "three_tag", nil},
// 					{"json", "three_tag", nil},
// 				},
// 			},
// 		},
// 		{
// 			"look for len3", (*Populated)(nil),
// 			func(_ string, tags textra.Tags) bool {
// 				return len(tags) == 3
// 			},
// 			textra.StructTags{
// 				"TwoTag": []textra.Tag{
// 					{"pg", "pgvalue", nil},
// 					{"json", "two_tag", []string{"omitempty"}},
// 					{"test", "two_tag", []string{"opt1", "opt2"}},
// 				},
// 			},
// 		},
// 	}

// 	for _, testCase := range testCases {
// 		_ = testCase

// 		got := textra.Extract(testCase.str).FilterFunc(testCase.fn)

// 		if !reflect.DeepEqual(got, testCase.want) && (len(got) != 0 && len(testCase.want) != 0) {
// 			t.Errorf("%s: got %v want %v", testCase.name, got, testCase.want)
// 		}
// 	}
// }
