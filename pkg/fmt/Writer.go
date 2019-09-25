// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package fmt

import (
	"io"
	"reflect"

	"github.com/pkg/errors"
)

// Writer formats and writes objects to the underlying writer as formatted lines.
type Writer struct {
	writer        io.Writer // writer for the underlying stream
	format        string    // the format string
	lineSeparator string    // the separator stirng to use, e.g, null byte or \n.
}

// NewWriter returns a writer for formating and writing objets to the underlying writer as JSON Lines (aka jsonl).
func NewWriter(w io.Writer, format string, lineSeparator string) *Writer {
	return &Writer{
		writer:        w,
		format:        format,
		lineSeparator: lineSeparator,
	}
}

// WriteObject formats and writes a single object to the underlying writer as a formatted line
// and appends the writer's line separator.
func (w *Writer) WriteObject(obj interface{}) error {
	format := w.format
	if len(w.lineSeparator) > 0 {
		format += w.lineSeparator
	}
	_, err := Fprintf(w.writer, format, obj)
	if err != nil {
		return errors.Wrap(err, "error writing to underlying writer")
	}
	return nil
}

// WriteObjects formats and writes the given objets to the underlying writer as formatted lines
// and separates the objects using the writer's line separator.
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
