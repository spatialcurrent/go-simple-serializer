// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package properties

import (
	"io"
)

// WriteInput provides the input for the Write function.
type WriteInput struct {
	Writer            io.Writer                                // the underlying writer
	LineSeparator     string                                   // the newline byte
	KeyValueSeparator string                                   // the separator for key-value pairs
	Object            interface{}                              // the object to write
	ValueSerializer   func(object interface{}) (string, error) // serializer for object properties
	Sorted            bool                                     // sort output
	EscapePrefix      string                                   // escape prefix, if empty then doesn't escape
	EscapeSpace       bool                                     // escape spaces
	EscapeEqual       bool                                     // escape =
	EscapeColon       bool                                     // escape :
	EscapeNewLine     bool                                     // escape \n
}
