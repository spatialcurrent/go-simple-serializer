// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package rapid

import (
	"bytes"
	"fmt"
)

func Marshal(object interface{}, limit int) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := Write(&WriteInput{
		Writer: buf,
		Limit:  limit,
	})
	if err != nil {
		return nil, fmt.Errorf("error writing json lines: %w", err)
	}
	return buf.Bytes(), nil
}
