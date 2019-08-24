// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package toml

import (
	"bytes"

	bstoml "github.com/BurntSushi/toml"
	"github.com/pkg/errors" // import the BurntSushi toml library as bstoml
)

// Marshal formats an object into a slice of bytes of TOML.
func Marshal(obj interface{}) ([]byte, error) {
	if obj == nil {
		return make([]byte, 0), ErrNilObject
	}
	buf := new(bytes.Buffer)
	err := bstoml.NewEncoder(buf).Encode(obj)
	if err != nil {
		return make([]byte, 0), errors.Wrap(err, "error marshaling TOML bytes")
	}
	return buf.Bytes(), nil
}
