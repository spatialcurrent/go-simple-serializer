// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

// Package gss provides simple functions for serializing/deserializing objects into common formats.
//
// Usage
//
// The simplest usage of gss is to call the DeserializeBytes, DeserializeString, SerializeBytes, and SerializeString functions.
//
//  inputObject, err := gss.DeserializeString(string(inputBytesPlain), inputFormat, inputHeader, inputComment, inputLazyQuotes, inputLimit, inputType, verbose)
//  if err != nil {
//    fmt.Println(errors.Wrap(err, "error deserializing input using format "+inputFormat))
//    os.Exit(1)
//  }
//  ...
//  str, err := gss.SerializeString(object, "json", header, -1)
//  if err != nil {
//    return "", errors.Wrap(err, "error serializing object")
//  }
//
// Usage with options
//
// You can also call [Serialize|Deserialize][Bytes|String] using an options object.
//
//  options := gss.Options{
//    Header: inputHeader,
//    Comment: inputComment,
//    LazyQuotes: inputLazyQuotes,
//    Limit: 1,
//    Type: reflect.TypeOf(map[string]interface{}{}),
//  }
//
//  if inputFormat == "jsonl" {
//    options.Format = "json"
//  } else {
//    options.Format = inputFormat
//  }
//
//  for inputLine := range inputLines {
//    inputObject, err := options.DeserializeBytes(inputLine, verbose)
//    if err != nil {
//      errorsChannel <- errors.Wrap(err, "error deserializing input using format "+objectFormat)
//      continue
//    }
//    ...
//  }
//
// Formats
//
// GSS supports the following formats:
//
//  - bson
//  - csv
//  - tsv
//  - hcl
//  - hcl2
//  - json
//  - jsonl
//  - properties
//  - toml
//  - yaml
//
package gss

import (
	"fmt"
	"reflect"
)

import (
	"github.com/spatialcurrent/go-simple-serializer/pkg/bson"
	"github.com/spatialcurrent/go-simple-serializer/pkg/json"
	"github.com/spatialcurrent/go-simple-serializer/pkg/toml"
	"github.com/spatialcurrent/go-simple-serializer/pkg/yaml"
)

const (
	NoSkip    = 0  // used as SkipLines parameter to indicate no skipping when reading
	NoLimit   = -1 // used to indicate that there is no limit on reading or writing, depending on context.
	NoComment = "" // used to indicate that there is no comment prefix to consider.
)

// UnmarshalTypeFunc is a function for unmarshaling bytes into an object of a given type.
type UnmarshalTypeFunc func(b []byte, t reflect.Type) (interface{}, error)

// UnmarshalFunc is a function for unmarshaling bytes into an object of a given type.
type UnmarshalFunc func(b []byte) (interface{}, error)

// MarshalFunc is a function for marshalling an object into bytes.
type MarshalFunc func(object interface{}) ([]byte, error)

var (
	// NoHeader is used to indicate that no defined header is given.
	// Derive the header from the input data.
	NoHeader = []string{}
	// Formats is a list of all the formats supported by GSS
	Formats = []string{"bson", "csv", "tsv", "hcl", "hcl2", "json", "jsonl", "properties", "toml", "yaml"}
	// UnmarshalFuncs contains a map of functions for unmarshaling formatted bytes into objects.
	UnmarshalFuncs = map[string]UnmarshalFunc{
		"bson": bson.Unmarshal,
		"json": json.Unmarshal,
		"toml": toml.Unmarshal,
		"yaml": yaml.Unmarshal,
	}
	// UnmarshalTypeFuncs contains a map of functions for unmarshaling formatted bytes into objects.
	UnmarshalTypeFuncs = map[string]UnmarshalTypeFunc{
		"bson": bson.UnmarshalType,
		"json": json.UnmarshalType,
		"toml": toml.UnmarshalType,
		"yaml": yaml.UnmarshalType,
	}
	// MarshalTypeFuncs contains a map of functions for marshaling objects into formatted bytes.
	MarshalFuncs = map[string]MarshalFunc{
		"go": func(object interface{}) ([]byte, error) {
			return []byte(fmt.Sprint(object)), nil
		},
		"toml": toml.Marshal,
		"yaml": yaml.Marshal,
	}
)
