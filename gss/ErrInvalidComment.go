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

// ErrInvalidComment is used when an invalid comment string is given.
type ErrInvalidComment struct {
	Value string // the value of the comment
}

// Error returns the error as a string.
func (e ErrInvalidComment) Error() string {
	return fmt.Sprintf("invalid comment %s", e.Value)
}
