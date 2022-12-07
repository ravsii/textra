package textra

import (
	"fmt"
)

// Field represents a one struct field.
type Field struct {
	Name string `json:"name"`
	// Type is stringified type, like "time.Time" or "*string"
	Type string `json:"type"`
	Tags Tags   `json:"tags"`
}

// FieldTag is like but it has only one tag.
// It's used as output of some functions (like Only()).
type FieldTag struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Tag  Tag    `json:"tag"`
}

func (f Field) String() string {
	return fmt.Sprintf("%s(%s):%s", f.Name, f.Type, f.Tags)
}

func (f FieldTag) String() string {
	return fmt.Sprintf("%s(%s):%s", f.Name, f.Type, f.Tag.String())
}
