// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gss

import (
	"encoding/gob"
	"io"
	"io/ioutil"
	"reflect"

	"github.com/pkg/errors"

	"github.com/spatialcurrent/go-pipe/pkg/pipe"
	"github.com/spatialcurrent/go-simple-serializer/pkg/iterator"
	"github.com/spatialcurrent/go-simple-serializer/pkg/serializer"
)

// DeserializeReaderInput provides the input for the DeserializeReader function.
type DeserializeReaderInput struct {
	Reader          io.Reader
	Format          string
	Header          []interface{}
	Comment         string
	LazyQuotes      bool
	SkipLines       int
	SkipBlanks      bool
	SkipComments    bool
	Trim            bool
	Limit           int
	LineSeparator   string
	DropCR          bool
	Type            reflect.Type
	EscapePrefix    string
	UnescapeSpace   bool
	UnescapeNewLine bool
	UnescapeColon   bool
	UnescapeEqual   bool
}

// DeserializeReader reads the serialized object from an io.Reader and returns the representative Go instance.
func DeserializeReader(input *DeserializeReaderInput) (interface{}, error) {

	switch input.Format {
	case "csv", "tsv", "jsonl", "geojsonl", "tags":

		var iteratorType reflect.Type
		if input.Type != nil {
			iteratorType = input.Type.Elem()
		}
		// These formats can be streamed.
		it, errorIterator := iterator.NewIterator(&iterator.NewIteratorInput{
			Reader:        input.Reader,
			Type:          iteratorType,
			Format:        input.Format,
			Header:        input.Header,
			Comment:       input.Comment,
			SkipLines:     input.SkipLines,
			SkipBlanks:    input.SkipBlanks,
			SkipComments:  input.SkipComments,
			LazyQuotes:    input.LazyQuotes,
			Trim:          input.Trim,
			Limit:         input.Limit,
			LineSeparator: input.LineSeparator,
			DropCR:        input.DropCR,
		})
		if errorIterator != nil {
			return nil, errors.Wrap(errorIterator, "error creating iterator")
		}
		p := pipe.NewBuilder().Input(it)
		var w *pipe.SliceWriter
		if input.Type != nil {
			w = pipe.NewSliceWriterWithValues(reflect.MakeSlice(input.Type, 0, 0).Interface())
			p = p.Output(w)
		} else {
			w = pipe.NewSliceWriterWithValues([]interface{}{})
			p = p.Output(w)
		}
		errorRun := p.Run()
		if errorRun != nil {
			return w.Values(), errors.Wrap(errorRun, "error deserializing")
		}
		return w.Values(), nil
	case "gob":
		obj := make([]interface{}, 0)
		d := gob.NewDecoder(input.Reader)
		err := d.Decode(&obj)
		return obj, err
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
				Comment(input.Comment).
				Trim(input.Trim).
				DropCR(input.DropCR).
				EscapePrefix(input.EscapePrefix).
				UnescapeSpace(input.UnescapeSpace).
				UnescapeColon(input.UnescapeColon).
				UnescapeNewLine(input.UnescapeNewLine).
				UnescapeEqual(input.UnescapeEqual)
		}

		// Deserialize bytes into object
		obj, err := s.Deserialize(b)
		if err != nil {
			return nil, errors.Wrap(err, "error deserializing object")
		}
		return obj, nil
	}

	return nil, errors.Wrap(&ErrUnknownFormat{Name: input.Format}, "could not deserialize bytes")
}
