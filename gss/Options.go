// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gss

import (
	"reflect"
)

// Options is a struct that includes option values used during iterative processing.
type Options struct {
	Format          string       // one of gss.Formats
	Header          []string     // if formt as csv or tsv, the column names
	Comment         string       // the line comment prefix
	LazyQuotes      bool         // if format is csv or tsv, allow LazyQuotes.
	SkipLines       int          // if format is csv, tsv, or jsonl, the number of lines to skip before processing.
	Limit           int          // if format is a csv, tsv, or jsonl, then limit the number of items processed.
	Type            reflect.Type // the type of the output object
	Async           bool         // async processing
	Pretty          bool         // pretty output
	LineSeparator   string       // new line character, used by properties and jsonl
	Sorted          bool         // sort output
	ValueSerializer func(object interface{}) (string, error)
}

// Deserialize the input bytes using the values in the options object.
func (o Options) DeserializeBytes(content []byte, verbose bool) (interface{}, error) {
	return DeserializeBytes(&DeserializeInput{
		Bytes:      content,
		Format:     o.Format,
		Header:     o.Header,
		Comment:    o.Comment,
		LazyQuotes: o.LazyQuotes,
		SkipLines:  o.SkipLines,
		Limit:      o.Limit,
		Type:       o.Type,
		Async:      o.Async,
		Verbose:    verbose,
	})
}

// Deserialize the input string using the values in the options object.
func (o Options) DeserializeString(content string, verbose bool) (interface{}, error) {
	return DeserializeString(content, o.Format, o.Header, o.Comment, o.LazyQuotes, o.SkipLines, o.Limit, o.Type, o.Async, verbose)
}

func (o Options) SerializeString(object interface{}) (string, error) {
	return SerializeString(&SerializeInput{
		Object:          object,
		Format:          o.Format,
		Header:          o.Header,
		Limit:           o.Limit,
		Pretty:          o.Pretty,
		Sorted:          o.Sorted,
		LineSeparator:   o.LineSeparator,
		ValueSerializer: o.ValueSerializer,
	})
}
