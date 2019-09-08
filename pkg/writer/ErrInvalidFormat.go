// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package writer

import (
	"fmt"
)

// ErrInvalidFormat is used when an invalid format is provided.
type ErrInvalidFormat struct {
	Format string // the name of the invalid format
}

// Error returns the error as a string.
func (e ErrInvalidFormat) Error() string {
	return fmt.Sprintf("invalid format %q, expecting csv, go, gob, jsonl, or tsv", e.Format)
}
