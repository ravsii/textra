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
	split := strings.Split(tagStr, ":")
	v := strings.Trim(split[1], "\"")
	vs := strings.Split(v, ",")
	value := strings.TrimSpace(vs[0])

	tag := Tag{
		Tag:   split[0],
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
	switch typ.Kind() {
	case reflect.Ptr:
		return "*" + typ.Elem().String()
	case reflect.Slice:
		return "[]" + typ.Elem().String()
	case reflect.Struct:
		if len(typ.PkgPath()) > 0 {
			return typ.PkgPath() + "." + typ.Name()
		}

		return typ.Kind().String()
	case reflect.Map:
		return "map[" + typ.Key().String() + "]" + typ.Elem().String()
	case reflect.Func:
		var args, results string

		for i := 0; i < typ.NumIn(); i++ {
			args += parseType(typ.In(i))

			if i != typ.NumIn()-1 {
				args += ", "
			}
		}

		for i := 0; i < typ.NumOut(); i++ {
			results += parseType(typ.Out(i))

			if i != typ.NumOut()-1 {
				results += ", "
			}
		}

		return "func(" + args + ") " + results
	case reflect.Interface:
		if _, ok := reflect.New(typ).Interface().(*error); ok {
			return "error"
		}

		fallthrough
	default:
		return typ.Kind().String()
	}
}
