// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package json

import (
	"fmt"
	"reflect"
)

// This example shows UnmarshalType accepts a type parameter and returns the object as that type if possible.
func ExampleUnmarshalType_map() {
	str := "{\"a\":1,\"b\":2,\"c\":3}"
	obj, err := UnmarshalType([]byte(str), reflect.TypeOf(map[string]float64{}))
	if err != nil {
		panic(err)
	}
	m, ok := obj.(map[string]float64)
	if ok {
		if v, ok := m["b"]; ok {
			fmt.Println("B:", v)
		}
	}
	// Output: B: 2
}

// This example shows that UnmarshalType accepts a type parameter and returns the object as that type if possible.
func ExampleUnmarshalType_slice() {
	str := "[1,2,3]"
	obj, err := UnmarshalType([]byte(str), reflect.TypeOf([]float64{}))
	if err != nil {
		panic(err)
	}
	slc, ok := obj.([]float64)
	if ok {
		fmt.Println("Length:", len(slc))
	}
	// Output: Length: 3
}
