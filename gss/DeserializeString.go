// =================================================================
//
// Copyright (C) 2018 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gss

import (
	"reflect"
)

// DeserializeString reads in an object as a string and returns the representative Go instance.
func DeserializeString(input string, format string, input_header []string, input_comment string, input_lazy_quotes bool, input_limit int, output_type reflect.Type, verbose bool) (interface{}, error) {

	return DeserializeBytes(
		[]byte(input),
		format,
		input_header,
		input_comment,
		input_lazy_quotes,
		input_limit,
		output_type,
		verbose)
}
