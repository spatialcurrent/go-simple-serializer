// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package tags

import (
	"bytes"
	"fmt"
)

import (
	"github.com/spatialcurrent/go-stringify/pkg/stringify"
)

// This example shows you can marshal a single map into a JSON object
func ExampleWrite_map() {
	obj := map[string]string{"a": "b", "c": "beautiful world", "d": "beautiful \"wide\" world"}
	buf := new(bytes.Buffer)
	valueSerializer := stringify.DefaultValueStringer("")
	err := Write(&WriteInput{
		Writer:          buf,
		LineSeparator:   "\n",
		Object:          obj,
		ValueSerializer: valueSerializer,
		Sorted:          true,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(buf.String())
	// Output: a=b c="beautiful world" d="beautiful \"wide\" world"
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
	valueSerializer := stringify.DefaultValueStringer("")
	err := Write(&WriteInput{
		Writer:          buf,
		LineSeparator:   "\n",
		Object:          obj,
		ValueSerializer: valueSerializer,
		Sorted:          true,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(buf.String())
	// Output: a=1 b=2 c=3
	// a=4 b=5 c=6
}
