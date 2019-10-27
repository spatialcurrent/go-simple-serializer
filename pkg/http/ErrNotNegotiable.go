// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package http

import (
	"fmt"
)

// ErrNotNegotiable is used when the server cannot negotiate a format given an accept header.
type ErrNotNegotiable struct {
	Value string // the name of the unknown format
}

// Error returns the error formatted as a string.
func (e ErrNotNegotiable) Error() string {
	return fmt.Sprintf("could not negotiate format from string %q", e.Value)
}
