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

	"github.com/pkg/errors"
	"github.com/spatialcurrent/go-simple-serializer/pkg/serializer"
)

// SerializeBytesInput provides the input for the SerializeString and SerializeBytes function.
type SerializeBytesInput struct {
	Object            interface{}
	Format            string
	FormatSpecifier   string
	Fit               bool
	Header            []interface{}
	Limit             int
	Pretty            bool
	Sorted            bool
	Reversed          bool
	LineSeparator     string
	KeyValueSeparator string
	KeySerializer     func(object interface{}) (string, error)
	ValueSerializer   func(object interface{}) (string, error)
	EscapePrefix      string
	EscapeSpace       bool
	EscapeNewLine     bool
	EscapeEqual       bool
	EscapeColon       bool
	ExpandHeader      bool
}

// SerializeBytes serializes an object to its representation given by format.
func SerializeBytes(input *SerializeBytesInput) ([]byte, error) {

	f := input.Format

	switch f {
	case "bson", "csv", "fmt", "go", "gob", "json", "jsonl", "properties", "tags", "toml", "tsv", "yaml":
		s := serializer.New(f)
		if f == serializer.FormatFmt {
			s = s.FormatSpecifier(input.FormatSpecifier)
		}
		if f == serializer.FormatGo || f == serializer.FormatJSON || f == serializer.FormatJSONL {
			s = s.Pretty(input.Pretty)
		}
		if f == serializer.FormatJSONL || f == serializer.FormatProperties || f == serializer.FormatTags {
			s = s.LineSeparator(input.LineSeparator)
		}
		if f == serializer.FormatProperties || f == serializer.FormatTags {
			s = s.KeyValueSeparator(input.KeyValueSeparator)
		}
		if f == serializer.FormatCSV || f == serializer.FormatProperties || f == serializer.FormatTags || f == serializer.FormatTSV {
			// Sort the order of the keys/properties
			// Does not sort the order of the records (if serializing multiples objects as tags)
			// If sorted and reversed, then sort in reverse alphabetical order.
			s = s.
				KeySerializer(input.KeySerializer).
				ValueSerializer(input.ValueSerializer).
				Sorted(input.Sorted).
				Reversed(input.Reversed)
		}
		if f == serializer.FormatCSV || f == serializer.FormatTSV || f == serializer.FormatTags {
			s = s.Header(input.Header).ExpandHeader(input.ExpandHeader)
		}
		if f == "properties" {
			s = s.
				EscapePrefix(input.EscapePrefix).
				EscapeSpace(input.EscapeSpace).
				EscapeColon(input.EscapeColon).
				EscapeNewLine(input.EscapeNewLine).
				EscapeEqual(input.EscapeEqual)
		}
		if f == "csv" || f == "jsonl" || f == "tags" || f == "tsv" {
			s = s.Limit(input.Limit)
		}
		if f == "gob" {
			s = s.
				Fit(input.Fit).
				Type(reflect.TypeOf(make([]interface{}, 0)))
		}
		return s.Serialize(input.Object)
	case "hcl", "hcl2":
		return make([]byte, 0), fmt.Errorf("cannot serialize to format %q", f)
	}
	return make([]byte, 0), errors.Wrap(&ErrUnknownFormat{Name: f}, "could not serialize object")
}
