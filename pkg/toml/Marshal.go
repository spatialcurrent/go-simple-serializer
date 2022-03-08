// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package toml

import (
	"bytes"
	"fmt"

	bstoml "github.com/BurntSushi/toml"
)

// Marshal formats an object into a slice of bytes of TOML.
func Marshal(obj interface{}) ([]byte, error) {
	if obj == nil {
		return make([]byte, 0), ErrNilObject
	}
	buf := new(bytes.Buffer)
	err := bstoml.NewEncoder(buf).Encode(obj)
	if err != nil {
		return make([]byte, 0), fmt.Errorf("error marshaling TOML bytes: %w", err)
	}
	return buf.Bytes(), nil
}
