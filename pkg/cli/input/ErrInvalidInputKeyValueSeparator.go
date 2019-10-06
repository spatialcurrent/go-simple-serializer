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

type ErrInvalidInputKeyValueSeparator struct {
	Value string
}

func (e *ErrInvalidInputKeyValueSeparator) Error() string {
	return fmt.Sprintf("invalid input key-value separator %q", e.Value)
}
