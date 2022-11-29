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
		tagSplitted := strings.Split(tagStr, ":")

		k := tagSplitted[0]

		// Removing quotes, "value,  optional" -> value,  optional
		v := tagSplitted[1][1 : len(tagSplitted[1])-1]
		// Removing extra spaces, value,  optional -> value,optional
		v = strings.ReplaceAll(v, " ", "")
		vs := strings.Split(v, ",")

		tag := Tag{
			Tag:   k,
			Value: vs[0],
		}

		if len(vs) > 1 {
			tag.Optional = vs[1:]
		}

		parsed = append(parsed, tag)
	}

	return parsed
}
