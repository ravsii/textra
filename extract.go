package textra

import (
	"reflect"
)

// Extract accept a struct (or a pointer to a struct) and returns a map
// of fields and their tags.
// If src is not a struct or a pointer to a struct, nil is returned.
func Extract(src interface{}) Struct {
	typ := reflect.TypeOf(src)

	// If str is a struct pointer
	if typ.Kind() == reflect.Ptr {
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
			Type: parseType(f.Type),
			Tags: parseTags(f.Tag),
		})
	}

	return result
}

// toUniqueMap accepts a slice of strings and returns a map of unique strings.
// This is used as an (arguably) better version of slices.Contains.
func toUniqueMap(strs ...string) map[string]struct{} {
	m := make(map[string]struct{})
	for _, s := range strs {
		m[s] = struct{}{}
	}

	return m
}
