// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package sv

import (
	"encoding/csv"
	"fmt"
	"io"
	"reflect"

	"github.com/spatialcurrent/go-object/pkg/object"
	"github.com/spatialcurrent/go-stringify/pkg/stringify"
)

type Writer struct {
	underlying      io.Writer
	writer          *csv.Writer
	columns         object.ObjectArray
	headerWritten   bool
	keySerializer   stringify.Stringer
	valueSerializer stringify.Stringer
	sorted          bool
	reversed        bool
}

// NewWriter returns a new Writer for writing objects to an underlying writer formatted as separated values.
// NewWriter is a streaming writer, so cannot dynamically expand the header.
// To dynamically expand the header, then use the Write function with ExpandHeader set to true.
func NewWriter(underlying io.Writer, separator rune, columns object.ObjectArray, keySerializer stringify.Stringer, valueSerializer stringify.Stringer, sorted bool, reversed bool) *Writer {

	// Create a new CSV writer.
	csvWriter := csv.NewWriter(underlying)

	// set the values separator
	csvWriter.Comma = separator

	if keySerializer == nil {
		keySerializer = stringify.NewStringer("", false, false, false)
	}

	if valueSerializer == nil {
		valueSerializer = stringify.NewStringer("", false, false, false)
	}

	return &Writer{
		underlying:      underlying,
		writer:          csvWriter,
		columns:         columns,
		headerWritten:   false,
		keySerializer:   keySerializer,
		valueSerializer: valueSerializer,
		sorted:          sorted,
		reversed:        reversed,
	}
}

func (w *Writer) WriteHeader() error {
	w.headerWritten = true

	// Stringify columns into strings
	h, err := stringify.StringifySlice(w.columns.Value(), w.keySerializer)
	if err != nil {
		return fmt.Errorf("error stringifying columns: %w", err)
	}

	// Write header to writer
	err = w.writer.Write(h)
	if err != nil {
		return fmt.Errorf("error writing header: %w", err)
	}
	return nil
}

func (w *Writer) WriteObject(obj interface{}) error {

	if slc, ok := obj.([]string); ok {
		errorWrite := w.writer.Write(slc)
		if errorWrite != nil {
			return fmt.Errorf("error writing object: %w", errorWrite)
		}
		return nil
	}

	concrete := object.NewObject(obj).Concrete()

	if !w.headerWritten {
		if w.columns.Empty() {
			w.columns = concrete.Keys().Append(concrete.FieldNames().ObjectArray())
		}
		if w.columns.Empty() {
			return fmt.Errorf("could not infer the header from the given value %#v with type %T", obj, obj)
		}
		err := w.WriteHeader()
		if err != nil {
			return fmt.Errorf("error writing header: %w", err)
		}
	}

	row, errorRow := w.columns.MapE(func(i int, v interface{}) (interface{}, error) {
		switch concrete.Kind() {
		case reflect.Map:
			return w.valueSerializer(concrete.Index(v).Value())
		case reflect.Struct:
			return w.valueSerializer(concrete.FieldByName(object.NewObject(v).String()).Value())
		}
		return nil, nil
	})
	if errorRow != nil {
		return fmt.Errorf("error serializing object as row: %w", errorRow)
	}

	errorWrite := w.writer.Write(row.StringArray().Value())
	if errorWrite != nil {
		return fmt.Errorf("error writing object: %w", errorWrite)
	}

	return nil
}

func (w *Writer) WriteObjects(objects interface{}) error {
	value := object.NewObject(objects).Concrete()
	if value.Kind() == reflect.Array || value.Kind() == reflect.Slice {
		for i := 0; i < value.Length(); i++ {
			err := w.WriteObject(value.Index(i).Value())
			if err != nil {
				return fmt.Errorf("error writing object: %w", err)
			}
		}
	}
	return nil
}

func (w *Writer) Flush() error {
	w.writer.Flush()
	err := w.writer.Error()
	if err != nil {
		return err
	}
	if flusher, ok := w.underlying.(Flusher); ok {
		err := flusher.Flush()
		if err != nil {
			return fmt.Errorf("error flushing underlying writer: %w", err)
		}
	}
	return nil
}

// Close closes the underlying writer, if it has a Close method.
func (w *Writer) Close() error {
	if closer, ok := w.underlying.(io.Closer); ok {
		err := closer.Close()
		if err != nil {
			return fmt.Errorf("error closing underlying writer: %w", err)
		}
	}
	return nil
}
