package textra

import (
	"reflect"
	"regexp"
	"strings"

	"golang.org/x/exp/slices"
)

// Tag represents a single struct tag, like
//
//	json:"value".
type Tag struct {
	Tag      string
	Value    string
	Optional []string
}

// Tags is a slice of tags.
type Tags []Tag

// StructTags is a map that stores a slice of tags related to each field in a struct.
type StructTags map[string]Tags

var whitespaceRegexp = regexp.MustCompile(`\s{2,}`)

func (t Tag) String() string {
	b := strings.Builder{}
	b.WriteString(t.Tag + ":" + t.Value)

	for _, v := range t.Optional {
		b.WriteString("," + v)
	}

	return b.String()
}

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
		tags := parseTags(tag)
		structTags[typ.Field(i).Name] = tags
	}

	return structTags
}

func parseTags(tag reflect.StructTag) Tags {
	s := string(tag)
	if len(s) == 0 {
		return nil
	}

	tagsStr := whitespaceRegexp.ReplaceAllLiteralString(s, " ")
	splitted := strings.Split(tagsStr, " ")
	tags := make(Tags, 0, len(splitted))

	for _, tagStr := range splitted {
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

		tags = append(tags, tag)
	}

	return tags
}

// Filter returns a map of fields and associated tag, if given tag key is present.
func (m StructTags) Filter(tag string) map[string]Tag {
	newMap := make(map[string]Tag, 0)

	for field, tags := range m {
		for _, t := range tags {
			if t.Tag == tag {
				newMap[field] = t
				break
			}
		}
	}

	return newMap
}

// FilterMany returns a map of fields and associated tags for given tag keys.
// If no tags are passed, nil is returned.
func (m StructTags) FilterMany(tags ...string) StructTags {
	if len(tags) == 0 {
		return nil
	}

	newMap := make(StructTags, 0)

	for field, strTags := range m {
		foundTags := make(Tags, 0)

		for _, t := range strTags {
			if !slices.Contains(tags, t.Tag) {
				foundTags = append(foundTags, t)
			}
		}

		newMap[field] = foundTags
	}

	return newMap
}
