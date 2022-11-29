package textra

import (
	"reflect"

	"golang.org/x/exp/slices"
)

// StructTags is a map that stores a slice of tags related to each field in a struct.
type StructTags map[string]Tags

// Extract accept a struct (or a pointer to a struct)
// and returns a map of fields and their tags.
func Extract(str any) StructTags {
	typ := reflect.TypeOf(str)

	// If str is a struct pointer
	if typ.Kind() == reflect.Pointer {
		typ = typ.Elem() // dereference it
	}

	if typ.Kind() != reflect.Struct {
		return nil
	}

	amount := typ.NumField()
	structTags := make(StructTags, amount)

	for i := 0; i < amount; i++ {
		tag := typ.Field(i).Tag

		s := string(tag)
		if len(s) == 0 {
			continue
		}

		tags := parseTags(tag)
		structTags[typ.Field(i).Name] = tags
	}

	return structTags
}

// Filter returns a map of fields and their tags, if a field has given tag.
func (m StructTags) Filter(tag string) StructTags {
	filtered := make(StructTags, len(m)/2)

	for field, tags := range m {
		for _, t := range tags {
			if t.Tag == tag {
				filtered[field] = tags
				break
			}
		}
	}

	return filtered
}

// FilterMany returns a map of fields and associated tags for given tag keys.
// If no tags are passed, nil is returned.
func (m StructTags) FilterMany(tags ...string) StructTags {
	if len(tags) == 0 {
		return nil
	}

	filtered := make(StructTags, len(m)/2)

	for field, strTags := range m {
		for _, t := range strTags {
			if slices.Contains(tags, t.Tag) {
				filtered[field] = strTags
				break
			}
		}
	}

	return filtered
}

// FilterFunc returns a map of fields and associated tags for given tag keys.
// fn is called for each field to decide whether that field should be included or not.
func (m StructTags) FilterFunc(fn func(string, Tags) bool) StructTags {
	filtered := make(StructTags, len(m)/2)

	for field, tags := range m {
		if fn(field, tags) {
			filtered[field] = tags
		}
	}

	return filtered
}
