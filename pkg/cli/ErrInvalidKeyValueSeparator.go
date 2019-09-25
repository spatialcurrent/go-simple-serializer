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

type ErrInvalidKeyValueSeparator struct {
	Value string
}

func (e *ErrInvalidKeyValueSeparator) String() {
	return fmt.Sprintf("invalid key-value separator %q", e.Value)
}
