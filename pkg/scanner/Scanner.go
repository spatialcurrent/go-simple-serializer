// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

// Package scanner provides an interface{} for scanning through a stream of bytes.
// The interface is compatible with bufio.Scanner.
//
// Examples:
//	- https://godoc.org/bufio#Scanner
package scanner

import (
	"bufio"
	"io"

	"github.com/spatialcurrent/go-simple-serializer/pkg/splitter"
)

// Scanner is an interface compatible with bufio.Scanner that is used by iterators.
// By using this interface, we can support streams separated by null bytes, new-line characters, or any separator.
type Scanner interface {
	Scan() bool    // advanced the scanner to the next block
	Bytes() []byte // returns the bytes of the current block
	Text() string  // returns the text of the current block
}

// New returns a new Scanner that reads from the given reader,
// splits on the given newLine byte, and drops carriage returns if indicated.
func New(reader io.Reader, separator byte, dropCR bool) Scanner {
	scanner := bufio.NewScanner(reader)
	scanner.Split(splitter.ScanLines(separator, dropCR))
	return scanner
}
