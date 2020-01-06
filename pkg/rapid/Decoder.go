// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package rapid

import (
	"fmt"
	"io"
)

type Decoder struct {
	literalDecoder *LiteralDecoder
	count          int
	eof            bool
}

//func NewDecoder(r io.Reader, separator byte, dropCR bool) *Decoder {
func NewDecoder(r io.ByteReader) *Decoder {
	return &Decoder{
		literalDecoder: NewLiteralDecoder(r),
		count:          0,
		eof:            false,
	}
}

func (d *Decoder) Reset(r io.ByteReader) {
	d.literalDecoder.Reset(r)
	d.count = 0
}

func (d *Decoder) Decode(v interface{}) error {

	if d.eof {
		return io.EOF
	}

	err := d.literalDecoder.Decode(v)
	if err != nil {
		d.eof = true
		if err == io.EOF {
			return io.EOF
		}
		return fmt.Errorf("error decoding object %#v after decoding %d: %w", v, d.count, err)
	}

	d.count += 1

	return nil
}
