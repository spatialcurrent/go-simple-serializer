// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package cli

import (
	"fmt"
)

type ErrInvalidLineSeparator struct {
	Value string
}

func (e *ErrInvalidLineSeparator) String() {
	return fmt.Sprintf("invalid line separator %q", e.Value)
}
