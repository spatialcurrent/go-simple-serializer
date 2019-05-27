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

// This example shows you can marshal a map into a JSON object
func ExampleMarshal_map() {
	obj := map[string]interface{}{
		"a": 1,
		"b": 2,
		"c": 3,
	}
	b, err := Marshal(obj, stringify.DefaultValueStringer(""), true)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	// Output: a=1 b=2 c=3
}
