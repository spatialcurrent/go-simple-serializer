// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package sv

import (
	"io"
	"reflect"
)

// ReadInput provides the input for the Read function.
type ReadInput struct {
	Type       reflect.Type // the output type
	Reader     io.Reader
	Separator  rune // the values separator
	Header     []string
	SkipLines  int
	Comment    string
	LazyQuotes bool
	Limit      int
}
