// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gss

import (
	"io"
	"io/ioutil"
	"reflect"
)

import (
	"github.com/pkg/errors"
)

import (
	"github.com/spatialcurrent/go-simple-serializer/pkg/iterator"
	"github.com/spatialcurrent/go-simple-serializer/pkg/serializer"
)

import (
	"github.com/spatialcurrent/go-pipe/pkg/pipe"
)

// DeserializeReaderInput provides the input for the DeserializeReader function.
type DeserializeReaderInput struct {
	Reader        io.Reader
	Format        string
	Header        []interface{}
	Comment       string
	LazyQuotes    bool
	SkipLines     int
	SkipBlanks    bool
	SkipComments  bool
	Trim          bool
	Limit         int
	LineSeparator string
	DropCR        bool
	Type          reflect.Type
	Async         bool
	Verbose       bool
}

// DeserializeReader reads the serialized object from an io.Reader and returns the representative Go instance.
func DeserializeReader(input *DeserializeReaderInput) (interface{}, error) {

	switch input.Format {
	case "csv", "tsv", "jsonl", "tags":
		// These formats can be streamed.
		it, errorIterator := iterator.NewIterator(&iterator.NewIteratorInput{
			Reader:        input.Reader,
			Type:          input.Type,
			Format:        input.Format,
			Header:        input.Header,
			Comment:       input.Comment,
			SkipLines:     input.SkipLines,
			SkipBlanks:    input.SkipBlanks,
			SkipComments:  input.SkipComments,
			LazyQuotes:    input.LazyQuotes,
			Trim:          input.Trim,
			Limit:         input.Limit,
			LineSeparator: []byte(input.LineSeparator)[0],
			DropCR:        input.DropCR,
		})
		if errorIterator != nil {
			return nil, errors.Wrap(errorIterator, "error creating iterator")
		}
		w := pipe.NewSliceWriterWithValues(reflect.MakeSlice(input.Type, 0, 0).Interface())
		errorRun := pipe.NewBuilder().Input(it).Output(w).Run()
		if errorRun != nil {
			return w.Values(), errors.Wrap(errorRun, "error deserializing")
		}
		return w.Values(), nil
	case "bson", "hcl", "hcl2", "json", "properties", "toml", "yaml":
		// These formats do not support streaming.
		b, err := ioutil.ReadAll(input.Reader)
		if err != nil {
			if err == io.EOF {
				return nil, io.EOF
			}
			return nil, errors.Wrap(err, "error reading bytes from reader")
		}

		// Set up Serializer
		s := serializer.New(input.Format).Type(input.Type)
		if input.Format == "properties" || input.Format == "yaml" {
			s = s.Comment(input.Comment)
		}
		if input.Format == "properties" {
			s = s.
				LineSeparator(input.LineSeparator).
				DropCR(input.DropCR).
				Trim(input.Trim).
				UnescapeSpace(true).
				UnescapeEqual(true).
				UnescapeColon(true).
				UnescapeNewLine(true)
		}

		// Deserialize bytes into object
		obj, err := serializer.New(input.Format).Type(input.Type).Deserialize(b)
		if err != nil {
			return nil, errors.Wrap(err, "error deserializing object")
		}
		return obj, nil
	}

	return nil, errors.Wrap(&ErrUnknownFormat{Name: input.Format}, "could not deserialize bytes")
}
