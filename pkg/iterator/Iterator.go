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
//	- github.com/spatialcurrent/go-simple-serializer/pkg/tags
package iterator

import (
	"io"
	"reflect"

	"github.com/pkg/errors"

	"github.com/spatialcurrent/go-simple-serializer/pkg/gob"
	"github.com/spatialcurrent/go-simple-serializer/pkg/jsonl"
	"github.com/spatialcurrent/go-simple-serializer/pkg/sv"
	"github.com/spatialcurrent/go-simple-serializer/pkg/tags"
)

var (
	ErrMissingLineSeparator = errors.New("missing line separator")
)

// Iterator is a simple interface that supports iterating over an input object source.
type Iterator interface {
	Next() (interface{}, error) // Returns the next object or error if any.  When input is exhausted, returns (nil, io.EOF).
}

// Input for NewIterator function.
type NewIteratorInput struct {
	Reader            io.Reader     // the underlying reader
	Format            string        // the format
	Header            []interface{} // for csv and tsv, the header.  If not given, then reads first line of stream as header.
	ScannerBufferSize int           // the initial buffer size for the scanner
	SkipLines         int           // Skip a given number of lines at the beginning of the stream.
	SkipBlanks        bool          // Skip blank lines.  If false, Next() returns a blank line as (nil, nil).  If true, Next() simply skips forward until it finds a non-blank line.
	SkipComments      bool          // Skip commented lines.  If false, Next() returns a commented line as (nil, nil).  If true, Next() simply skips forward until it finds a non-commented line.
	Comment           string        // The comment line prefix.  CSV and TSV only support single characters.  JSON Lines support any string.
	Trim              bool          // Trim each input line before parsing into an object.
	LazyQuotes        bool          // for csv and tsv, parse with lazy quotes
	Limit             int           // Limit the number of objects to read and return from the underlying stream.
	KeyValueSeparator string        // For tags, the key-value separator.
	LineSeparator     string        // For JSON Lines, the new line byte.
	DropCR            bool          // For JSON Lines, drop carriage returns at the end of lines.
	Type              reflect.Type  //
}

// NewIterator returns an Iterator for the given input source, format, and other options.
// Supports formats:
//	- csv - Comma-Separated Values
//	- jsonl - JSON Lines
//	- tags - Tags (key-value pairs)
//	- tsv - Tab-Separated Values
func NewIterator(input *NewIteratorInput) (Iterator, error) {

	if len(input.LineSeparator) == 0 {
		return nil, ErrMissingLineSeparator
	}

	switch input.Format {
	case "jsonl":
		it := jsonl.NewIterator(&jsonl.NewIteratorInput{
			Type:              input.Type,
			Reader:            input.Reader,
			ScannerBufferSize: input.ScannerBufferSize,
			SkipLines:         input.SkipLines,
			SkipBlanks:        input.SkipBlanks,
			SkipComments:      input.SkipComments,
			Comment:           input.Comment,
			Trim:              input.Trim,
			Limit:             input.Limit,
			LineSeparator:     []byte(input.LineSeparator)[0],
			DropCR:            input.DropCR,
		})
		return it, nil
	case "csv", "tags", "tsv":
		var inputType reflect.Type
		switch input.Format {
		case "csv":
			inputType := input.Type
			if inputType == nil {
				inputType = reflect.TypeOf(map[string]string{})
			}
			return sv.NewIterator(&sv.NewIteratorInput{
				Reader:     input.Reader,
				Type:       inputType,
				Separator:  ',',
				Header:     input.Header,
				SkipLines:  input.SkipLines,
				Comment:    input.Comment,
				LazyQuotes: input.LazyQuotes,
				Limit:      input.Limit,
			})
		case "tags":
			return tags.NewIterator(&tags.NewIteratorInput{
				Reader:            input.Reader,
				Type:              inputType,
				SkipLines:         input.SkipLines,
				SkipBlanks:        input.SkipBlanks,
				SkipComments:      input.SkipComments,
				Comment:           input.Comment,
				KeyValueSeparator: input.KeyValueSeparator,
				LineSeparator:     []byte(input.LineSeparator)[0],
				DropCR:            input.DropCR,
				Limit:             input.Limit,
			})
		case "tsv":
			return sv.NewIterator(&sv.NewIteratorInput{
				Reader:     input.Reader,
				Type:       inputType,
				Separator:  '\t',
				Header:     input.Header,
				SkipLines:  input.SkipLines,
				Comment:    input.Comment,
				LazyQuotes: input.LazyQuotes,
				Limit:      input.Limit,
			})
		}
	case "gob":
		it := gob.NewIterator(&gob.NewIteratorInput{
			Reader: input.Reader,
			Type:   input.Type,
			Limit:  input.Limit,
		})
		return it, nil
	}
	return nil, &ErrInvalidFormat{Format: input.Format}
}
