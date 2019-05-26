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
func DeserializeString(input string, inputFormat string, inputHeader []string, inputComment string, inputLazyQuotes bool, inputSkipLines int, inputLimit int, outputType reflect.Type, async bool, verbose bool) (interface{}, error) {

	return DeserializeBytes(&DeserializeInput{
		Bytes:      []byte(input),
		Format:     inputFormat,
		Header:     inputHeader,
		Comment:    inputComment,
		LazyQuotes: inputLazyQuotes,
		SkipLines:  inputSkipLines,
		Limit:      inputLimit,
		Type:       outputType,
		Async:      async,
		Verbose:    verbose,
	})

}
