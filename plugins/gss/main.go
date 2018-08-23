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

var GO_GSS_VERSION = "0.0.2"

func main() {}

//export Convert
func Convert(input_string *C.char, input_format *C.char, input_header *C.char, input_comment *C.char, output_format *C.char, output_string **C.char) *C.char {

	s, err := gss.Convert(
		C.GoString(input_string),
		C.GoString(input_format),
		strings.Split(C.GoString(input_header), ","),
		C.GoString(input_comment),
		C.GoString(output_format))
	if err != nil {
		return C.CString(err.Error())
	}

	*output_string = C.CString(s)

	return nil
}
