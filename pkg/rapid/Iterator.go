// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package rapid

import (
	"io"
	"reflect"
)

// Iterator iterates trough a stream of bytes
// returning a new object on each call of Next()
// until it reaches the end and returns io.EOF.
type Iterator struct {
	Type    reflect.Type // the type to unmarshal for each line
	Decoder *Decoder     // the decoder that splits the underlying stream of bytes
	Limit   int          // Limit the number of objects to read and return from the underlying stream.
	Count   int          // The current count of the number of objects read.
}

// NewIteratorInput provides the input parameters for the NewIterator function.
type NewIteratorInput struct {
	Type   reflect.Type // the type to unmarshal for each line
	Reader io.ByteReader
	Limit  int  // Limit the number of objects to read and return from the underlying stream.
	DropCR bool // Drop carriage returns at the end of lines.
}

// NewIterator returns a new JSON Lines (aka jsonl) Iterator base on the given input.
func NewIterator(input *NewIteratorInput) *Iterator {
	return &Iterator{
		Type:    input.Type,
		Decoder: NewDecoder(input.Reader),
		Limit:   input.Limit,
		Count:   0,
	}
}

// Next reads from the underlying reader and returns the next object and error, if any.
// If a blank line is found and SkipBlanks is false, then returns (nil, nil).
// If a commented line is found and SkipComments is false, then returns (nil, nil).
// When the input stream is exhausted, returns (nil, io.EOF).
func (it *Iterator) Next() (interface{}, error) {

	// If reached limit, return io.EOF
	if it.Limit > 0 && it.Count >= it.Limit {
		return nil, io.EOF
	}

	// Increment Counter
	it.Count++

	m := map[string]interface{}{}

	err := it.Decoder.Decode(&m)
	if err != nil {
		return nil, err
	}
	return m, err
}
