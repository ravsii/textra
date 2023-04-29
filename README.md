# Textra

[![Go Reference](https://pkg.go.dev/badge/github.com/ravsii/textra.svg)](https://pkg.go.dev/github.com/ravsii/textra) [![codecov](https://codecov.io/gh/ravsii/textra/branch/main/graph/badge.svg?token=C8WA38GNFV)](https://codecov.io/gh/ravsii/textra)

Textra is a zero-dependency, simple and fast struct tags parser library. It also has json tags for all structs, in case of JSON output.

Initially I built it for another private project, but decided to try to open source it, since it could be useful for someone. Because of that, it has some features that feel redundant, like having field type as a part of returned data

## Installation

```shell
go get github.com/ravsii/textra
```

## Examples

Basic usage:

```go
type Tester struct {
 NoTags   bool
 WithTag  string `json:"with_tag,omitempty"`
 WithTags string `json:"with_tags"          sql:"with_tag"`
 SqlOnly  string `sql:"sql_only"`
}

func main() {
 basic := textra.Extract((*Tester)(nil))
 for _, field := range basic {
  fmt.Println(field)
 }
}

```

```text
NoTags(bool):[]
WithTag(string):[json:"with_tag,omitempty"]
WithTags(string):[json:"with_tags" sql:"with_tag"]
SqlOnly(string):[sql:"sql_only"]
```

You can look at return types at [pkg.go.dev](https://pkg.go.dev/github.com/Ravcii/textra), but basically it returns a slice of fields with its types (as strings) and a slice of Tags for each field.

Now let's apply some functions:

```go
 removed := basic.RemoveEmpty()
 for _, field := range removed {
  fmt.Println(field)
 }
```

```text
SqlOnly(*[]string):[pg:"sql_only" sql:"sql_only"]
WithTag(string):[json:"with_tag,omitempty"]
WithTags(*string):[json:"with_tags" sql:"with_tag"]
```

What if we care only about SQL tags?

```go
 onlySQL := removed.OnlyTag("sql")
 for _, field := range onlySQL {
  fmt.Println(field)
 }
```

```text
SqlOnly(*[]string):sql:"sql_only"
WithTags(*string):sql:"with_tag"
```

_Only() is a bit special as it returns a Field of a different type, with `Tag` rather than `Tags`(=`[]Tag`)_

API is built like standard's `time` package, where chaining function will create new values, instead of modifying them.

Although it may be redundant, it also parses types a their string representation (for easier comparison or output, if you need it)

```go
type Types struct {
 intType        int
 intPType       *int
 byteType       byte
 bytePType      *byte
 byteArrType    []byte
 byteArrPType   []*byte
 bytePArrPType  *[]*byte
 runeType       rune
 runePType      *rune
 stringType     string
 stringPType    *string
 booleanType    bool
 booleanPType   *bool
 mapType        map[string]string
 mapPType       map[*string]*string
 mapPImportType map[*string]*time.Time
 chanType       chan int
 funcType       func() error
 funcParamsType func(arg1 int, arg2 string, arg3 map[*string]*time.Time) (int, error)
 importType     time.Time
 pointerType    *string
}

func main() {
 fields := textra.Extract((*Types)(nil))
 for _, field := range fields {
  fmt.Println(field.Name, field.Type)
 }
}
```

```text
intType int
intPType *int
byteType uint8
bytePType *uint8
byteArrType []uint8
byteArrPType []*uint8
bytePArrPType *[]*uint8
runeType int32
runePType *int32
stringType string
stringPType *string
booleanType bool
booleanPType *bool
mapType map[string]string
mapPType map[*string]*string
mapPImportType map[*string]*time.Time
chanType chan
funcType func() error
funcParamsType func(int, string, map[*string]*time.Time) int, error
importType time.Time
pointerType *string
```

### TODO

- [ ] Examples for go.dev
- [ ] ...?
