// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gob

import (
	stdgob "encoding/gob"
	"io"
	"reflect"
)

// Iterator iterates trough a stream of bytes
// returning a new object on each call of Next()
// until it reaches the end and returns io.EOF.
type Iterator struct {
	Type    reflect.Type    // the type to unmarshal for each item
	Decoder *stdgob.Decoder // the scanner that splits the underlying stream of bytes
	Limit   int             // Limit the number of objects to read and return from the underlying stream.
	Count   int             // The current count of the number of objects read.
}

// NewIteratorInput provides the input parameters for the NewIterator function.
type NewIteratorInput struct {
	Reader io.Reader
	Type   reflect.Type // the type to unmarshal for each line
	Limit  int          // Limit the number of objects to read and return from the underlying stream.
}

// NewIterator returns a new gob iterator base on the given input.
func NewIterator(input *NewIteratorInput) *Iterator {

	return &Iterator{
		Type:    input.Type,
		Decoder: stdgob.NewDecoder(input.Reader),
		Limit:   input.Limit,
		Count:   0,
	}
}

// Next reads from the underlying reader and returns the next object and error, if any.
// If a blank line is found and SkipBlanks is false, then returns (nil, nil).
// If a commented line is found and SkipComments is false, then returns (nil, nil).
// When the input stream is exhausted, returns (nil, io.EOF).
func (it *Iterator) Next() (interface{}, error) {

	if it.Type == nil {
		return nil, ErrMissingType
	}

	// If reached limit, return io.EOF
	if it.Limit > 0 && it.Count >= it.Limit {
		return nil, io.EOF
	}

	// Increment Counter
	it.Count++

	ptr := reflect.New(it.Type)
	if it.Type.Kind() == reflect.Map {
		ptr.Elem().Set(reflect.MakeMap(it.Type))
	}

	err := it.Decoder.Decode(ptr.Interface())
	if err != nil {
		if err == io.EOF {
			return nil, io.EOF
		}
		return nil, err
	}

	return ptr.Elem().Interface(), nil
}
