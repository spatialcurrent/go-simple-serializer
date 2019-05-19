// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

// Package iterator provides an easy API to create a iterator to read objects from a file.
package iterator

import (
	"io"
)

import (
	"github.com/spatialcurrent/go-simple-serializer/pkg/jsonl"
	"github.com/spatialcurrent/go-simple-serializer/pkg/sv"
)

type Iterator interface {
	Next() (interface{}, error)
}

// Input for NewIterator function.
type NewIteratorInput struct {
	Reader       io.Reader
	Format       string
	Header       []string
	SkipLines    int
	SkipBlanks   bool
	SkipComments bool
	Comment      string
	Trim         bool
	LazyQuotes   bool
	Limit        int
}

func NewIterator(input *NewIteratorInput) (Iterator, error) {
	if input.Format == "jsonl" {
		it := jsonl.NewIterator(&jsonl.NewIteratorInput{
			Reader:       input.Reader,
			SkipLines:    input.SkipLines,
			SkipBlanks:   input.SkipBlanks,
			SkipComments: input.SkipComments,
			Comment:      input.Comment,
			Trim:         input.Trim,
			Limit:        input.Limit,
		})
		return it, nil
	} else if input.Format == "csv" {
		return sv.NewIterator(&sv.NewIteratorInput{
			Reader:     input.Reader,
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
