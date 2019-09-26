// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package input

import (
	"fmt"
)

type ErrInvalidInputLineSeparator struct {
	Value string
}

func (e *ErrInvalidInputLineSeparator) Error() string {
	return fmt.Sprintf("invalid input line separator %q", e.Value)
}
