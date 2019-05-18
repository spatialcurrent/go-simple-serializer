// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package sv

import (
	"fmt"
)

// ErrInvalidFormat indicates an invalid format was given
type ErrInvalidFormat struct {
	Format string
}

// Error returns the string representation of the error.
func (e *ErrInvalidFormat) Error() string {
	return fmt.Sprintf("invalid format %q, expecting \"csv\" or \"tsv\"", e.Format)
}
