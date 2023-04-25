package textra

import (
	"strings"
)

// Struct represents a single struct.
type Struct []Field

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
	tagsUnique := toUniqueMap(tags...)

	for _, field := range s {
		for _, t := range field.Tags {
			if _, ok := tagsUnique[t.Tag]; ok {
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
	tagsUnique := toUniqueMap(tags...)
	shouldMatch := len(tags)

	for _, field := range s {
		matched := 0

		for _, t := range field.Tags {
			if _, ok := tagsUnique[t.Tag]; ok {
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
	fieldsUnique := toUniqueMap(fields...)

	for _, field := range s {
		if _, ok := fieldsUnique[field.Name]; !ok {
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
