# Textra

Textra is a simple and fast *****t*****ags *****extra*****ctor package that helps to work with structs tags.

Textra parses struct tags and returns them as a slice. It also does provide extra functionality like filtering.

Gathering json tags to feed it some other service is probably the most common usecase of this.

_Initially I built it for another private project, but decided to try open-source it, since it could be useful in some use-cases_

## Go Get

```shell
go get github.com/Ravcii/textra
```

## Examples

Basic usage:

```go
package main

import (
	"fmt"

	"github.com/Ravcii/textra"
)

type Tester struct {
	NoTags   bool
	WithTag  string `json:"with_tag,omitempty"`
	WithTags string `json:"with_tags"          sql:"with_tag"`
	SqlOnly  string `sql:"sql_only"`
}

func main() {
	tags := textra.Extract((*Tester)(nil))
	fmt.Println("Basic: ", tags)
	filtered := tags.Filter("sql")
	fmt.Println("Filtered: ", filtered)
	filteredMany := tags.FilterMany("sql", "json")
	fmt.Println("Filtered Many: ", filteredMany)
}

```

```
Basic: 		 map[SqlOnly:[sql:sql_only] WithTag:[json:with_tag,omitempty] WithTags:[json:with_tags sql:with_tag]]
Filtered: 	 map[SqlOnly:[sql:sql_only] WithTags:[json:with_tags sql:with_tag]]
Filtered Many: 	 map[SqlOnly:[sql:sql_only] WithTag:[json:with_tag,omitempty] WithTags:[json:with_tags sql:with_tag]]
```

## TODO:

- [ ] Add blacklisting
- [ ] Add better filtering options (FilterFunc for example)
- [ ] Some sugar for common tags, like json's omitempty
- [ ] Better string representation