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

type ErrMissingInputFormat struct {
	Expected []string
}

func (e *ErrMissingInputFormat) Error() string {
	return fmt.Sprintf("missing input format, expecting on of %s", strings.Join(e.Expected, ", "))
}
