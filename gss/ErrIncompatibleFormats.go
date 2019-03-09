// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gss

import (
	"fmt"
)

// ErrIncompatibleFormats is used when an input format and output format are incompatible.
type ErrIncompatibleFormats struct {
	Input  string // the name of the input format
	Output string // the name of the output format
}

// Error returns the error as a string.
func (e ErrIncompatibleFormats) Error() string {
	return fmt.Sprintf("incompatible formats %s input with %s output", e.Input, e.Output)
}
