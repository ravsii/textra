package textra

import (
	"fmt"
	"reflect"
	"strings"

	"golang.org/x/exp/slices"
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

// Struct represents a single struct.
type Struct []Field

// Extract accept a struct (or a pointer to a struct) and returns a map of fields and their tags.
// if src is not a struct or a pointer to a struct, nil is returned.
func Extract(src any) Struct {
	typ := reflect.TypeOf(src)

	// If str is a struct pointer
	if typ.Kind() == reflect.Pointer {
		typ = typ.Elem() // dereference it
	}

	if typ.Kind() != reflect.Struct {
		return nil
	}

	amount := typ.NumField()
	result := make(Struct, 0, amount)

	var f reflect.StructField

	for i := 0; i < amount; i++ {
		f = typ.Field(i)

		result = append(result, Field{
			Name: f.Name,
			Type: f.Type.Name(),
			Tags: parseTags(f.Tag),
		})
	}

	return result
}

// ByTagName returns a slice of fields which contain given tag.
func (s Struct) ByTagName(tag string) Struct {
	filtered := make(Struct, 0)

	for _, field := range s {
		for _, t := range field.Tags {
			if t.Tag == tag {
				filtered = append(filtered, field)
				break
			}
		}
	}

	return filtered
}

// ByAnyTagName returns a slice of fields which contain at least one tag of the given tags.
func (s Struct) ByAnyTagName(tags ...string) Struct {
	filtered := make(Struct, 0)

	for _, field := range s {
		for _, t := range field.Tags {
			if slices.Contains(tags, t.Tag) {
				filtered = append(filtered, field)
				break
			}
		}
	}

	return filtered
}

// ByTagNames returns a slice of fields which contain all of the given tags.
func (s Struct) ByTagNames(tags ...string) Struct {
	filtered := make(Struct, 0)

	shouldMatch := len(tags)

	for _, field := range s {
		matched := 0

		for _, t := range field.Tags {
			if slices.Contains(tags, t.Tag) {
				matched++
			}
		}

		if matched != 0 && matched == shouldMatch {
			filtered = append(filtered, field)
		}
	}

	return filtered
}

// Field returns a field by name.
func (s Struct) Field(name string) (Field, bool) {
	for _, field := range s {
		if field.Name == name {
			return field, true
		}
	}

	return Field{}, false
}

// FilterFunc returns a slice of fields, filtered by fn(field) == true.
func (s Struct) FilterFunc(fn func(Field) bool) Struct {
	filtered := make(Struct, 0)

	for _, f := range s {
		if fn(f) {
			filtered = append(filtered, f)
		}
	}

	return filtered
}

// RemoveEmpty returns a map without fields that has no tags.
func (s Struct) RemoveEmpty() Struct {
	filtered := make(Struct, 0)

	for _, field := range s {
		if len(field.Tags) != 0 {
			filtered = append(filtered, field)
		}
	}

	return filtered
}

// RemoveFields copies original map but skips given fields on each field.
func (s Struct) RemoveFields(fields ...string) Struct {
	filtered := make(Struct, 0)

	for _, field := range s {
		if !slices.Contains(fields, field.Name) {
			filtered = append(filtered, field)
		}
	}

	return filtered
}

// Only returns StructTag (instead of Struct like most other) of a
// field and a tag with a given name.
func (s Struct) Only(name string) []FieldTag {
	filtered := make([]FieldTag, 0)

	for _, field := range s {
		if tag, ok := field.Tags.ByName(name); ok {
			filtered = append(filtered, FieldTag{
				Name: field.Name,
				Type: field.Type,
				Tag:  tag,
			})
		}
	}

	return filtered
}

func (s Struct) String() string {
	b := strings.Builder{}

	for _, field := range s {
		b.WriteString(field.String())
	}

	return b.String()
}

func (f Field) String() string {
	return fmt.Sprintf("%s(%s):%s", f.Name, f.Type, f.Tags)
}
