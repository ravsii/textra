package textra

import (
	"strings"
)

// Tags is a slice of tags.
type Tags []Tag

// ByName returns a tag by its name, if it exists.
func (t Tags) ByName(name string) (Tag, bool) {
	for _, tag := range t {
		if tag.Tag == name {
			return tag, true
		}
	}

	return Tag{}, false
}

func (t Tags) String() string {
	tags := make([]string, 0, len(t))

	for _, tag := range t {
		tags = append(tags, tag.String())
	}

	return "[" + strings.Join(tags, " ") + "]"
}

// Tag represents a single struct tag, like
//
//	`json:"value"`.
type Tag struct {
	// Tag holds tag's name.
	Tag   string `json:"tag"`
	Value string `json:"value"`
	// Optional holds the rest of the value which comes after the comma.
	// Example: In a tag like
	// 	`json:"id,pk,omitempty"`
	// Optional will contain
	// 	["pk", "omitempty"]
	Optional []string `json:"optional,omitempty"`
}

// OmitEmpty returns true if t.Optional contains "omitempty".
func (t Tag) OmitEmpty() bool {
	m := toUniqueMap(t.Optional...)
	_, ok := m["omitempty"]
	return ok
}

// Ignored is a shortcut for t.Value == "-".
func (t Tag) Ignored() bool {
	return t.Value == "-"
}

func (t Tag) String() string {
	b := strings.Builder{}
	b.WriteString(t.Tag + ":\"" + t.Value)

	for _, v := range t.Optional {
		b.WriteString("," + v)
	}

	b.WriteRune('"')

	return b.String()
}
