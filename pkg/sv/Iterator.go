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

	"github.com/pkg/errors"
)

// Iterator is used to iterate over a table of separated values.
type Iterator struct {
	Reader *csv.Reader
	Type   reflect.Type
	header []interface{}
	limit  int
	count  int
}

// NewIteratorInput provides the input parameters for NewIterator function.
type NewIteratorInput struct {
	Reader     io.Reader
	Type       reflect.Type // required
	Separator  rune         // the values separator
	Header     []interface{}
	SkipLines  int
	Comment    string
	LazyQuotes bool
	Limit      int
}

// NewIterator returns a new iterator for iterating over a table of separated values.
func NewIterator(input *NewIteratorInput) (*Iterator, error) {

	if input.Type == nil || input.Type.Kind() != reflect.Map {
		return nil, errors.New("input type must be of kind map")
	}

	reader := csv.NewReader(input.Reader)
	reader.Comma = input.Separator
	reader.LazyQuotes = input.LazyQuotes
	reader.FieldsPerRecord = -1 // records may have a variable number of fields

	if len(input.Comment) > 1 {
		return nil, errors.New("go's encoding/csv package only supports single character comment characters")
	} else if len(input.Comment) == 1 {
		reader.Comment = []rune(input.Comment)[0]
	}

	for i := 0; i < input.SkipLines; i++ {
		if _, err := reader.Read(); err != nil {
			if err == io.EOF {
				return nil, err
			}
			return nil, errors.Wrap(err, "error skipping lines")
		}
	}

	header := input.Header
	if len(input.Header) == 0 {
		h, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				return nil, err
			}
			return nil, errors.Wrap(err, "error reading header")
		}
		header = make([]interface{}, 0, len(h))
		for _, str := range h {
			header = append(header, str)
		}
	}

	return &Iterator{Reader: reader, Type: input.Type, header: header, limit: input.Limit, count: 0}, nil
}

// Next reads from the underlying reader and returns the next object and error, if any.
// When finished, returns (nil, io.EOF).
func (it *Iterator) Next() (interface{}, error) {
	// If reached limit, return io.EOF
	if it.limit > 0 && it.count >= it.limit {
		return nil, io.EOF
	}

	// Increment Counter
	it.count++

	row, err := it.Reader.Read()
	if err != nil {
		if err == io.EOF {
			return nil, err
		}
		return nil, errors.Wrap(err, "error reading next line")
	}
	m := reflect.MakeMap(it.Type)
	for i, h := range it.header {
		if i < len(row) {
			m.SetMapIndex(reflect.ValueOf(h), reflect.ValueOf(row[i]))
		}
	}
	return m.Interface(), nil
}

func (it *Iterator) Header() []interface{} {
	return it.header
}
