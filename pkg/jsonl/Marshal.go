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

func Marshal(object interface{}, lineSeparator string, pretty bool) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := Write(&WriteInput{
		Writer:        buf,
		LineSeparator: lineSeparator,
		Object:        object,
		Pretty:        pretty,
	})
	if err != nil {
		return nil, errors.Wrap(err, "error writing json lines")
	}
	return buf.Bytes(), nil
}
