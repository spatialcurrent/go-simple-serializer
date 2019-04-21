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

// DeserializeSV deserializes a CSV or TSV string into a Go instance.
//  - https://golang.org/pkg/encoding/csv/
func DeserializeSV(input io.Reader, format string, inputHeader []string, inputComment string, inputLazyQuotes bool, inputSkipLines int, inputLimit int, outputType reflect.Type) (interface{}, error) {

	if outputType.Kind() == reflect.Map {
		if inputLimit != 1 {
			return nil, errors.Wrap(&ErrInvalidLimit{Value: inputLimit}, "DeserializeSV expects input limit of 1 when output type is of kind map.")
		}
		if len(inputHeader) == 0 {
			return nil, errors.New("deserializeSV when returning a map type expects a input header")
		}
	} else if !(outputType.Kind() == reflect.Array || outputType.Kind() == reflect.Slice) {
		return nil, &ErrInvalidKind{Value: outputType.Kind(), Valid: []reflect.Kind{reflect.Array, reflect.Slice, reflect.Map}}
	}

	reader := csv.NewReader(input)
	if format == "tsv" {
		reader.Comma = '\t'
	}
	reader.LazyQuotes = inputLazyQuotes
	reader.FieldsPerRecord = -1 // records may have a variable number of fields

	if len(inputComment) > 1 {
		return nil, errors.New("go's encoding/csv package only supports single character comment characters")
	} else if len(inputComment) == 1 {
		reader.Comment = []rune(inputComment)[0]
	}

	if outputType.Kind() == reflect.Map {
		inRow, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				return nil, err
			}
			return nil, errors.Wrap(err, "error reading row into map from input with format csv")
		}
		if len(inRow) == 0 {
			return nil, &ErrEmptyRow{}
		}
		m := reflect.MakeMap(outputType)
		for i, h := range inputHeader {
			// if the number of columns in the header is greater than in the row,
			// then break and return data for the columns that are there
			if i >= len(inRow) {
				break
			}
			m.SetMapIndex(reflect.ValueOf(h), reflect.ValueOf(inRow[i]))
		}
		return m.Interface(), nil
	}

	if len(inputHeader) == 0 {
		h, err := reader.Read()
		if err != nil {
			if err != io.EOF {
				return nil, errors.Wrap(err, "Error reading header from input with format csv")
			}
		}
		inputHeader = h
	}

	output := reflect.MakeSlice(outputType, 0, 0)
	for {
		inRow, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, errors.Wrap(err, "error reading row into slice from input with format csv")
			}
		}
		m := reflect.MakeMap(outputType.Elem())
		for i, h := range inputHeader {
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
