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

// This example shows you can unmarshal a JSON object into a map.
func ExampleUnmarshal_map() {
	str := "{\"a\":1,\"b\":2,\"c\":3}"
	obj, err := Unmarshal([]byte(str))
	if err != nil {
		panic(err)
	}
	m, ok := obj.(map[string]interface{})
	if ok {
		if v, ok := m["b"]; ok {
			fmt.Println("B:", v)
		}
	}
	// Output: B: 2
}

// This examples shows that you can unmarshal a JSON array into a slice.
func ExampleUnmarshal_slice() {
	str := "[1,2,3]"
	obj, err := Unmarshal([]byte(str))
	if err != nil {
		panic(err)
	}
	slc, ok := obj.([]interface{})
	if ok {
		fmt.Println("Length:", len(slc))
	}
	// Output: Length: 3
}
