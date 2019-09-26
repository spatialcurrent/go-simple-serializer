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

type ErrInvalidInputType struct {
	Value string
}

func (e *ErrInvalidInputType) Error() string {
	return fmt.Sprintf("invalid input type %q", e.Value)
}
