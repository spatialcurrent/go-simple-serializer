// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

// Package iterator provides an easy API to create an iterator to read objects from a file.
// Depends on the following packages in go-simple-serializer.
//	- github.com/spatialcurrent/go-simple-serializer/pkg/jsonl
//	- github.com/spatialcurrent/go-simple-serializer/pkg/sv
package iterator

import (
	"io"
	"reflect"
)

import (
	"github.com/spatialcurrent/go-simple-serializer/pkg/jsonl"
	"github.com/spatialcurrent/go-simple-serializer/pkg/sv"
)

// Iterator is a simple interface that supports iterating over an input object source.
type Iterator interface {
	Next() (interface{}, error) // Returns the next object or error if any.  When input is exhausted, returns (nil, io.EOF).
}

// Input for NewIterator function.
type NewIteratorInput struct {
	Reader        io.Reader    // the underlying reader
	Format        string       // the format
	Header        []interface{}     // for csv and tsv, the header.  If not given, then reads first line of stream as header.
	SkipLines     int          // Skip a given number of lines at the beginning of the stream.
	SkipBlanks    bool         // Skip blank lines.  If false, Next() returns a blank line as (nil, nil).  If true, Next() simply skips forward until it finds a non-blank line.
	SkipComments  bool         // Skip commented lines.  If false, Next() returns a commented line as (nil, nil).  If true, Next() simply skips forward until it finds a non-commented line.
	Comment       string       // The comment line prefix.  CSV and TSV only support single characters.  JSON Lines support any string.
	Trim          bool         // Trim each input line before parsing into an object.
	LazyQuotes    bool         // for csv and tsv, parse with lazy quotes
	Limit         int          // Limit the number of objects to read and return from the underlying stream.
	LineSeparator byte         // For JSON Lines, the new line byte.
	DropCR        bool         // For JSON Lines, drop carriage returns at the end of lines.
	Type          reflect.Type //
}

// NewIterator returns an Iterator for the given input source, format, and other options.
// Supports formats:
//	- csv - Comma-Separated Values
//	- jsonl - JSON Lines
//	- tsv - Tab-Separated Values
func NewIterator(input *NewIteratorInput) (Iterator, error) {
	if input.Format == "jsonl" {
		it := jsonl.NewIterator(&jsonl.NewIteratorInput{
			Reader:        input.Reader,
			SkipLines:     input.SkipLines,
			SkipBlanks:    input.SkipBlanks,
			SkipComments:  input.SkipComments,
			Comment:       input.Comment,
			Trim:          input.Trim,
			Limit:         input.Limit,
			LineSeparator: input.LineSeparator,
			DropCR:        input.DropCR,
		})
		return it, nil
	} else if input.Format == "csv" {
		return sv.NewIterator(&sv.NewIteratorInput{
			Reader:     input.Reader,
			Type:       input.Type.Elem(),
			Separator:  ',',
			Header:     input.Header,
			SkipLines:  input.SkipLines,
			Comment:    input.Comment,
			LazyQuotes: input.LazyQuotes,
			Limit:      input.Limit,
		})
	} else if input.Format == "tsv" {
		return sv.NewIterator(&sv.NewIteratorInput{
			Reader:     input.Reader,
			Type:       input.Type.Elem(),
			Separator:  '\t',
			Header:     input.Header,
			SkipLines:  input.SkipLines,
			Comment:    input.Comment,
			LazyQuotes: input.LazyQuotes,
			Limit:      input.Limit,
		})
	}
	return nil, &ErrInvalidFormat{Format: input.Format}
}
