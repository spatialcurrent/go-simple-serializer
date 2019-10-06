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

type ErrMissingOutputFormat struct {
	Expected []string
}

func (e *ErrMissingOutputFormat) Error() string {
	return fmt.Sprintf("missing output format, expecting on of %s", strings.Join(e.Expected, ", "))
}
