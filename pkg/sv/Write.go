// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package sv

import (
	"fmt"
	"io"
	"reflect"

	"github.com/spatialcurrent/go-object/pkg/object"
	"github.com/spatialcurrent/go-pipe/pkg/pipe"
	"github.com/spatialcurrent/go-stringify/pkg/stringify"
)

// WriteInput provides the input for the Write function.
type WriteInput struct {
	Writer          io.Writer // the underlying writer
	Separator       rune      // the values separator
	Header          object.ObjectArray
	KeySerializer   stringify.Stringer
	ValueSerializer stringify.Stringer
	Object          object.Object // the object to write
	ExpandHeader    bool          // dynamically expand header, requires caching output in memory
	Sorted          bool          // sort columns
	Reversed        bool          // if sorted, sort in reverse alphabetical order.
	Limit           int
}

// Write writes the given object(s) as separated values, e.g., CSV or T
// If the type of the input object is of kind Array or Slice, then writes each object on its own line.
// Otherwise, just writes a CSV with a header and one row for the object.
func Write(input *WriteInput) error {

	concrete := input.Object.Concrete()

	// If the input is a slice of strings, simply write it as a row
	if slc, ok := concrete.Value().([]string); ok {
		errWriteTable := WriteTable(&WriteTableInput{
			Writer:    input.Writer,
			Separator: input.Separator,
			Header:    make([]string, 0),
			Rows:      [][]string{slc},
			Sorted:    input.Sorted, // if sorted and no specific wilcard position
			Reversed:  input.Reversed,
		})
		if errWriteTable != nil {
			return fmt.Errorf("error writing table to underlying writer: %w", errWriteTable)
		}
	}

	if input.ExpandHeader {

		// set the key serializer
		keySerializer := input.KeySerializer
		if keySerializer == nil {
			keySerializer = stringify.NewStringer("", false, false, false)
		}

		// set the value serializer
		valueSerializer := input.ValueSerializer
		if valueSerializer == nil {
			valueSerializer = stringify.NewStringer("", false, false, false)
		}

		// initialize header, wildcard, and known keys
		header := input.Header
		wildcard := header.Unique().Has("*")

		rows := make([][]string, 0)

		switch concrete.Kind() {
		case reflect.Map, reflect.Struct:

			if header.Empty() {
				header = concrete.Keys().Append(concrete.FieldNames().ObjectArray())
			} else if wildcard {
				header = header.Replace(
					"*",
					header.Append(concrete.Keys().Append(concrete.FieldNames().ObjectArray()).Unique().Subtract(header.Unique()).Array()).Value()...,
				)
			} else {
				// Get keys from object, subtract the known keys, and add the difference to the header
				header = header.Append(concrete.Keys().Append(concrete.FieldNames().ObjectArray()).Unique().Subtract(header.Unique()).Array())
			}

			// Create row from object
			row, errorRow := header.MapE(func(i int, v interface{}) (interface{}, error) {
				switch concrete.Kind() {
				case reflect.Map:
					return valueSerializer(concrete.Index(v).Value())
				case reflect.Struct:
					return valueSerializer(concrete.FieldByName(object.NewObject(v).String()).Value())
				}
				return nil, nil
			})
			if errorRow != nil {
				return fmt.Errorf("error serializing object as row: %w", errorRow)
			}
			rows = append(rows, row.StringArray().Value())

		case reflect.Array, reflect.Slice:
			if concrete.Empty() {
				// If there are no records then just return an empty string
				return nil
			}
			// create header
			for i := 0; i < concrete.Length() && (input.Limit < 0 || i <= input.Limit); i++ {
				x := object.NewObject(concrete.Index(i).Value()).Concrete()
				if header.Empty() {
					header = x.Keys().Append(x.FieldNames().ObjectArray())
				} else if wildcard {
					// replace wildcard with new columns
					header = header.Replace(
						"*",
						header.Append(x.Keys().Append(x.FieldNames().ObjectArray()).Unique().Subtract(header.Unique()).Array()).Value()...,
					)
				} else {
					// Get keys from object, subtract the known keys, and add the difference to the header
					header = header.Append(x.Keys().Append(x.FieldNames().ObjectArray()).Unique().Subtract(header.Unique()).Array())
				}
			}
			// create rows
			for i := 0; i < concrete.Length() && (input.Limit < 0 || i <= input.Limit); i++ {
				x := object.NewObject(concrete.Index(i).Value()).Concrete()
				// Create row from element
				row, errorRow := header.MapE(func(i int, v interface{}) (interface{}, error) {
					switch x.Kind() {
					case reflect.Map:
						return valueSerializer(x.Index(v).Value())
					case reflect.Struct:
						return valueSerializer(x.FieldByName(object.NewObject(v).String()).Value())
					}
					return nil, nil
				})
				if errorRow != nil {
					return fmt.Errorf("error serializing object as row: %w", errorRow)
				}
				rows = append(rows, row.StringArray().Value())
			}
		}

		// sort header
		if input.Sorted {
			header = header.Sort(input.Reversed)
		}

		outputHeader, errStringifyHeader := stringify.StringifySlice(header, keySerializer)
		if errStringifyHeader != nil {
			return fmt.Errorf("error stringifying header %v: %w", header.Value(), errStringifyHeader)
		}

		errWriteTable := WriteTable(&WriteTableInput{
			Writer:    input.Writer,
			Separator: input.Separator,
			Header:    outputHeader,
			Rows:      rows,
			Sorted:    input.Sorted && !wildcard, // if sorted and no specific wilcard position
			Reversed:  input.Reversed,
		})
		if errWriteTable != nil {
			return fmt.Errorf("error writing table to underlying writer: %w", errWriteTable)
		}
		return nil
	}

	// if streaming and not expanding header.
	p := pipe.NewBuilder().OutputLimit(input.Limit)
	if concrete.Kind() == reflect.Array || concrete.Kind() == reflect.Slice {
		it, errorIterator := pipe.NewSliceIterator(concrete.Value())
		if errorIterator != nil {
			return fmt.Errorf("error creating slice iterator: %w", errorIterator)
		}
		p = p.Input(it)
		p = p.Output(NewWriter(
			input.Writer,
			input.Separator,
			input.Header,
			input.KeySerializer,
			input.ValueSerializer,
			input.Sorted,
			input.Reversed,
		))
		errorRun := p.Run()
		if errorRun != nil {
			return fmt.Errorf("error serializing separated values: %w", errorRun)
		}
		return nil
	}

	w := NewWriter(
		input.Writer,
		input.Separator,
		input.Header,
		input.KeySerializer,
		input.ValueSerializer,
		input.Sorted,
		input.Reversed,
	)

	errorWrite := w.WriteObject(concrete.Value())
	if errorWrite != nil {
		return fmt.Errorf("error serializing separated values: %w", errorWrite)
	}

	errorFlush := w.Flush()
	if errorFlush != nil {
		return fmt.Errorf("error serializing separated values: %w", errorFlush)
	}

	return nil
}
