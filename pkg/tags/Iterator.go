// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package tags

import (
	"io"
	"reflect"
	"strings"
	"unicode/utf8"

	"github.com/pkg/errors"

	"github.com/spatialcurrent/go-simple-serializer/pkg/scanner"
)

// Iterator iterates trough a stream of bytes
// returning a new object on each call of Next()
// until it reaches the end and returns io.EOF.
type Iterator struct {
	Scanner           scanner.Scanner // the scanner that splits the underlying stream of bytes
	Type              reflect.Type
	KeyValueSeparator rune   // the key value separator
	Comment           string // The comment line prefix.  Can be any string.
	SkipBlanks        bool   // Skip blank lines.  If false, Next() returns a blank line as (nil, nil).  If true, Next() simply skips forward until it finds a non-blank line.
	SkipComments      bool   // Skip commented lines.  If false, Next() returns a commented line as (nil, nil).  If true, Next() simply skips forward until it finds a non-commented line.
	Limit             int    // Limit the number of objects to read and return from the underlying stream.
	Count             int    // The current count of the number of objects read.
}

// NewIteratorInput provides the input parameters for the NewIterator function.
type NewIteratorInput struct {
	Reader            io.Reader
	Type              reflect.Type
	SkipLines         int    // Skip a given number of lines at the beginning of the stream.
	SkipBlanks        bool   // Skip blank lines.  If false, Next() returns a blank line as (nil, nil).  If true, Next() simply skips forward until it finds a non-blank line.
	SkipComments      bool   // Skip commented lines.  If false, Next() returns a commented line as (nil, nil).  If true, Next() simply skips forward until it finds a non-commented line.
	Comment           string // The comment line prefix. Can be any string.
	Limit             int    // Limit the number of objects to read and return from the underlying stream.
	KeyValueSeparator string // the key value separator
	LineSeparator     byte   // The new line byte.
	DropCR            bool   // Drop carriage returns at the end of lines.
}

// NewIterator returns a new JSON Lines (aka jsonl) Iterator base on the given input.
func NewIterator(input *NewIteratorInput) (*Iterator, error) {
	if len(input.KeyValueSeparator) == 0 {
		return nil, ErrMissingKeyValueSeparator
	}

	KeyValueSeparator, n := utf8.DecodeRuneInString(input.KeyValueSeparator)
	if KeyValueSeparator == utf8.RuneError && n == 1 {
		return nil, errors.Wrap(ErrInvalidUTF8, "error decoding key-value separator")
	}

	s := scanner.New(input.Reader, input.LineSeparator, input.DropCR)
	for i := 0; i < input.SkipLines; i++ {
		if !s.Scan() {
			break
		}
	}

	it := &Iterator{
		Scanner:           s,
		Type:              input.Type,
		KeyValueSeparator: KeyValueSeparator,
		Comment:           input.Comment,
		SkipBlanks:        input.SkipBlanks,
		SkipComments:      input.SkipComments,
		Limit:             input.Limit,
		Count:             0,
	}

	return it, nil
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

	if it.Scanner.Scan() {
		line := strings.TrimSpace(it.Scanner.Text())
		if len(line) == 0 {
			if it.SkipBlanks {
				return it.Next()
			}
			return nil, nil
		}
		if len(it.Comment) > 0 && strings.HasPrefix(line, it.Comment) {
			if it.SkipComments {
				return it.Next()
			}
			return nil, nil
		}
		if it.Type != nil {
			obj, err := UnmarshalType([]byte(line), it.KeyValueSeparator, it.Type)
			if err != nil {
				return obj, errors.Wrap(err, "error unmarshaling next tags object")
			}
			return obj, nil
		}
		obj, err := Unmarshal([]byte(line), it.KeyValueSeparator)
		if err != nil {
			return obj, errors.Wrap(err, "error unmarshaling next tags object")
		}
		return obj, nil
	}
	return nil, io.EOF
}
