// =================================================================
//
// Copyright (C) 2018 Spatial Current, Inc. - All Rights Reserved
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

var NoSkip = 0
var NoLimit = -1
var NoHeader = []string{}
var NoComment = ""
