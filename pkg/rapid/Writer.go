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
	"reflect"
)

// Writer formats and writes objects to the underlying writer as JSON Lines (aka jsonl).
type Writer struct {
	encoder *Encoder
}

// NewWriter returns a writer for formating and writing objets to the underlying writer as JSON Lines (aka jsonl).
func NewWriter(w io.Writer) *Writer {
	return &Writer{
		encoder: NewEncoder(w),
	}
}

// WriteObject formats and writes a single object to the underlying writer as JSON
// and appends the writer's line separator.
func (w *Writer) WriteObject(obj interface{}) error {
	return w.encoder.Encode(obj)
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
				return fmt.Errorf("error writing object: %w", err)
			}
		}
	}
	return nil
}

// Flush flushes the underlying writer, if it has a Flush method.
// This writer itself does no buffering.
func (w *Writer) Flush() error {
	return w.encoder.Flush()
}

// Close closes the underlying writer, if it has a Close method.
func (w *Writer) Close() error {
	return w.encoder.Close()
}
