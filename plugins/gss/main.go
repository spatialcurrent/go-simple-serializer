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
	"fmt"
	"strings"
)

//"github.com/spatialcurrent/go-stringify/pkg/stringify"

import (
	"github.com/spatialcurrent/go-simple-serializer/pkg/gss"
)

var gitBranch string
var gitCommit string

func main() {}

//lint:ignore ST1020
//export Version
func Version() *C.char {
	var b strings.Builder
	if len(gitBranch) > 0 {
		b.WriteString(fmt.Sprintf("Branch: %q\n", gitBranch))
	}
	if len(gitCommit) > 0 {
		b.WriteString(fmt.Sprintf("Commit: %q\n", gitCommit))
	}
	return C.CString(b.String())
}

//lint:ignore ST1020
//export Convert
func Convert(inputString *C.char, inputFormat *C.char, outputFormat *C.char, outputPretty *C.char, outputSorted *C.char, outputString **C.char) *C.char {

	s, err := gss.Convert(&gss.ConvertInput{
		InputBytes:  []byte(C.GoString(inputString)),
		InputFormat: C.GoString(inputFormat),
		//InputHeader:     stringify.StringSliceToInterfaceSlice(strings.Split(C.GoString(inputHeader), ",")),
		InputHeader:         []interface{}{},
		InputComment:        "",
		InputLazyQuotes:     false,
		InputSkipLines:      0,
		InputLimit:          -1,
		InputLineSeparator:  "\n",
		InputEscapePrefix:   "//",
		OutputFormat:        C.GoString(outputFormat),
		OutputHeader:        []interface{}{},
		OutputLimit:         -1,
		OutputPretty:        C.GoString(outputFormat) == "1",
		OutputSorted:        C.GoString(outputSorted) == "1",
		OutputEscapePrefix:  "//",
		OutputLineSeparator: "\n",
	})
	if err != nil {
		return C.CString(err.Error())
	}

	*outputString = C.CString(string(s))

	return nil
}
