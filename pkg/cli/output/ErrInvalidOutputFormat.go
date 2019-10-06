// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package output

import (
	"fmt"
	"strings"
)

type ErrInvalidOutputFormat struct {
	Value    string
	Expected []string
}

func (e *ErrInvalidOutputFormat) Error() string {
	return fmt.Sprintf("invalid output format %q, expecting on of %s", e.Value, strings.Join(e.Expected, ", "))
}
