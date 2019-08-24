// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package jsonl

import (
	"bytes"
	"fmt"

	"github.com/spatialcurrent/go-stringify/pkg/stringify"
)

// This example shows you can marshal a single map into a JSON object
func ExampleWrite_map() {
	obj := map[string]interface{}{
		"a": 1,
		"b": 2,
		"c": 3,
	}
	buf := new(bytes.Buffer)
	err := Write(&WriteInput{
		Writer:        buf,
		LineSeparator: "\n",
		KeySerializer: stringify.NewStringer("", false, false, false),
		Object:        obj,
		Pretty:        false,
		Limit:         -1,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(buf.String())
	// Output: {"a":1,"b":2,"c":3}
}

// This example shows you can marshal a single map into a JSON object
func ExampleWrite_map_upper() {
	obj := map[string]interface{}{
		"a": 1,
		"b": 2,
		"c": 3,
	}
	buf := new(bytes.Buffer)
	err := Write(&WriteInput{
		Writer:        buf,
		LineSeparator: "\n",
		KeySerializer: stringify.NewStringer("", false, false, true),
		Object:        obj,
		Pretty:        false,
		Limit:         -1,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(buf.String())
	// Output: {"A":1,"B":2,"C":3}
}

// This examples shows that you can marshal a slice of maps into lines of JSON.
func ExampleWrite_slice() {
	obj := []interface{}{
		map[string]interface{}{
			"a": 1,
			"b": 2,
			"c": 3,
		},
		map[string]interface{}{
			"a": 4,
			"b": 5,
			"c": 6,
		},
	}
	buf := new(bytes.Buffer)
	err := Write(&WriteInput{
		Writer:        buf,
		LineSeparator: "\n",
		KeySerializer: stringify.NewStringer("", false, false, false),
		Object:        obj,
		Pretty:        false,
		Limit:         -1,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(buf.String())
	// Output: {"a":1,"b":2,"c":3}
	// {"a":4,"b":5,"c":6}
}

// This example shows you can marshal a slice maps into lines of JSON objects with pretty formatting.
func ExampleWrite_pretty() {
	obj := []interface{}{
		map[string]interface{}{
			"a": 1,
			"b": 2,
			"c": 3,
		},
		map[string]interface{}{
			"a": 4,
			"b": 5,
			"c": 6,
		},
	}
	buf := new(bytes.Buffer)
	err := Write(&WriteInput{
		Writer:        buf,
		LineSeparator: "\n",
		KeySerializer: stringify.NewStringer("", false, false, false),
		Object:        obj,
		Pretty:        true,
		Limit:         -1,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(buf.String())
	// Output: {
	//   "a": 1,
	//   "b": 2,
	//   "c": 3
	//}
	//{
	//   "a": 4,
	//   "b": 5,
	//   "c": 6
	//}

}
