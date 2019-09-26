// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package input

import (
	"fmt"
	"strings"
)

type ErrInvalidInputFormat struct {
	Value    string
	Expected []string
}

func (e *ErrInvalidInputFormat) Error() string {
	return fmt.Sprintf("invalid input format %q, expecting on of %s", e.Value, strings.Join(e.Expected, ", "))
}
