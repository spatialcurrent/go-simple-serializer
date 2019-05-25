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
	"strings"
)

func ExampleRead() {
	in := `
	a=b
	hello="beautiful world"
	hello="beautiful \"wide\" world"
  `
	out, err := Read(&ReadInput{
		Type:          reflect.TypeOf([]interface{}{}),
		Reader:        strings.NewReader(in),
		SkipLines:     0,
		SkipBlanks:    true,
		SkipComments:  false,
		LineSeparator: []byte("\n")[0],
		DropCR:        true,
		Comment:       "",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(out)
	// Output: [map[a:b] map[hello:beautiful world] map[hello:beautiful "wide" world]]
}
