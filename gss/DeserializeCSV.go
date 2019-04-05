// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gss

import (
	"encoding/csv"
	"github.com/pkg/errors"
	"io"
	"reflect"
)

// DeserializeCSV deserializes a CSV or TSV string into a Go instance.
//  - https://golang.org/pkg/encoding/csv/
func DeserializeCSV(input io.Reader, format string, input_header []string, input_comment string, inputLazyQuotes bool, inputSkipLines int, inputLimit int, output_type reflect.Type) (interface{}, error) {

	if output_type.Kind() == reflect.Map {
		if inputLimit != 1 {
			return nil, errors.Wrap(&ErrInvalidLimit{Value: inputLimit}, "DeserializeCSV expects input limit of 1 when output type is of kind map.")
		}
		if len(input_header) == 0 {
			return nil, errors.New("deserializeCSV when returning a map type expects a input header")
		}
	} else if !(output_type.Kind() == reflect.Array || output_type.Kind() == reflect.Slice) {
		return nil, &ErrInvalidKind{Value: output_type.Kind(), Valid: []reflect.Kind{reflect.Array, reflect.Slice, reflect.Map}}
	}

	reader := csv.NewReader(input)
	if format == "tsv" {
		reader.Comma = '\t'
	}
	reader.LazyQuotes = inputLazyQuotes
	reader.FieldsPerRecord = -1 // records may have a variable number of fields

	if len(input_comment) > 1 {
		return nil, errors.New("go's encoding/csv package only supports single character comment characters")
	} else if len(input_comment) == 1 {
		reader.Comment = []rune(input_comment)[0]
	}

	if output_type.Kind() == reflect.Map {
		inRow, err := reader.Read()
		if err != nil {
			return nil, errors.Wrap(err, "error reading row from input with format csv")
		}
		if len(inRow) == 0 {
			return nil, &ErrEmptyRow{}
		}
		m := reflect.MakeMap(output_type)
		for i, h := range input_header {
			// if the number of columns in the header is greater than in the row,
			// then break and return data for the columns that are there
			if i >= len(inRow) {
				break
			}
			m.SetMapIndex(reflect.ValueOf(h), reflect.ValueOf(inRow[i]))
		}
		return m.Interface(), nil
	}

	if len(input_header) == 0 {
		h, err := reader.Read()
		if err != nil {
			if err != io.EOF {
				return nil, errors.Wrap(err, "Error reading header from input with format csv")
			}
		}
		input_header = h
	}

	output := reflect.MakeSlice(output_type, 0, 0)
	for {
		inRow, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, errors.Wrap(err, "error reading row from input with format csv")
			}
		}
		m := reflect.MakeMap(output_type.Elem())
		for i, h := range input_header {
			if i < len(inRow) {
				m.SetMapIndex(reflect.ValueOf(h), reflect.ValueOf(inRow[i]))
			}
			//else {
			//	m.SetMapIndex(reflect.ValueOf(h), "")
			//}
		}
		output = reflect.Append(output, m)

		if inputLimit > 0 && output.Len() >= inputLimit {
			break
		}
	}

	return output.Interface(), nil
}
