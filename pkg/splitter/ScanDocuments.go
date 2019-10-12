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

// ScanDocuments returns a function that splits a stream of bytes on
// the given separator byte slice and whether to drop line-ending carriage returns.
// Returns a new bufio.SplitFunc compatible with bufio.Scanner
//
// Examples:
//	- ScanDocuments([]byte("---\n")[0], true) - split on YAML boundary marker and drop carriage returns at the end of a line.
func ScanDocuments(separator []byte, dropCR bool) bufio.SplitFunc {
	return bufio.SplitFunc(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}
		if i := bytes.Index(data, separator); i >= 0 {
			// We have a full separator-terminated line.
			if dropCR {
				return i + len(separator), DropCarriageReturn(data[0:i]), nil
			}
			return i + len(separator), data[0:i], nil
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
