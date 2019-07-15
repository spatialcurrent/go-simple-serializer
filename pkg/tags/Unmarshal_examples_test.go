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

// This example shows you can unmarshal a line of tags into a map.
func ExampleUnmarshal_map() {
	str := "a=1 b=2 c=3"
	obj, err := Unmarshal([]byte(str), '=')
	if err != nil {
		panic(err)
	}
	fmt.Println(obj)
	// Output: map[a:1 b:2 c:3]
}
