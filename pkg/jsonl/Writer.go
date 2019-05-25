// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package jsonl

import (
	"io"
	"reflect"
)

import (
	"github.com/spatialcurrent/go-simple-serializer/pkg/json"
)

import (
	"github.com/pkg/errors"
)

// Writer formats and writes objects to the underlying writer as JSON Lines (aka jsonl).
type Writer struct {
	writer    io.Writer // writer for the underlying stream
	separator byte      // the separator byte to use, e.g, null byte or \n.
	pretty    bool      // write pretty output
}

// NewWriter returns a writer for formating and writing objets to the underlying writer as JSON Lines (aka jsonl).
func NewWriter(w io.Writer, separator byte, pretty bool) *Writer {
	return &Writer{
		writer:    w,
		separator: separator,
		pretty:    pretty,
	}
}

// WriteObject formats and writes a single object to the underlying writer as JSON
// and appends the writer's line separator.
func (w *Writer) WriteObject(obj interface{}) error {
	b, err := json.Marshal(obj, w.pretty)
	if err != nil {
		return errors.Wrap(err, "error marshaling object")
	}
	_, err = w.writer.Write(append(b, w.separator))
	if err != nil {
		return errors.Wrap(err, "error writing to underlying writer")
	}
	return nil
}

// WriteObjects formats and writes the given objets to the underlying writer as JSON lines
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
	if flusher, ok := w.writer.(Flusher); ok {
		err := flusher.Flush()
		if err != nil {
			return errors.Wrap(err, "error flushing underlying writer")
		}
	}
	return nil
}
