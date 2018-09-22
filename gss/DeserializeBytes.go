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

func deserializeBSON(input_bytes []byte, output_type reflect.Type) (interface{}, error) {
	if output_type.Kind() == reflect.Map {
		ptr := reflect.New(output_type)
		ptr.Elem().Set(reflect.MakeMap(output_type))
		err := bson.Unmarshal(input_bytes, ptr.Interface())
		if err != nil {
			return nil, errors.Wrap(err, "error unmarshalling bytes into BSON")
		}
		return ptr.Elem().Interface(), nil
	} else if output_type.Kind() == reflect.Slice {
		ptr := reflect.New(output_type)
		ptr.Elem().Set(reflect.MakeSlice(output_type, 0, 0))
		err := bson.Unmarshal(input_bytes, ptr.Interface())
		if err != nil {
			return nil, errors.Wrap(err, "error unmarshalling bytes into BSON")
		}
		return ptr.Elem().Interface(), nil
	}
	return nil, errors.New("Invalid output type for bson " + fmt.Sprint(output_type))
}

// DeserializeBytes reads in an object as string bytes and returns the representative Go instance.
func DeserializeBytes(input []byte, format string, input_header []string, input_comment string, input_lazy_quotes bool, input_limit int, output_type reflect.Type, verbose bool) (interface{}, error) {

	if format == "csv" || format == "tsv" {
		return DeserializeCSV(string(input), format, input_header, input_comment, input_lazy_quotes, input_limit, output_type)
	} else if format == "properties" {
		return DeserializeProperties(string(input), input_comment, output_type)
	} else if format == "bson" {
		return deserializeBSON(input, output_type)
	} else if format == "json" {
		return DeserializeJSON(input, output_type)
	} else if format == "jsonl" {
		return DeserializeJSONL(string(input), input_comment, input_limit, output_type)
	} else if format == "hcl" {
		ptr := reflect.New(output_type)
		ptr.Elem().Set(reflect.MakeMap(output_type))
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
		return DeserializeTOML(string(input), output_type)
	} else if format == "yaml" {
		return DeserializeYAML(input, output_type)
	}

	return nil, nil
}
