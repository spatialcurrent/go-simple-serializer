// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package tags

import (
	"fmt"
	"reflect"
)

// This example shows you can unmarshal a line of tags into a map.
func ExampleUnmarshalType_map() {
	str := "a=1 b=2 c=3"
	obj, err := UnmarshalType([]byte(str), reflect.TypeOf(map[string]string{}))
	if err != nil {
		panic(err)
	}
	fmt.Println(obj)
	// Output: map[a:1 b:2 c:3]
}
