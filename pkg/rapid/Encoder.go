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

	"github.com/spatialcurrent/go-simple-serializer/pkg/fit"
)

type Encoder struct {
	writer             io.Writer
	literalEncoder     *LiteralEncoder
	magicNumberWritten bool
	count              int
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		writer:             w,
		literalEncoder:     NewLiteralEncoder(w),
		magicNumberWritten: false,
		count:              0,
	}
}

func (e *Encoder) Reset(w io.Writer) {
	e.writer = w
	e.literalEncoder.Reset(w)
	e.magicNumberWritten = false
	e.count = 0
}

func (e *Encoder) Encode(v interface{}) error {
	if !e.magicNumberWritten {
		_, err := e.writer.Write(MagicNumber)
		if err != nil {
			return fmt.Errorf("error writing magic number: %w", err)
		}
		e.magicNumberWritten = true
	}
	err := e.literalEncoder.Encode(fit.Fit(v))
	if err != nil {
		return fmt.Errorf("error encoding value %#v: %w", v, err)
	}
	e.count += 1
	return nil
}

// Flush flushes the underlying writer, if it has a Flush method.
// This writer itself does no buffering.
func (e *Encoder) Flush() error {
	if flusher, ok := e.writer.(interface{ Flush() error }); ok {
		err := flusher.Flush()
		if err != nil {
			return fmt.Errorf("error flushing underlying writer: %w", err)
		}
	}
	return nil
}

// Close closes the underlying writer, if it has a Close method.
func (e *Encoder) Close() error {
	if closer, ok := e.writer.(interface{ Close() error }); ok {
		err := closer.Close()
		if err != nil {
			return fmt.Errorf("error closing underlying writer: %w", err)
		}
	}
	return nil
}
