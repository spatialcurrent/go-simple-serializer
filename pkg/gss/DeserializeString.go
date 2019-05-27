// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gss

import (
	"reflect"
)

// DeserializeString reads in an object as a string and returns the representative Go instance.
func DeserializeString(inputString string, inputFormat string, inputHeader []interface{}, inputComment string, inputLazyQuotes bool, inputSkipLines int, inputLineSeparator string, inputLimit int, outputType reflect.Type, async bool, verbose bool) (interface{}, error) {

	return DeserializeBytes(&DeserializeInput{
		Bytes:         []byte(inputString),
		Format:        inputFormat,
		Header:        inputHeader,
		Comment:       inputComment,
		LazyQuotes:    inputLazyQuotes,
		SkipLines:     inputSkipLines,
		LineSeparator: inputLineSeparator,
		Limit:         inputLimit,
		Type:          outputType,
		Async:         async,
		Verbose:       verbose,
	})

}
