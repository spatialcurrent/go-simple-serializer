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

	"github.com/pkg/errors"

	"github.com/spatialcurrent/go-simple-serializer/pkg/json"
	"github.com/spatialcurrent/go-stringify/pkg/stringify"
)

// Writer formats and writes objects to the underlying writer as JSON Lines (aka jsonl).
type Writer struct {
	writer        io.Writer // writer for the underlying stream
	separator     string    // the separator stirng to use, e.g, null byte or \n.
	keySerializer stringify.Stringer
	pretty        bool // write pretty output
}

// NewWriter returns a writer for formating and writing objets to the underlying writer as JSON Lines (aka jsonl).
func NewWriter(w io.Writer, separator string, keySerializer stringify.Stringer, pretty bool) *Writer {
	return &Writer{
		writer:        w,
		separator:     separator,
		keySerializer: keySerializer,
		pretty:        pretty,
	}
}

// WriteObject formats and writes a single object to the underlying writer as JSON
// and appends the writer's line separator.
func (w *Writer) WriteObject(obj interface{}) error {
	obj, err := stringify.StringifyMapKeys(obj, w.keySerializer)
	if err != nil {
		return errors.Wrap(err, "error stringify map keys")
	}
	b, err := json.Marshal(obj, w.pretty)
	if err != nil {
		return errors.Wrap(err, "error marshaling object")
	}
	if len(w.separator) > 0 {
		b = append(b, []byte(w.separator)...)
	}
	_, err = w.writer.Write(b)
	if err != nil {
		return errors.Wrap(err, "error writing to underlying writer")
	}
	return nil
}

// WriteObjects formats and writes the given objects to the underlying writer as JSON lines
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
