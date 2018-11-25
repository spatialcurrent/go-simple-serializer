// =================================================================
//
// Copyright (C) 2018 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gss

import (
	"encoding/csv"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"reflect"
	"strings"
)

// DeserializeCSV deserializes a CSV or TSV string into a Go instance.
//  - https://golang.org/pkg/encoding/csv/
func DeserializeCSV(input string, format string, input_header []string, input_comment string, input_lazy_quotes bool, inputSkipLines int, input_limit int, output_type reflect.Type) (interface{}, error) {

	if output_type.Kind() == reflect.Map {
		if input_limit != 1 {
			return nil, errors.New("deserializeCSV when returning a map type expects input_limit to be set to 1 but got " + fmt.Sprint(input_limit))
		}
		if len(input_header) == 0 {
			return nil, errors.New("deserializeCSV when returning a map type expects a input header")
		}
	} else if !(output_type.Kind() == reflect.Array || output_type.Kind() == reflect.Slice) {
		return nil, errors.New("deserializeCSV expects an array, map, or slice type but got " + fmt.Sprint(output_type))
	}

	reader := csv.NewReader(strings.NewReader(input))
	if format == "tsv" {
		reader.Comma = '\t'
	}
	reader.LazyQuotes = input_lazy_quotes

	if len(input_comment) > 1 {
		return nil, errors.New("go's encoding/csv package only supports single character comment characters")
	} else if len(input_comment) == 1 {
		reader.Comment = []rune(input_comment)[0]
	}

	if output_type.Kind() == reflect.Map {
		inRow, err := reader.Read()
		if err != nil {
			return nil, errors.Wrap(err, "Error reading row from input with format csv")
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
				return nil, errors.Wrap(err, "Error reading row from input with format csv")
			}
		}
		m := reflect.MakeMap(output_type.Elem())
		for i, h := range input_header {
			m.SetMapIndex(reflect.ValueOf(h), reflect.ValueOf(inRow[i]))
		}
		output = reflect.Append(output, m)
	}

	return output.Interface(), nil
}
