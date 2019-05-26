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
	"github.com/hashicorp/hcl"
	hcl2 "github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hcl/hclsyntax"
	"github.com/pkg/errors"
)

import (
	"github.com/spatialcurrent/go-simple-serializer/pkg/bson"
	"github.com/spatialcurrent/go-simple-serializer/pkg/iterator"
	"github.com/spatialcurrent/go-simple-serializer/pkg/json"
	"github.com/spatialcurrent/go-simple-serializer/pkg/properties"
	"github.com/spatialcurrent/go-simple-serializer/pkg/toml"
	"github.com/spatialcurrent/go-simple-serializer/pkg/yaml"
)

import (
	"github.com/spatialcurrent/go-pipe/pkg/pipe"
)

// DeserializeReaderInput provides the input for the DeserializeReader function.
type DeserializeReaderInput struct {
	Reader        io.Reader
	Format        string
	Header        []string
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
	case "csv", "tsv", "jsonl":
		it, errorIterator := iterator.NewIterator(&iterator.NewIteratorInput{
			Reader:        input.Reader,
			Format:        input.Format,
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
	case "properties":
		return properties.Read(&properties.ReadInput{
			Type:            input.Type,
			Reader:          input.Reader,
			LineSeparator:   []byte(input.LineSeparator)[0],
			DropCR:          input.DropCR,
			Comment:         input.Comment,
			Trim:            input.Trim,
			UnescapeSpace:   true,
			UnescapeEqual:   true,
			UnescapeColon:   true,
			UnescapeNewLine: true,
		})
	case "bson", "hcl", "hcl2", "json", "toml", "yaml":
		b, err := ioutil.ReadAll(input.Reader)
		if err != nil {
			if err == io.EOF {
				return nil, io.EOF
			}
			return nil, errors.Wrap(err, "error reading bytes from reader")
		}
		switch input.Format {
		case "bson":
			return bson.UnmarshalType(b, input.Type)
		case "hcl":
			ptr := reflect.New(input.Type)
			ptr.Elem().Set(reflect.MakeMap(input.Type))
			obj, err := hcl.Parse(string(b))
			if err != nil {
				return nil, errors.Wrap(err, "Error parsing hcl")
			}
			if err := hcl.DecodeObject(ptr.Interface(), obj); err != nil {
				return nil, errors.Wrap(err, "Error decoding hcl")
			}
			return ptr.Elem().Interface(), nil
		case "hcl2":
			file, diags := hclsyntax.ParseConfig(b, "<stdin>", hcl2.Pos{Byte: 0, Line: 1, Column: 1})
			if diags.HasErrors() {
				return nil, errors.Wrap(errors.New(diags.Error()), "Error parsing hcl2")
			}
			return &file.Body, nil
		case "json":
			return json.UnmarshalType(b, input.Type)
		case "toml":
			return toml.UnmarshalType(b, input.Type)
		case "yaml":
			return yaml.UnmarshalType(b, input.Type)
		}
	}

	return nil, errors.Wrap(&ErrUnknownFormat{Name: input.Format}, "could not deserialize bytes")
}
