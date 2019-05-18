// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package sv

import (
	"io"
)

// WriteInput provides the input for the Write function.
type WriteInput struct {
	Writer    io.Writer  // the underlying writer
	Separator rune       // the values separator
	Header    []string   // the header row
	Rows      [][]string // the row of values to write to the underlying writer
}
