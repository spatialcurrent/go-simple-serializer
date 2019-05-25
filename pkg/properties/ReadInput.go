// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package properties

import (
	"io"
	"reflect"
)

// ReadInput provides the input for the Read function.
type ReadInput struct {
	Type            reflect.Type // the output type
	Reader          io.Reader    // the underlying reader
	LineSeparator   byte         // the newline byte
	DropCR          bool         // drop carriage return
	Comment         string       // the comment prefix
	Trim            bool         // trim lines
	EscapePrefix    string       // escape prefix
	UnescapeSpace   bool         // unescape spaces
	UnescapeEqual   bool         // unescape =
	UnescapeColon   bool         // unescape :
	UnescapeNewLine bool         // unescape \n
}
