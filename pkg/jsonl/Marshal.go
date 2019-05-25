// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package jsonl

import (
	"bytes"
)

import (
	"github.com/pkg/errors"
)

import (
	"github.com/spatialcurrent/go-stringify/pkg/stringify"
)

func Marshal(object interface{}, lineSeparator string, keySerializer stringify.Stringer, pretty bool, limit int) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := Write(&WriteInput{
		Writer:        buf,
		LineSeparator: lineSeparator,
		KeySerializer: keySerializer,
		Object:        object,
		Pretty:        pretty,
		Limit:         limit,
	})
	if err != nil {
		return nil, errors.Wrap(err, "error writing json lines")
	}
	return buf.Bytes(), nil
}
