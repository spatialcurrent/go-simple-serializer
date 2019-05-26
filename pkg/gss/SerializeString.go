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
func SerializeString(input *SerializeInput) (string, error) {
	if input.Format == "golang" || input.Format == "go" {
		return fmt.Sprint(input.Object), nil
	}
	b, err := SerializeBytes(input)
	return string(b), err
}
