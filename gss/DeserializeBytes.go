// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gss

import (
	"bytes"
	"fmt"
	"reflect"
)

import (
	"github.com/hashicorp/hcl"
	hcl2 "github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hcl/hclsyntax"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2/bson"
)

import (
	"github.com/spatialcurrent/go-simple-serializer/pkg/iterator"
	"github.com/spatialcurrent/go-simple-serializer/pkg/properties"
)

import (
	"github.com/spatialcurrent/go-pipe/pkg/pipe"
)

func deserializeBSON(inputBytes []byte, outputType reflect.Type) (interface{}, error) {
	if outputType.Kind() == reflect.Map {
		ptr := reflect.New(outputType)
		ptr.Elem().Set(reflect.MakeMap(outputType))
		err := bson.Unmarshal(inputBytes, ptr.Interface())
		if err != nil {
			return nil, errors.Wrap(err, "error unmarshalling bytes into BSON")
		}
		return ptr.Elem().Interface(), nil
	} else if outputType.Kind() == reflect.Slice {
		ptr := reflect.New(outputType)
		ptr.Elem().Set(reflect.MakeSlice(outputType, 0, 0))
		err := bson.Unmarshal(inputBytes, ptr.Interface())
		if err != nil {
			return nil, errors.Wrap(err, "error unmarshalling bytes into BSON")
		}
		return ptr.Elem().Interface(), nil
	}
	return nil, errors.New("Invalid output type for bson " + fmt.Sprint(outputType))
}

// DeserializeBytes reads in an object as string bytes and returns the representative Go instance.
func DeserializeBytes(input *DeserializeInput) (interface{}, error) {

	switch input.Format {
	case "csv", "tsv", "jsonl":
		it, errorIterator := iterator.NewIterator(&iterator.NewIteratorInput{
			Reader:        bytes.NewReader(input.Bytes),
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
			Reader:          bytes.NewReader(input.Bytes),
			LineSeparator:   []byte(input.LineSeparator)[0],
			DropCR:          input.DropCR,
			Comment:         input.Comment,
			Trim:            input.Trim,
			UnescapeSpace:   true,
			UnescapeEqual:   true,
			UnescapeColon:   true,
			UnescapeNewLine: true,
		})
	case "bson":
		return deserializeBSON(input.Bytes, input.Type)
	case "json":
		return DeserializeJSON(input.Bytes, input.Type)
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
	case "toml":
		return DeserializeTOML(string(input.Bytes), input.Type)
	case "yaml":
		return DeserializeYAML(input.Bytes, input.Type)
	}

	return nil, errors.Wrap(&ErrUnknownFormat{Name: input.Format}, "could not deserialize bytes")
}
