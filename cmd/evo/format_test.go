package main

import (
	"testing"

	"github.com/hiway-media/evo/lib/gpath"
)

func Test_FormatStruct(t *testing.T) {
	b, err := gpath.ReadFile("./test.go")
	if err != nil {
		panic(err)
	}

	code := formatStruct(string(b))
	_ = code
	f, err := gpath.Open("./test.go")
	if err != nil {
		panic(err)
	}

	f.WriteString(code)

}
