// =================================================================
//
// Copyright (C) 2018 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

// gss.so creates a shared library of Go that can be called by C, C++, or Python
//

package main

import (
	"C"
	"strings"
)

import (
	"github.com/spatialcurrent/go-simple-serializer/gss"
)

func main() {}

//export Convert
func Convert(input_string *C.char, input_format *C.char, input_header *C.char, input_comment *C.char, input_limit C.long, output_format *C.char, output_limit C.long, output_string **C.char) *C.char {

	s, err := gss.Convert(
		[]byte(C.GoString(input_string)),
		C.GoString(input_format),
		strings.Split(C.GoString(input_header), ","),
		C.GoString(input_comment),
		int(input_limit),
		C.GoString(output_format),
		int(output_limit),
		false)
	if err != nil {
		return C.CString(err.Error())
	}

	*output_string = C.CString(s)

	return nil
}

//export Version
func Version() *C.char {
	return C.CString(gss.VERSION)
}
