// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package properties

import (
	"bytes"
	"fmt"

	"github.com/spatialcurrent/go-stringify/pkg/stringify"
)

// This example shows you can write a map into properties text.
func ExampleWrite_map() {
	obj := map[string]interface{}{
		"a": 1,
		"b": 2,
		"c": 3,
	}
	buf := new(bytes.Buffer)
	err := Write(&WriteInput{
		Writer:            buf,
		KeyValueSeparator: "=",
		LineSeparator:     "\n",
		Object:            obj,
		KeySerializer:     stringify.NewStringer("", false, false, false),
		ValueSerializer:   stringify.NewStringer("", false, false, false),
		Sorted:            true,
		EscapePrefix:      "\\",
		EscapeSpace:       true,
		EscapeEqual:       true,
		EscapeColon:       false,
		EscapeNewLine:     false,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(buf.String())
	// Output: a=1
	//b=2
	//c=3
}

// This example shows you can write a map into properties text using a custom key-value separator, in this case a semicolon.
func ExampleWrite_semicolon() {
	obj := map[string]interface{}{
		"a": 1,
		"b": 2,
		"c": 3,
	}
	buf := new(bytes.Buffer)
	err := Write(&WriteInput{
		Writer:            buf,
		KeyValueSeparator: ";", // specify the separator for each key-value pair
		LineSeparator:     "\n",
		Object:            obj,
		KeySerializer:     stringify.NewStringer("", false, false, false),
		ValueSerializer:   stringify.NewStringer("", false, false, false),
		Sorted:            true,
		EscapePrefix:      "\\",
		EscapeSpace:       true,
		EscapeEqual:       true,
		EscapeColon:       false,
		EscapeNewLine:     false,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(buf.String())
	// Output: a;1
	//b;2
	//c;3
}

// This example shows you can write a map into properties text using a custom value serializer.
// In this case, nil values are written as a dash.
func ExampleWrite_valueSerializer() {
	obj := map[string]interface{}{
		"a": 1,
		"b": 2,
		"c": nil,
	}
	buf := new(bytes.Buffer)
	err := Write(&WriteInput{
		Writer:            buf,
		KeyValueSeparator: "=", // specify the separator for each key-value pair
		LineSeparator:     "\n",
		Object:            obj,
		KeySerializer:     stringify.NewStringer("", false, false, false),
		ValueSerializer:   stringify.NewStringer("-", false, false, false), // specify the no-data value
		Sorted:            true,
		EscapePrefix:      "\\",
		EscapeSpace:       true,
		EscapeEqual:       true,
		EscapeColon:       false,
		EscapeNewLine:     false,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(buf.String())
	// Output: a=1
	//b=2
	//c=-
}
