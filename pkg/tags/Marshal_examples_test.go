// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package tags

import (
	"fmt"
)

import (
	"github.com/spatialcurrent/go-stringify/pkg/stringify"
)

// This example shows how to marshal a map into a line of tags.
func ExampleMarshal_map() {
	obj := map[string]interface{}{
		"a": 1,
		"b": 2,
		"c": 3,
	}
	b, err := Marshal(obj, stringify.NewStringer("", false, false, false), stringify.NewStringer("", false, false, false), true, false)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	// Output: a=1 b=2 c=3
}

// This example shows how to marshal a struct into a line of tags.
func ExampleMarshal_stuct() {
	in := struct {
		A string
		B string
		C string
	}{A: "1", B: "2", C: "3"}
	keySerializer := stringify.NewStringer("", false, false, false)
	valueSerializer := stringify.NewStringer("", false, false, false)
	b, err := Marshal(in, keySerializer, valueSerializer, true, false)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	// Output: A=1 B=2 C=3
}
