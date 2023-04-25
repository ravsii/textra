package main

import (
	"fmt"
	"time"

	"github.com/ravsii/textra"
)

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
