// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

// Package contains a writer for writing objects as formatted lines by wrapping the "fmt" package in the standard library.
// This package has a minimal API and should not replace general use of the standard library package.
package fmt

import (
	"fmt"
	"io"
)

// Fprintf formats according to a format specifier and writes to w. It returns the number of bytes written and any write error encountered.
func Fprintf(w io.Writer, format string, a ...interface{}) (int, error) {
	return fmt.Fprintf(w, format, a...)
}
