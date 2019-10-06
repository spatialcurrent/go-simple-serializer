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

type ErrInvalidInputComment struct {
	Value string
}

func (e *ErrInvalidInputComment) Error() string {
	return fmt.Sprintf("invalid input comment %q", e.Value)
}
