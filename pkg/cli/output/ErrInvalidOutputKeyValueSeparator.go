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

type ErrInvalidOutputKeyValueSeparator struct {
	Value string
}

func (e *ErrInvalidOutputKeyValueSeparator) Error() string {
	return fmt.Sprintf("invalid output key-value separator %q", e.Value)
}
