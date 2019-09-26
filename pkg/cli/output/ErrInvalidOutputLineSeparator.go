// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package output

import (
	"fmt"
)

type ErrInvalidOutputLineSeparator struct {
	Value string
}

func (e *ErrInvalidOutputLineSeparator) Error() string {
	return fmt.Sprintf("invalid output line separator %q", e.Value)
}
