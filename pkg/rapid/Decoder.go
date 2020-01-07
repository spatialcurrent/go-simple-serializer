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
	"io"
)

type Decoder struct {
	reader         io.ByteReader
	literalDecoder *LiteralDecoder
	count          int
	eof            bool
}

//func NewDecoder(r io.Reader, separator byte, dropCR bool) *Decoder {
func NewDecoder(r io.ByteReader) *Decoder {
	return &Decoder{
		reader:         r,
		literalDecoder: NewLiteralDecoder(r),
		count:          0,
		eof:            false,
	}
}

func (d *Decoder) Reset(r io.ByteReader) {
	d.reader = r
	d.literalDecoder.Reset(r)
	d.count = 0
}

func (d *Decoder) Decode(v interface{}) error {

	if d.eof {
		return io.EOF
	}

	if d.count == 0 {
		h := make([]byte, 0, len(MagicNumber))
		for i := 0; i < len(MagicNumber); i++ {
			b, err := d.reader.ReadByte()
			if err != nil {
				return fmt.Errorf("error reading magic number byte %d: %w", i, err)
			}
			h = append(h, b)
		}
		if !bytes.Equal(h, MagicNumber) {
			return fmt.Errorf("invalid magic number, expecting \"% x\" but found \"% x\"", MagicNumber, h)
		}
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
