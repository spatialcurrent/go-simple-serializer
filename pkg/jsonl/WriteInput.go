// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package jsonl

import (
	"io"
)

// WriteInput provides the input for the Write function.
type WriteInput struct {
	Writer        io.Writer   // the underlying writer
	LineSeparator byte        // the newline byte
	Object        interface{} // the object to write
}
