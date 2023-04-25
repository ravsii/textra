package textra

import (
	"reflect"
	"regexp"
	"strings"
)

var tagRegexp = regexp.MustCompile(`(\w+:\"[^\"]+\")`)

func parseTags(tag reflect.StructTag) Tags {
	tags := tagRegexp.FindAllString(string(tag), -1)
	parsed := make(Tags, 0, len(tags))

	for _, tagStr := range tags {
		parsed = append(parsed, parseTag(tagStr))
	}

	return parsed
}

func parseTag(tagStr string) Tag {
	splitted := strings.Split(tagStr, ":")

	// Removing quotes, "value,  optional" -> value,  optional
	v := splitted[1][1 : len(splitted[1])-1]

	vs := strings.Split(v, ",")
	value := strings.TrimSpace(vs[0])

	tag := Tag{
		Tag:   splitted[0],
		Value: value,
	}

	if len(vs) > 1 {
		tag.Optional = make([]string, 0, len(vs)-1)
		for _, opt := range vs[1:] {
			tag.Optional = append(tag.Optional, strings.TrimSpace(opt))
		}
	}

	return tag
}

func parseType(typ reflect.Type) string {
	if typ.Kind() == reflect.Ptr {
		return "*" + typ.Elem().String()
	}

	if typ.Kind() == reflect.Struct {
		if len(typ.PkgPath()) > 0 {
			return typ.PkgPath() + "." + typ.Name()
		}

		return typ.Kind().String()
	}

	return typ.Kind().String()
}
