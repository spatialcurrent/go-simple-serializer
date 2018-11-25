// =================================================================
//
// Copyright (C) 2018 Spatial Current, Inc. - All Rights Reserved
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
func DeserializeBytes(input []byte, format string, inputHeader []string, inputComment string, inputLazyQuotes bool, inputSkipLines int, inputLimit int, outputType reflect.Type, verbose bool) (interface{}, error) {

	if format == "csv" || format == "tsv" {
		return DeserializeCSV(string(input), format, inputHeader, inputComment, inputLazyQuotes, inputSkipLines, inputLimit, outputType)
	} else if format == "properties" {
		return DeserializeProperties(string(input), inputComment, outputType)
	} else if format == "bson" {
		return deserializeBSON(input, outputType)
	} else if format == "json" {
		return DeserializeJSON(input, outputType)
	} else if format == "jsonl" {
		return DeserializeJSONL(string(input), inputComment, inputSkipLines, inputLimit, outputType)
	} else if format == "hcl" {
		ptr := reflect.New(outputType)
		ptr.Elem().Set(reflect.MakeMap(outputType))
		obj, err := hcl.Parse(string(input))
		if err != nil {
			return nil, errors.Wrap(err, "Error parsing hcl")
		}
		if err := hcl.DecodeObject(ptr.Interface(), obj); err != nil {
			return nil, errors.Wrap(err, "Error decoding hcl")
		}
		return ptr.Elem().Interface(), nil
	} else if format == "hcl2" {
		file, diags := hclsyntax.ParseConfig([]byte(input), "<stdin>", hcl2.Pos{Byte: 0, Line: 1, Column: 1})
		if diags.HasErrors() {
			return nil, errors.Wrap(errors.New(diags.Error()), "Error parsing hcl2")
		}
		return &file.Body, nil
	} else if format == "toml" {
		return DeserializeTOML(string(input), outputType)
	} else if format == "yaml" {
		return DeserializeYAML(input, outputType)
	}

	return nil, errors.Wrap(&ErrUnknownFormat{Name: format}, "could not deserialize bytes")
}
