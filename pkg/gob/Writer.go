// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gob

import (
	"io"
	"reflect"

	"github.com/pkg/errors"

	pkgfit "github.com/spatialcurrent/go-fit/pkg/fit"
)

// Writer formats and writes objects to the underlying writer as gob-encoded items.
type Writer struct {
	writer  io.Writer // writer for the underlying stream
	encoder *Encoder  // GOB encoder
	fit     bool      // fit values before writing
}

// NewWriter returns a writer for formating and writing objets to the underlying writer as gob-encoded items.
func NewWriter(w io.Writer, fit bool) *Writer {
	return &Writer{
		writer:  w,
		encoder: NewEncoder(w),
		fit:     fit,
	}
}

// WriteObject formats and writes a single object to the underlying writer as a gob-encoded item.
func (w *Writer) WriteObject(obj interface{}) error {
	if w.fit {
		err := w.encoder.EncodeValue(pkgfit.FitValue(reflect.ValueOf(obj)))
		if err != nil {
			return errors.Wrap(err, "error writing to underlying writer")
		}
		return nil
	}
	err := w.encoder.EncodeValue(reflect.ValueOf(obj))
	if err != nil {
		return errors.Wrap(err, "error writing to underlying writer")
	}
	return nil
}

// WriteObjects formats and writes the given objects to the underlying writer.
func (w *Writer) WriteObjects(objects interface{}) error {
	value := reflect.ValueOf(objects)
	k := value.Type().Kind()
	if k == reflect.Ptr {
		value = value.Elem()
		k = value.Type().Kind()
	}
	if k == reflect.Array || k == reflect.Slice {
		for i := 0; i < value.Len(); i++ {
			err := w.WriteObject(value.Index(i).Interface())
			if err != nil {
				return errors.Wrap(err, "error writing object")
			}
		}
	}
	return nil
}

// Flush flushes the underlying writer, if it has a Flush method.
// This writer itself does no buffering.
func (w *Writer) Flush() error {
	if flusher, ok := w.writer.(interface{ Flush() error }); ok {
		err := flusher.Flush()
		if err != nil {
			return errors.Wrap(err, "error flushing underlying writer")
		}
	}
	return nil
}

// Close closes the underlying writer, if it has a Close method.
func (w *Writer) Close() error {
	if closer, ok := w.writer.(io.Closer); ok {
		err := closer.Close()
		if err != nil {
			return errors.Wrap(err, "error closing underlying writer")
		}
	}
	return nil
}
