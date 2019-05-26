// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gss

import (
	"bytes"
	"reflect"
)

import (
	"github.com/pkg/errors"
)

import (
	"github.com/spatialcurrent/go-stringify/pkg/stringify"
)

import (
	"github.com/spatialcurrent/go-simple-serializer/pkg/bson"
	"github.com/spatialcurrent/go-simple-serializer/pkg/json"
	"github.com/spatialcurrent/go-simple-serializer/pkg/jsonl"
	"github.com/spatialcurrent/go-simple-serializer/pkg/properties"
)

// Serializer is a struct for serializing
type Serializer struct {
	format            string       // one of gss.Formats
	header            []string     // if formt as csv or tsv, the column names
	comment           string       // the line comment prefix
	LazyQuotes        bool         // if format is csv or tsv, allow LazyQuotes.
	SkipLines         int          // if format is csv, tsv, or jsonl, the number of lines to skip before processing.
	limit             int          // if format is a csv, tsv, or jsonl, then limit the number of items processed.
	objectType        reflect.Type // the type of the output object
	Async             bool         // async processing
	pretty            bool         // pretty output
	lineSeparator     string       // new line character, used by properties and jsonl
	keyValueSeparator string
	sorted            bool // sort output
	valueSerializer   func(object interface{}) (string, error)
	escapePrefix      string
	escapeSpace       bool
	escapeNewLine     bool
	escapeColon       bool
	escapeEqual       bool
	unescapeSpace     bool
	unescapeNewLine   bool
	unescapeColon     bool
	unescapeEqual     bool
	trim              bool
	dropCR            bool
}

func NewSerializer(format string) *Serializer {
	return &Serializer{
		format: format,
	}
}
func (s *Serializer) Format(format string) *Serializer {
	s.format = format
	return s
}

func (s *Serializer) Header(header []string) *Serializer {
	s.header = header
	return s
}

func (s *Serializer) Comment(comment string) *Serializer {
	s.comment = comment
	return s
}

func (s *Serializer) Limit(limit int) *Serializer {
	s.limit = limit
	return s
}

func (s *Serializer) LineSeparator(lineSeparator string) *Serializer {
	s.lineSeparator = lineSeparator
	return s
}

func (s *Serializer) KeyValueSeparator(keyValueSeparator string) *Serializer {
	s.keyValueSeparator = keyValueSeparator
	return s
}

func (s *Serializer) Pretty(pretty bool) *Serializer {
	s.pretty = pretty
	return s
}

func (s *Serializer) Sorted(sorted bool) *Serializer {
	s.sorted = sorted
	return s
}

func (s *Serializer) Type(t reflect.Type) *Serializer {
	s.objectType = t
	return s
}

func (s *Serializer) EscapePrefix(escapePrefix string) *Serializer {
	s.escapePrefix = escapePrefix
	return s
}

func (s *Serializer) EscapeSpace(escapeSpace bool) *Serializer {
	s.escapeSpace = escapeSpace
	return s
}

func (s *Serializer) EscapeNewLine(escapeNewLine bool) *Serializer {
	s.escapeNewLine = escapeNewLine
	return s
}

func (s *Serializer) EscapeEqual(escapeEqual bool) *Serializer {
	s.escapeEqual = escapeEqual
	return s
}

func (s *Serializer) EscapeColon(escapeColon bool) *Serializer {
	s.escapeColon = escapeColon
	return s
}

func (s *Serializer) UnescapeSpace(unescapeSpace bool) *Serializer {
	s.unescapeSpace = unescapeSpace
	return s
}

func (s *Serializer) UnescapeNewLine(unescapeNewLine bool) *Serializer {
	s.unescapeNewLine = unescapeNewLine
	return s
}

func (s *Serializer) UnescapeEqual(unescapeEqual bool) *Serializer {
	s.unescapeEqual = unescapeEqual
	return s
}

func (s *Serializer) UnescapeColon(unescapeColon bool) *Serializer {
	s.unescapeColon = unescapeColon
	return s
}

func (s *Serializer) Trim(trim bool) *Serializer {
	s.trim = trim
	return s
}

func (s *Serializer) DropCR(dropCR bool) *Serializer {
	s.dropCR = dropCR
	return s
}

func (s *Serializer) ValueSerializer(valueSerializer func(object interface{}) (string, error)) *Serializer {
	s.valueSerializer = valueSerializer
	return s
}

func (s *Serializer) Deserialize(b []byte) (interface{}, error) {
	switch s.format {
	case "bson", "json", "toml", "yaml":
		if s.objectType != nil {
			return UnmarshalTypeFuncs[s.format](b, s.objectType)
		}
		return UnmarshalFuncs[s.format](b)
	case "properties":
		return properties.Read(&properties.ReadInput{
			Type:            s.objectType,
			Reader:          bytes.NewReader(b),
			LineSeparator:   []byte(s.lineSeparator)[0],
			DropCR:          s.dropCR,
			Comment:         s.comment,
			Trim:            s.trim,
			UnescapeSpace:   s.unescapeSpace,
			UnescapeEqual:   s.unescapeEqual,
			UnescapeColon:   s.unescapeColon,
			UnescapeNewLine: s.unescapeNewLine,
		})
	}
	return nil, &ErrUnknownFormat{Name: s.format}
}

func (s *Serializer) Serialize(object interface{}) ([]byte, error) {
	switch s.format {
	case "bson":
		return bson.Marshal(stringify.StringifyMapKeys(object))
	case "json":
		return json.Marshal(object, s.pretty)
	case "jsonl":
		return jsonl.Marshal(object, s.lineSeparator, s.pretty)
	case "properties":

		valueSerializer := s.valueSerializer
		if valueSerializer == nil {
			valueSerializer = stringify.DefaultValueStringer("")
		}

		buf := new(bytes.Buffer)
		err := properties.Write(&properties.WriteInput{
			Writer:            buf,
			LineSeparator:     s.lineSeparator,
			KeyValueSeparator: s.keyValueSeparator,
			Object:            object,
			ValueSerializer:   valueSerializer,
			Sorted:            s.sorted,
			EscapePrefix:      s.escapePrefix,
			EscapeSpace:       s.escapeSpace,
			EscapeColon:       s.escapeColon,
			EscapeNewLine:     s.escapeNewLine,
			EscapeEqual:       s.escapeEqual,
		})
		if err != nil {
			return make([]byte, 0), errors.Wrap(err, "error writing properties")
		}
		return buf.Bytes(), err
	case "go", "toml", "yaml":
		return MarshalFuncs[s.format](object)
	}
	return make([]byte, 0), &ErrUnknownFormat{Name: s.format}
}
