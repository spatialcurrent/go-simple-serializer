// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gss

import (
	"fmt"
)

// SerializeString serializes an object to its representation given by format.
func SerializeString(input interface{}, format string, header []string, limit int) (string, error) {
	if format == "golang" || format == "go" {
		return fmt.Sprint(input), nil
	}
	b, err := SerializeBytes(input, format, header, limit)
	return string(b), err
}
