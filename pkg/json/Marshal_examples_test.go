// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package json

import (
	"fmt"
)

// This example shows you can marshal a map into a JSON object
func ExampleMarshal_map() {
	obj := map[string]interface{}{
		"a": 1,
		"b": 2,
		"c": 3,
	}
	b, err := Marshal(obj, false)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	// Output: {"a":1,"b":2,"c":3}
}

// This examples shows that you can marshal a slice into a JSON array.
func ExampleMarshal_slice() {
	obj := []interface{}{"1", "2", "3"}
	b, err := Marshal(obj, false)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	// Output: ["1","2","3"]
}

// This example shows you can marshal a map into a JSON object with pretty formatting.
func ExampleMarshal_pretty() {
	obj := map[string]interface{}{
		"a": 1,
		"b": 2,
		"c": 3,
	}
	b, err := Marshal(obj, true)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	// Output: {
	//   "a": 1,
	//   "b": 2,
	//   "c": 3
	//}
}
