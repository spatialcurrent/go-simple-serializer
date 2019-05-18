// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package jsonl

import (
	"encoding/json"
	"io"
	"reflect"
)

import (
	"github.com/pkg/errors"
)

type Writer struct {
	writer  io.Writer
	newline byte
}

// WriteSV writes the given rows as separated values.
func NewWriter(w io.Writer) *Writer {
	return &Writer{
		writer:  w,
		newline: []byte("\n")[0],
	}
}

func (w *Writer) WriteObject(obj interface{}) error {
	b, err := json.Marshal(obj)
	if err != nil {
		return errors.Wrap(err, "error marshaling object")
	}
	_, err = w.writer.Write(append(b, w.newline))
	if err != nil {
		return errors.Wrap(err, "error writing to underlying writer")
	}
	return nil
}

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

func (w *Writer) Flush() error {
	if flusher, ok := w.writer.(Flusher); ok {
		err := flusher.Flush()
		if err != nil {
			return errors.Wrap(err, "error flushing underlying writer")
		}
	}
	return nil
}
