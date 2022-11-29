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
	// Removing extra spaces, value,  optional -> value,optional
	v = strings.ReplaceAll(v, " ", "")
	vs := strings.Split(v, ",")

	tag := Tag{
		Tag:   splitted[0],
		Value: vs[0],
	}

	if len(vs) > 1 {
		tag.Optional = vs[1:]
	}

	return tag
}
