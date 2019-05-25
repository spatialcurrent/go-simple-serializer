// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package splitter

import (
	"bufio"
	"bytes"
)

// ScanLines returns a function that splits a stream of bytes on
// the given separator character and whether to drop line-ending carriage returns.
// Returns a new bufio.SplitFunc compatible with bufio.Scanner
//
// Examples:
//	- ScanLines([]byte("\n")[0], true) - split on new lines and drop carriage returns at the end of a line.
//	- ScanLines(byte(0), false) - split on null byte
func ScanLines(separator byte, dropCR bool) bufio.SplitFunc {
	return bufio.SplitFunc(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}
		if i := bytes.IndexByte(data, separator); i >= 0 {
			// We have a full separator-terminated line.
			if dropCR {
				return i + 1, DropCarriageReturn(data[0:i]), nil
			}
			return i + 1, data[0:i], nil
		}
		// If we're at EOF, we have a final, non-terminated line. Return it.
		if atEOF {
			if dropCR {
				return len(data), DropCarriageReturn(data), nil
			}
			return len(data), data, nil
		}
		// Request more data.
		return 0, nil, nil
	})
}
