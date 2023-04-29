package main

import (
	"fmt"

	"github.com/ravsii/textra"
)

type Tester struct {
	NoTags   bool
	WithTag  string    `json:"with_tag,omitempty"`
	WithTags *string   `json:"with_tags"          sql:"with_tag"`
	SqlOnly  *[]string `pg:"sql_only"             sql:"sql_only"`
}

func main() {
	fmt.Println("basic")
	basic := textra.Extract((*Tester)(nil))
	for _, field := range basic {
		fmt.Println(field)
	}

	fmt.Println("removed")
	removed := basic.RemoveEmpty()
	for _, field := range removed {
		fmt.Println(field)
	}

	fmt.Println("onlySQL")
	onlySQL := removed.OnlyTag("sql")
	for _, field := range onlySQL {
		fmt.Println(field)
	}
}
