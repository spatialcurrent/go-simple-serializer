// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gss

import (
	"fmt"
	"reflect"
	"strings"
)

import (
	"github.com/hashicorp/hcl"
	hcl2 "github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hcl/hclsyntax"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2/bson"
)

func unescapePropertyText(in string) string {
	out := in
	out = strings.Replace(fmt.Sprint(out), "\\ ", " ", -1)
	out = strings.Replace(fmt.Sprint(out), "\\\\", "\\", -1)
	return out
}

func deserializeBSON(input_bytes []byte, outputType reflect.Type) (interface{}, error) {
	if outputType.Kind() == reflect.Map {
		ptr := reflect.New(outputType)
		ptr.Elem().Set(reflect.MakeMap(outputType))
		err := bson.Unmarshal(input_bytes, ptr.Interface())
		if err != nil {
			return nil, errors.Wrap(err, "error unmarshalling bytes into BSON")
		}
		return ptr.Elem().Interface(), nil
	} else if outputType.Kind() == reflect.Slice {
		ptr := reflect.New(outputType)
		ptr.Elem().Set(reflect.MakeSlice(outputType, 0, 0))
		err := bson.Unmarshal(input_bytes, ptr.Interface())
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
	case "csv", "tsv":
		return DeserializeCSV(string(input.Bytes), input.Format, input.Header, input.Comment, input.LazyQuotes, input.SkipLines, input.Limit, input.Type)
	case "properties":
		return DeserializeProperties(string(input.Bytes), input.Comment, input.Type)
	case "bson":
		return deserializeBSON(input.Bytes, input.Type)
	case "json":
		return DeserializeJSON(input.Bytes, input.Type)
	case "jsonl":
		return DeserializeJSONL(string(input.Bytes), input.Comment, input.SkipLines, input.Limit, input.Type, input.Async)
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
