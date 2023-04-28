package textra

import (
	"fmt"
)

// Field represents one struct field.
type Field struct {
	Name string `json:"name"`
	// Type is a type of a field, like "time.Time" or "*string"
	Type string `json:"type"`
	Tags Tags   `json:"tags,omitempty"`
}

// FieldTag is like Field but it has only one tag.
// It's used as an output of some functions (like Only()).
type FieldTag struct {
	Name string `json:"name"`
	// Type is a type of a field, like "time.Time" or "*string"
	Type string `json:"type"`
	Tag  Tag    `json:"tag,omitempty"`
}

func (f Field) String() string {
	return fmt.Sprintf("%s(%s):%s", f.Name, f.Type, f.Tags.String())
}

func (f FieldTag) String() string {
	return fmt.Sprintf("%s(%s):%s", f.Name, f.Type, f.Tag.String())
}
