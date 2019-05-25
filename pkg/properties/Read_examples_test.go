// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package properties

import (
	"fmt"
	"reflect"
	"strings"
)

// This example shows you can write a map into properties text.
func ExampleRead_map() {
	in := `
  a=1
  b:2
  c true
  d=nil
  e=
  `

	out, err := Read(&ReadInput{
		Type:            reflect.TypeOf(map[string]string{}),
		Reader:          strings.NewReader(in),
		LineSeparator:   []byte("\n")[0],
		Comment:         "",
		Trim:            true,
		UnescapeSpace:   false,
		UnescapeEqual:   false,
		UnescapeColon:   false,
		UnescapeNewLine: false,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(out)
	// Output: map[a:1 b:2 c:true d:nil e:]
}

// This example shows you can read from properties text with comments
func ExampleRead_comment() {
	in := `
  a=1
  b:2
  # ignore this comment
  c true
  d=nil
  # and ignore this one too
  e=
  `

	out, err := Read(&ReadInput{
		Type:            reflect.TypeOf(map[string]string{}),
		Reader:          strings.NewReader(in),
		LineSeparator:   []byte("\n")[0],
		Comment:         "#",
		Trim:            true,
		UnescapeSpace:   false,
		UnescapeEqual:   false,
		UnescapeColon:   false,
		UnescapeNewLine: false,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(out)
	// Output: map[a:1 b:2 c:true d:nil e:]
}

// This example shows you can use a custom line separator to read from a series of semi-colon separated key-value pairs.
func ExampleRead_semicolon() {

	in := `a=1;b:2;c true;d=nil;e=`

	out, err := Read(&ReadInput{
		Type:            reflect.TypeOf(map[string]string{}),
		Reader:          strings.NewReader(in),
		LineSeparator:   []byte(";")[0],
		Comment:         "",
		Trim:            true,
		UnescapeSpace:   false,
		UnescapeEqual:   false,
		UnescapeColon:   false,
		UnescapeNewLine: false,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(out)

	// Output: map[a:1 b:2 c:true d:nil e:]
}
