// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

// gss.so creates a shared library of Go that can be called by C, C++, or Python
//
//
//  - https://godoc.org/cmd/cgo
//  - https://blog.golang.org/c-go-cgo
//
package main

import (
	"C"
	"strings"
)

import (
	"github.com/spatialcurrent/go-simple-serializer/pkg/gss"
)

func main() {}

//export Convert
func Convert(inputString *C.char, inputFormat *C.char, inputHeader *C.char, inputComment *C.char, inputLazyQuotes C.int, inputSkipLines C.long, inputLimit C.long, outputFormat *C.char, outputHeader *C.char, outputLimit C.long, outputPretty C.int, outputSorted C.int, async C.int, outputString **C.char) *C.char {

	s, err := gss.Convert(&gss.ConvertInput{
		InputBytes:      []byte(C.GoString(inputString)),
		InputFormat:     C.GoString(inputFormat),
		InputHeader:     strings.Split(C.GoString(inputHeader), ","),
		InputComment:    C.GoString(inputComment),
		InputLazyQuotes: int(inputLazyQuotes) > 0,
		InputSkipLines:  int(inputSkipLines),
		InputLimit:      int(inputLimit),
		OutputFormat:    C.GoString(outputFormat),
		OutputHeader:    strings.Split(C.GoString(outputHeader), ","),
		OutputLimit:     int(outputLimit),
		OutputPretty:    int(outputPretty) > 0,
		OutputSorted:    int(outputSorted) > 0,
		Async:           int(async) > 0,
		Verbose:         false,
	})
	if err != nil {
		return C.CString(err.Error())
	}

	*outputString = C.CString(s)

	return nil
}
