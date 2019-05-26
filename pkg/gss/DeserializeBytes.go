// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gss

import (
	"bytes"
	"reflect"
)

import (
	"github.com/hashicorp/hcl"
	hcl2 "github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hcl/hclsyntax"
	"github.com/pkg/errors"
)

import (
	"github.com/spatialcurrent/go-simple-serializer/pkg/iterator"
)

import (
	"github.com/spatialcurrent/go-pipe/pkg/pipe"
)

// DeserializeBytes reads in an object as string bytes and returns the representative Go instance.
func DeserializeBytes(input *DeserializeInput) (interface{}, error) {

	switch input.Format {
	case "csv", "tsv":
		it, errorIterator := iterator.NewIterator(&iterator.NewIteratorInput{
			Reader:        bytes.NewReader(input.Bytes),
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
	case "bson", "json", "jsonl", "properties", "toml", "yaml":
		s := NewSerializer(input.Format).Type(input.Type)
		if input.Format == "jsonl" {
			s = s.Limit(input.Limit)
		}
		if input.Format == "jsonl" || input.Format == "properties" {
			s = s.
				LineSeparator(input.LineSeparator).
				Comment(input.Comment).
				Trim(input.Trim)
		}
		if input.Format == "properties" {
			s = s.
				DropCR(input.DropCR).
				EscapePrefix(input.EscapePrefix).
				UnescapeSpace(input.UnescapeSpace).
				UnescapeColon(input.UnescapeColon).
				UnescapeNewLine(input.UnescapeNewLine).
				UnescapeEqual(input.UnescapeEqual)
		}
		return s.Deserialize(input.Bytes)
	case "hcl":
		ptr := reflect.New(input.Type)
		ptr.Elem().Set(reflect.MakeMap(input.Type))
		obj, err := hcl.Parse(string(input.Bytes))
		if err != nil {
			return nil, errors.Wrap(err, "Error parsing hcl")
		}
		if err := hcl.DecodeObject(ptr.Interface(), obj); err != nil {
			return nil, errors.Wrap(err, "Error decoding hcl")
		}
		return ptr.Elem().Interface(), nil
	case "hcl2":
		file, diags := hclsyntax.ParseConfig([]byte(input.Bytes), "<stdin>", hcl2.Pos{Byte: 0, Line: 1, Column: 1})
		if diags.HasErrors() {
			return nil, errors.Wrap(errors.New(diags.Error()), "Error parsing hcl2")
		}
		return &file.Body, nil
	}

	return nil, errors.Wrap(&ErrUnknownFormat{Name: input.Format}, "could not deserialize bytes")
}
