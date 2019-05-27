// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package sv

import (
	"encoding/csv"
	"io"
	"reflect"
)

import (
	"github.com/pkg/errors"
)

import (
	"github.com/spatialcurrent/go-simple-serializer/pkg/inspector"
	"github.com/spatialcurrent/go-stringify/pkg/stringify"
)

type Writer struct {
	underlying      io.Writer
	writer          *csv.Writer
	columns         []interface{}
	headerWritten   bool
	valueSerializer func(object interface{}) (string, error)
}

// WriteSV writes the given rows as separated values.
func NewWriter(underlying io.Writer, separator rune, columns []interface{}, valueSerializer func(object interface{}) (string, error)) *Writer {

	// Create a new CSV writer.
	csvWriter := csv.NewWriter(underlying)

	// set the values separator
	csvWriter.Comma = separator

	if valueSerializer == nil {
		valueSerializer = stringify.DefaultValueStringer("")
	}

	return &Writer{
		underlying:      underlying,
		writer:          csvWriter,
		columns:         columns,
		headerWritten:   false,
		valueSerializer: valueSerializer,
	}
}

func (w *Writer) WriteHeader() error {
	w.headerWritten = true

	// Stringify columns into strings
	h, err := stringify.StringifySlice(w.columns, w.valueSerializer)
	if err != nil {
		return errors.Wrap(err, "error stringifying columns")
	}

	// Write header to writer
	err = w.writer.Write(h)
	if err != nil {
		return errors.Wrap(err, "error writing header")
	}
	return nil
}

func (w *Writer) ToRow(obj interface{}) ([]string, error) {
	return ToRow(obj, w.columns, w.valueSerializer)
}

func (w *Writer) WriteObject(obj interface{}) error {
	if !w.headerWritten {
		if len(w.columns) == 0 {
			objectValue := reflect.ValueOf(obj)
			objectType := objectValue.Type()
			objectKind := objectType.Kind()
			if objectKind == reflect.Ptr {
				objectValue = objectValue.Elem()
				objectType = objectValue.Type()
				objectKind = objectType.Kind()
			}
			if objectKind == reflect.Map {
				w.columns = inspector.GetKeysFromValue(objectValue, false)
			} else if objectKind == reflect.Struct {
				fieldNames := make([]interface{}, 0)
				for _, fieldName := range inspector.GetFieldNamesFromValue(objectValue, false) {
					fieldNames = append(fieldNames, fieldName)
				}
				w.columns = fieldNames
			}
		}
		if len(w.columns) == 0 {
			return errors.New("could not infer the header from the given value")
		}
		err := w.WriteHeader()
		if err != nil {
			return errors.Wrap(err, "error writing header")
		}
	}
	row, errorRow := w.ToRow(obj)
	if errorRow != nil {
		return errors.Wrap(errorRow, "error serializing object as row")
	}
	errorWrite := w.writer.Write(row)
	if errorWrite != nil {
		return errors.Wrap(errorWrite, "error writing object")
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
	w.writer.Flush()
	err := w.writer.Error()
	if err != nil {
		return err
	}
	if flusher, ok := w.underlying.(Flusher); ok {
		err := flusher.Flush()
		if err != nil {
			return errors.Wrap(err, "error flushing underlying writer")
		}
	}
	return nil
}
