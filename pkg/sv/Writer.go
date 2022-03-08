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

	"github.com/spatialcurrent/go-simple-serializer/pkg/inspector"
	"github.com/spatialcurrent/go-stringify/pkg/stringify"
)

type Writer struct {
	underlying      io.Writer
	writer          *csv.Writer
	columns         []interface{}
	headerWritten   bool
	keySerializer   stringify.Stringer
	valueSerializer stringify.Stringer
	sorted          bool
	reversed        bool
}

// NewWriter returns a new Writer for writing objects to an underlying writer formatted as separated values.
// NewWriter is a streaming writer, so cannot dynamically expand the header.
// To dynamically expand the header, then use the Write function with ExpandHeader set to true.
func NewWriter(underlying io.Writer, separator rune, columns []interface{}, keySerializer stringify.Stringer, valueSerializer stringify.Stringer, sorted bool, reversed bool) *Writer {

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
	fmt.Println("WriteHeader():", w.columns)
	w.headerWritten = true

	// Stringify columns into strings
	h, err := stringify.StringifySlice(w.columns, w.keySerializer)
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

func (w *Writer) ToRow(obj interface{}) ([]string, error) {
	return ToRow(obj, w.columns, w.valueSerializer)
}

func (w *Writer) WriteObject(obj interface{}) error {
	row := make([]string, 0)
	if slc, ok := obj.([]string); ok {
		row = slc
	} else {
		if !w.headerWritten {
			if len(w.columns) == 0 {
				inputObjectValue := reflect.ValueOf(obj)
				for reflect.TypeOf(inputObjectValue.Interface()).Kind() == reflect.Ptr {
					inputObjectValue = inputObjectValue.Elem()
				}
				inputObjectValue = reflect.ValueOf(inputObjectValue.Interface()) // sets value to concerete type
				inputObjectKind := inputObjectValue.Type().Kind()
				if inputObjectKind == reflect.Map {
					w.columns = inspector.GetKeysFromValue(inputObjectValue, w.sorted, w.reversed)
				} else if inputObjectKind == reflect.Struct {
					fieldNames := make([]interface{}, 0)
					for _, fieldName := range inspector.GetFieldNamesFromValue(inputObjectValue, w.sorted, w.reversed) {
						fieldNames = append(fieldNames, fieldName)
					}
					fmt.Println("Field Names:", fieldNames)
					w.columns = fieldNames
				}
			}
			if len(w.columns) == 0 {
				return fmt.Errorf("could not infer the header from the given value %#v with type %T", obj, obj)
			}
			err := w.WriteHeader()
			if err != nil {
				return fmt.Errorf("error writing header: %w", err)
			}
		}
		r, errorRow := w.ToRow(obj)
		if errorRow != nil {
			return fmt.Errorf("error serializing object as row: %w", errorRow)
		}
		row = r
	}
	errorWrite := w.writer.Write(row)
	if errorWrite != nil {
		return fmt.Errorf("error writing object: %w", errorWrite)
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
