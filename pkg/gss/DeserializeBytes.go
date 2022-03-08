// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gss

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"reflect"

	"github.com/hashicorp/hcl"

	"github.com/spatialcurrent/go-pipe/pkg/pipe"
	"github.com/spatialcurrent/go-simple-serializer/pkg/iterator"
	"github.com/spatialcurrent/go-simple-serializer/pkg/serializer"
)

// DeserializeBytesInput provides the input for the DeserializeBytes function.
type DeserializeBytesInput struct {
	Bytes             []byte
	Format            string
	Header            []interface{}
	Comment           string
	LazyQuotes        bool
	ScannerBufferSize int
	SkipLines         int
	SkipBlanks        bool
	SkipComments      bool
	Trim              bool
	Limit             int
	LineSeparator     string
	DropCR            bool
	Type              reflect.Type
	EscapePrefix      string
	UnescapeSpace     bool
	UnescapeNewLine   bool
	UnescapeColon     bool
	UnescapeEqual     bool
}

// DeserializeBytes reads in an object as string bytes and returns the representative Go instance.
func DeserializeBytes(input *DeserializeBytesInput) (interface{}, error) {

	switch input.Format {
	case "csv", "tsv", "jsonl", "geojsonl", "tags":
		it, errorIterator := iterator.NewIterator(&iterator.NewIteratorInput{
			Reader:            bytes.NewReader(input.Bytes),
			Type:              input.Type,
			Format:            input.Format,
			Header:            input.Header,
			Comment:           input.Comment,
			ScannerBufferSize: input.ScannerBufferSize,
			SkipLines:         input.SkipLines,
			SkipBlanks:        input.SkipBlanks,
			SkipComments:      input.SkipComments,
			LazyQuotes:        input.LazyQuotes,
			Trim:              input.Trim,
			Limit:             input.Limit,
			LineSeparator:     input.LineSeparator,
			DropCR:            input.DropCR,
		})
		if errorIterator != nil {
			return nil, fmt.Errorf("error creating iterator: %w", errorIterator)
		}
		w := pipe.NewSliceWriterWithValues(reflect.MakeSlice(input.Type, 0, 0).Interface())
		errorRun := pipe.NewBuilder().Input(it).Output(w).Run()
		if errorRun != nil {
			return w.Values(), fmt.Errorf("error deserializing: %w", errorRun)
		}
		return w.Values(), nil
	case "bson", "json", "properties", "toml", "yaml":
		s := serializer.New(input.Format).Type(input.Type)
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
		return s.Deserialize(input.Bytes)
	case "gob":
		obj := make([]interface{}, 0)
		d := gob.NewDecoder(bytes.NewReader(input.Bytes))
		err := d.Decode(&obj)
		return obj, err
	case "hcl":
		ptr := reflect.New(input.Type)
		ptr.Elem().Set(reflect.MakeMap(input.Type))
		obj, err := hcl.Parse(string(input.Bytes))
		if err != nil {
			return nil, fmt.Errorf("Error parsing hcl: %w", err)
		}
		if err := hcl.DecodeObject(ptr.Interface(), obj); err != nil {
			return nil, fmt.Errorf("Error decoding hcl: %w", err)
		}
		return ptr.Elem().Interface(), nil
	}

	return nil, fmt.Errorf("could not deserialize bytes: %w", &ErrUnknownFormat{Name: input.Format})
}
