// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gob

import (
	"bytes"
)

// Marshal formats an object into a slice of bytes of GOB.
func Marshal(obj interface{}, fit bool) ([]byte, error) {
	b := new(bytes.Buffer)
	e := NewEncoder(b, fit)
	err := e.Encode(obj)
	if err != nil {
		return make([]byte, 0), err
	}
	return b.Bytes(), nil
}
