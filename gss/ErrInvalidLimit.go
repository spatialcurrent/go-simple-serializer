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

// ErrInvalidLimit is used when an invalid limit int is given.
type ErrInvalidLimit struct {
	Value int // the value of the limit
}

// Error returns the error as a string.
func (e ErrInvalidLimit) Error() string {
	return fmt.Sprintf("invalid limit %d", e.Value)
}
