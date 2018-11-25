// =================================================================
//
// Copyright (C) 2018 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gss

import (
	"reflect"
)

// Options is a struct that includes option values used during iterative processing.
type Options struct {
	Format     string       // one of gss.Formats
	Header     []string     // if formt as csv or tsv, the column names
	Comment    string       // the line comment prefix
	LazyQuotes bool         // if format is csv or tsv, allow LazyQuotes.
	SkipLines  int          // if format is csv, tsv, or jsonl, the number of lines to skip before processing.
	Limit      int          // if format is a csv, tsv, or jsonl, then limit the number of items processed.
	Type       reflect.Type // the type of the output object
}

// Deserialize the input bytes using the values in the options object.
func (o Options) DeserializeBytes(content []byte, verbose bool) (interface{}, error) {
	return DeserializeBytes(content, o.Format, o.Header, o.Comment, o.LazyQuotes, o.SkipLines, o.Limit, o.Type, verbose)
}

// Deserialize the input string using the values in the options object.
func (o Options) DeserializeString(content string, verbose bool) (interface{}, error) {
	return DeserializeString(content, o.Format, o.Header, o.Comment, o.LazyQuotes, o.SkipLines, o.Limit, o.Type, verbose)
}

func (o Options) SerializeString(object interface{}) (string, error) {
	return SerializeString(object, o.Format, o.Header, o.Limit)
}
