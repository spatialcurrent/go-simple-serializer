// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gob

import (
	"bytes"
	"reflect"

	pkgfit "github.com/spatialcurrent/go-fit/pkg/fit"
)

// Marshal formats an object into a slice of gob-encoded bytes.
func Marshal(obj interface{}, fit bool) ([]byte, error) {
	b := new(bytes.Buffer)
	e := NewEncoder(b)
	if fit {
		err := e.EncodeValue(pkgfit.FitValue(reflect.ValueOf(obj)))
		if err != nil {
			return make([]byte, 0), err
		}
	} else {
		err := e.Encode(obj)
		if err != nil {
			return make([]byte, 0), err
		}
	}
	return b.Bytes(), nil
}
