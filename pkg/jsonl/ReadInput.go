// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package jsonl

import (
	"io"
	"reflect"
)

// ReadInput provides the input for the Read function.
type ReadInput struct {
	Type          reflect.Type // the output type
	Reader        io.Reader    // the underlying reader
	SkipLines     int
	SkipBlanks    bool
	SkipComments  bool
	Comment       string // the comment prefix
	Trim          bool   // trim lines
	LineSeparator byte   // the newline byte
	DropCR        bool   // drop carriage return
	Limit         int
}
