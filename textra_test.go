package textra_test

import (
	"reflect"
	"testing"

	"github.com/Ravcii/textra"
)

type Empty struct{}

type Populated struct {
	NoTag    Empty
	OneTag   Empty `test:"notag"`
	TwoTag   Empty `pg:"pgvalue"   json:"two_tag,omitempty" test:"two_tag,opt1,opt2"`
	ThreeTag Empty `pg:"three_tag" json:"three_tag"`
}

func TestExtract(t *testing.T) {
	expected := textra.StructTags{
		"OneTag": []textra.Tag{{Tag: "test", Value: "notag"}},
		"TwoTag": []textra.Tag{
			{"pg", "pgvalue", nil},
			{"json", "two_tag", []string{"omitempty"}},
			{"test", "two_tag", []string{"opt1", "opt2"}},
		},
		"ThreeTag": []textra.Tag{
			{"pg", "three_tag", nil},
			{"json", "three_tag", nil},
		},
	}

	testCases := []struct {
		name string
		str  any
		want textra.StructTags
	}{
		{"Empty, non-pointer", Empty{}, nil},
		{"Empty, pointer", &Empty{}, nil},
		{"Empty, nil-check", (*Empty)(nil), nil},
		{"Populated, non-pointer", Populated{}, expected},
		{"Populated, pointer", &Populated{}, expected},
		{"Populated, nil-check", (*Populated)(nil), expected},
	}

	for _, testCase := range testCases {
		testCase := testCase

		got := textra.Extract(testCase.str)

		if !reflect.DeepEqual(got, testCase.want) && (len(got) != 0 && len(testCase.want) != 0) {
			t.Errorf("%s: got %v want %v", testCase.name, got, testCase.want)
		}
	}
}

func TestFilter(t *testing.T) {
	testCases := []struct {
		name    string
		str     any
		lookFor string
		want    textra.StructTags
	}{
		{"look for empty", (*Populated)(nil), "", nil},
		{
			"look for test", (*Populated)(nil), "test",
			textra.StructTags{
				"OneTag": []textra.Tag{{Tag: "test", Value: "notag"}},
				"TwoTag": []textra.Tag{
					{"pg", "pgvalue", nil},
					{"json", "two_tag", []string{"omitempty"}},
					{"test", "two_tag", []string{"opt1", "opt2"}},
				},
			},
		},
		{
			"look for json", (*Populated)(nil), "json",
			textra.StructTags{
				"TwoTag": []textra.Tag{
					{"pg", "pgvalue", nil},
					{"json", "two_tag", []string{"omitempty"}},
					{"test", "two_tag", []string{"opt1", "opt2"}},
				},
				"ThreeTag": []textra.Tag{
					{"pg", "three_tag", nil},
					{"json", "three_tag", nil},
				},
			},
		},
		{
			"look for pg", (*Populated)(nil), "pg",
			textra.StructTags{
				"TwoTag": []textra.Tag{
					{"pg", "pgvalue", nil},
					{"json", "two_tag", []string{"omitempty"}},
					{"test", "two_tag", []string{"opt1", "opt2"}},
				},
				"ThreeTag": []textra.Tag{
					{"pg", "three_tag", nil},
					{"json", "three_tag", nil},
				},
			},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		got := textra.Extract(testCase.str).Filter(testCase.lookFor)

		if !reflect.DeepEqual(got, testCase.want) && (len(got) != 0 && len(testCase.want) != 0) {
			t.Errorf("%s: got %v want %v", testCase.name, got, testCase.want)
		}
	}
}

func TestFilterMany(t *testing.T) {
	testCases := []struct {
		name    string
		str     any
		lookFor []string
		want    textra.StructTags
	}{
		{"look for nil", (*Populated)(nil), nil, nil},
		{"look for empty", (*Populated)(nil), []string{}, nil},
		{
			"look for json", (*Populated)(nil),
			[]string{"json"},
			textra.StructTags{
				"TwoTag": []textra.Tag{
					{"pg", "pgvalue", nil},
					{"json", "two_tag", []string{"omitempty"}},
					{"test", "two_tag", []string{"opt1", "opt2"}},
				},
				"ThreeTag": []textra.Tag{
					{"pg", "three_tag", nil},
					{"json", "three_tag", nil},
				},
			},
		},
		{
			"look for json and test", (*Populated)(nil),
			[]string{"json", "test"},
			textra.StructTags{
				"OneTag": []textra.Tag{{Tag: "test", Value: "notag"}},
				"TwoTag": []textra.Tag{
					{"pg", "pgvalue", nil},
					{"json", "two_tag", []string{"omitempty"}},
					{"test", "two_tag", []string{"opt1", "opt2"}},
				},
				"ThreeTag": []textra.Tag{
					{"pg", "three_tag", nil},
					{"json", "three_tag", nil},
				},
			},
		},
	}

	for _, testCase := range testCases {
		_ = testCase

		got := textra.Extract(testCase.str).FilterMany(testCase.lookFor...)

		if !reflect.DeepEqual(got, testCase.want) && (len(got) != 0 && len(testCase.want) != 0) {
			t.Errorf("%s: got %v want %v", testCase.name, got, testCase.want)
		}
	}
}

func TestFilterFunc(t *testing.T) {
	testCases := []struct {
		name string
		str  any
		fn   func(string, textra.Tags) bool
		want textra.StructTags
	}{
		{
			"look for threetag", (*Populated)(nil),
			func(field string, _ textra.Tags) bool {
				return field == "ThreeTag"
			},
			textra.StructTags{
				"ThreeTag": []textra.Tag{
					{"pg", "three_tag", nil},
					{"json", "three_tag", nil},
				},
			},
		},
		{
			"look for len3", (*Populated)(nil),
			func(_ string, tags textra.Tags) bool {
				return len(tags) == 3
			},
			textra.StructTags{
				"TwoTag": []textra.Tag{
					{"pg", "pgvalue", nil},
					{"json", "two_tag", []string{"omitempty"}},
					{"test", "two_tag", []string{"opt1", "opt2"}},
				},
			},
		},
	}

	for _, testCase := range testCases {
		_ = testCase

		got := textra.Extract(testCase.str).FilterFunc(testCase.fn)

		if !reflect.DeepEqual(got, testCase.want) && (len(got) != 0 && len(testCase.want) != 0) {
			t.Errorf("%s: got %v want %v", testCase.name, got, testCase.want)
		}
	}
}
