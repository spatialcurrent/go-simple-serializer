// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gob

import (
	"encoding/gob"
	"io"
)

type Encoder struct {
	*gob.Encoder
}

// Encoder returns a new gob encoder given the underlying writer.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		Encoder: gob.NewEncoder(w),
	}
}
