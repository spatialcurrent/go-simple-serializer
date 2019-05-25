// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

// Package serializer provides a middle layer between the gss package and the lower-level packages.
// This package provides a simple api in the builder pattern for serializing/deserializing objects.
package serializer

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
)

import (
	"github.com/spatialcurrent/go-stringify/pkg/stringify"
)

import (
	"github.com/spatialcurrent/go-simple-serializer/pkg/bson"
	"github.com/spatialcurrent/go-simple-serializer/pkg/json"
	"github.com/spatialcurrent/go-simple-serializer/pkg/jsonl"
	"github.com/spatialcurrent/go-simple-serializer/pkg/properties"
	"github.com/spatialcurrent/go-simple-serializer/pkg/sv"
	"github.com/spatialcurrent/go-simple-serializer/pkg/tags"
	"github.com/spatialcurrent/go-simple-serializer/pkg/toml"
	"github.com/spatialcurrent/go-simple-serializer/pkg/yaml"
)

const (
	FormatBSON       = "bson"       // Binary JSON
	FormatCSV        = "csv"        // Comma-Separated Values
	FormatGo         = "go"         // Native Golang print format
	FormatJSON       = "json"       // JSON
	FormatJSONL      = "jsonl"      // JSON Lines
	FormatProperties = "properties" // Properties
	FormatTags       = "tags"       // Tags (a=b c=d ...)
	FormatTOML       = "toml"       // TOML
	FormatTSV        = "tsv"        // Tab-Separated Values
	FormatHCL        = "hcl"        // HashiCorp Configuration Language
	FormatHCL2       = "hcl2"       // HashiCorp Configuration Language Version 2.x
	FormatYAML       = "yaml"       // YAML

	NoLimit = -1
)

var (
	Formats = []string{
		FormatBSON,
		FormatCSV,
		FormatGo,
		FormatJSON,
		FormatJSONL,
		FormatProperties,
		FormatTags,
		FormatTOML,
		FormatTSV,
		FormatHCL,
		FormatHCL2,
		FormatYAML,
	}
	ErrMissingLineSeparator = errors.New("missing line separator")
)

// UnmarshalTypeFunc is a function for unmarshaling bytes into an object of a given type.
type UnmarshalTypeFunc func(b []byte, t reflect.Type) (interface{}, error)

// UnmarshalFunc is a function for unmarshaling bytes into an object of a given type.
type UnmarshalFunc func(b []byte) (interface{}, error)

// MarshalFunc is a function for marshalling an object into bytes.
type MarshalFunc func(object interface{}) ([]byte, error)

var (
	// UnmarshalFuncs contains a map of functions for unmarshaling formatted bytes into objects.
	UnmarshalFuncs = map[string]UnmarshalFunc{
		FormatBSON: bson.Unmarshal,
		FormatJSON: json.Unmarshal,
		FormatTags: tags.Unmarshal,
		FormatTOML: toml.Unmarshal,
		FormatYAML: yaml.Unmarshal,
	}
	// UnmarshalTypeFuncs contains a map of functions for unmarshaling formatted bytes into objects.
	UnmarshalTypeFuncs = map[string]UnmarshalTypeFunc{
		FormatBSON: bson.UnmarshalType,
		FormatJSON: json.UnmarshalType,
		FormatTags: tags.UnmarshalType,
		FormatTOML: toml.UnmarshalType,
		FormatYAML: yaml.UnmarshalType,
	}
	// MarshalTypeFuncs contains a map of functions for marshaling objects into formatted bytes.
	MarshalFuncs = map[string]MarshalFunc{
		FormatTOML: toml.Marshal,
		FormatYAML: yaml.Marshal,
	}
)

// Serializer is a struct for serializing/deserializing objects.  This is the workhorse of the gss package.
type Serializer struct {
	format            string        // one of gss.Formats
	header            []interface{} // if formt as csv or tsv, the column names
	comment           string        // the line comment prefix
	lazyQuotes        bool          // if format is csv or tsv, allow LazyQuotes.
	skipLines         int           // if format is csv, tsv, or jsonl, the number of lines to skip before processing.
	skipBlanks        bool          // Skip blank lines.  If false, Next() returns a blank line as (nil, nil).  If true, Next() simply skips forward until it finds a non-blank line.
	skipComments      bool          // Skip commented lines.  If false, Next() returns a commented line as (nil, nil).  If true, Next() simply skips forward until it finds a non-commented line.
	limit             int           // if format is a csv, tsv, or jsonl, then limit the number of items processed.
	objectType        reflect.Type  // the type of the output object
	pretty            bool          // pretty output
	lineSeparator     string        // new line character, used by properties and jsonl
	keyValueSeparator string
	sorted            bool // sort output
	reversed          bool // if sorted, sort in reverse alphabetical order
	keySerializer     stringify.Stringer
	valueSerializer   stringify.Stringer
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
	expandHeader      bool // dynamically expand header, requires caching output in memory
}

// New returns a new serializer with the given format.
func New(format string) *Serializer {
	return &Serializer{
		format: format,
	}
}

// NewWithOptions returns a new serializer configured with the given options.
func NewWithOptions(format string, options ...map[string]interface{}) (*Serializer, error) {
	s := &Serializer{
		format: format,
	}
	for _, opt := range options {
		for key, value := range opt {
			switch key {
			case "escapePrefix":
				s = s.EscapePrefix(fmt.Sprint(value))
			case "comment":
				s = s.Comment(fmt.Sprint(value))
			case "lineSeparator":
				s = s.LineSeparator(fmt.Sprint(value))
			case "keyValueSeparator":
				s = s.KeyValueSeparator(fmt.Sprint(value))
			case "limit":
				switch v := value.(type) {
				case int:
					s = s.Limit(v)
				case float64:
					s = s.Limit(int(v))
				}
			case "pretty":
				switch v := value.(type) {
				case bool:
					s = s.Pretty(v)
				case int:
					s = s.Pretty(v > 0)
				case float64:
					s = s.Pretty(v > 0.0)
				}
			case "trim":
				switch v := value.(type) {
				case bool:
					s = s.Trim(v)
				case int:
					s = s.Trim(v > 0)
				case float64:
					s = s.Trim(v > 0.0)
				}
			case "sorted":
				switch v := value.(type) {
				case bool:
					s = s.Sorted(v)
				case int:
					s = s.Sorted(v > 0)
				case float64:
					s = s.Sorted(v > 0.0)
				}
			case "reversed":
				switch v := value.(type) {
				case bool:
					s = s.Reversed(v)
				case int:
					s = s.Reversed(v > 0)
				case float64:
					s = s.Reversed(v > 0.0)
				}
			case "expandHeader":
				switch v := value.(type) {
				case bool:
					s = s.ExpandHeader(v)
				case int:
					s = s.ExpandHeader(v > 0)
				case float64:
					s = s.ExpandHeader(v > 0.0)
				}
			case "header":
				s = s.Header(toInterfaceSlice(value))
			default:
				return s, &ErrUnknownOption{Name: key}
			}
		}
	}

	return s, nil
}

// Format sets the format of the serializer.
func (s *Serializer) Format(format string) *Serializer {
	s.format = format
	return s
}

// Header sets the header of the serializer.
func (s *Serializer) Header(header []interface{}) *Serializer {
	s.header = header
	return s
}

// ExpandHeader enables/disables dynamically expanding the header.
// Dynamically expanding the header requires buffering the output in memory.
func (s *Serializer) ExpandHeader(expandHeader bool) *Serializer {
	s.expandHeader = expandHeader
	return s
}

// Comment sets the comment of the serializer.
func (s *Serializer) Comment(comment string) *Serializer {
	s.comment = comment
	return s
}

// SkipLines sets the number of lines to skip from the beginning of the input.
func (s *Serializer) SkipLines(skipLines int) *Serializer {
	s.skipLines = skipLines
	return s
}

// Comment sets the comment of the serializer.
func (s *Serializer) SkipBlanks(skipBlanks bool) *Serializer {
	s.skipBlanks = skipBlanks
	return s
}

// Comment sets the comment of the serializer.
func (s *Serializer) SkipComments(skipComments bool) *Serializer {
	s.skipComments = skipComments
	return s
}

// Limit seets the limit of the serializer.
func (s *Serializer) Limit(limit int) *Serializer {
	s.limit = limit
	return s
}

// LineSeparator sets the line separator of the serializer.
func (s *Serializer) LineSeparator(lineSeparator string) *Serializer {
	s.lineSeparator = lineSeparator
	return s
}

// KeyValueSeparator sets the key-value separator of the serializer.
func (s *Serializer) KeyValueSeparator(keyValueSeparator string) *Serializer {
	s.keyValueSeparator = keyValueSeparator
	return s
}

// Pretty enables/disables pretty output.
func (s *Serializer) Pretty(pretty bool) *Serializer {
	s.pretty = pretty
	return s
}

// Sorted enables/disables sorted output.
func (s *Serializer) Sorted(sorted bool) *Serializer {
	s.sorted = sorted
	return s
}

// Reversed enables/disables reversed output.
func (s *Serializer) Reversed(reversed bool) *Serializer {
	s.reversed = reversed
	return s
}

// Type sets the optional type for deserialization.
// If no type is given, then the type is inferred from the source.
func (s *Serializer) Type(t reflect.Type) *Serializer {
	s.objectType = t
	return s
}

// EscapePrefix sets the prefix for escaping text.  Used with the properties format.
// If the escape prefix is not set, then the serializer doesn't escape/unescape any text.
func (s *Serializer) EscapePrefix(escapePrefix string) *Serializer {
	s.escapePrefix = escapePrefix
	return s
}

// EscapeSpace enables/disables escaping the whitespace character.
func (s *Serializer) EscapeSpace(escapeSpace bool) *Serializer {
	s.escapeSpace = escapeSpace
	return s
}

// EscapeSpace enables/disables escaping the new line character.
func (s *Serializer) EscapeNewLine(escapeNewLine bool) *Serializer {
	s.escapeNewLine = escapeNewLine
	return s
}

// EscapeSpace enables/disables escaping the equal character.
func (s *Serializer) EscapeEqual(escapeEqual bool) *Serializer {
	s.escapeEqual = escapeEqual
	return s
}

// EscapeSpace enables/disables escaping the colon character.
func (s *Serializer) EscapeColon(escapeColon bool) *Serializer {
	s.escapeColon = escapeColon
	return s
}

// UnescapeSpace enables/disables unescaping the whitespace character.
func (s *Serializer) UnescapeSpace(unescapeSpace bool) *Serializer {
	s.unescapeSpace = unescapeSpace
	return s
}

// UnescapeNewLine enables/disables unescaping the new line character.
func (s *Serializer) UnescapeNewLine(unescapeNewLine bool) *Serializer {
	s.unescapeNewLine = unescapeNewLine
	return s
}

// UnescapeEqual enables/disables unescpaing the equal character.
func (s *Serializer) UnescapeEqual(unescapeEqual bool) *Serializer {
	s.unescapeEqual = unescapeEqual
	return s
}

// UnescapeColon enables/disables unescaping the colon character.
func (s *Serializer) UnescapeColon(unescapeColon bool) *Serializer {
	s.unescapeColon = unescapeColon
	return s
}

// Trim enables/disables trimming whitespace from input lines.
func (s *Serializer) Trim(trim bool) *Serializer {
	s.trim = trim
	return s
}

// DropCR enables/disables dropping the carriage return character if it terminates a line.
func (s *Serializer) DropCR(dropCR bool) *Serializer {
	s.dropCR = dropCR
	return s
}

// KeySerializer sets the function for serializing keys as strings for the csv, tsv, and properties formats.
func (s *Serializer) KeySerializer(keySerializer stringify.Stringer) *Serializer {
	s.keySerializer = keySerializer
	return s
}

// ValueSerializer sets the function for serializing values as strings for the csv, tsv, and properties formats.
func (s *Serializer) ValueSerializer(valueSerializer stringify.Stringer) *Serializer {
	s.valueSerializer = valueSerializer
	return s
}

// LazyQuotes enables/disables lazy quotes when reading from an input formatted as separated values, e.g., CSV or TSV.
func (s *Serializer) LazyQuotes(lazyQuotes bool) *Serializer {
	s.lazyQuotes = lazyQuotes
	return s
}

// Deserialize deserializes the input slice of bytes into an object and returns an error, if any.
// Formats jsonl and tags return slices.  If the type is not set, then returns a slice of type []interface{}.
func (s *Serializer) Deserialize(b []byte) (interface{}, error) {
	switch s.format {
	case FormatBSON, FormatJSON, FormatTOML, FormatYAML:
		if s.objectType != nil {
			return UnmarshalTypeFuncs[s.format](b, s.objectType)
		}
		return UnmarshalFuncs[s.format](b)
	case FormatCSV, FormatTSV:
		separator, errSeparator := sv.FormatToSeparator(s.format)
		if errSeparator != nil {
			return make([]byte, 0), errSeparator
		}
		return sv.Read(&sv.ReadInput{
			Type:       s.objectType,
			Reader:     bytes.NewReader(b),
			Separator:  separator,
			Header:     s.header,
			SkipLines:  0,
			Comment:    s.comment,
			LazyQuotes: s.lazyQuotes,
			Limit:      s.limit,
		})
	case FormatJSONL, FormatProperties, FormatTags:
		if len(s.lineSeparator) == 0 {
			return nil, ErrMissingLineSeparator
		}
		switch s.format {
		case FormatJSONL:
			return jsonl.Read(&jsonl.ReadInput{
				Type:          s.objectType,
				Reader:        bytes.NewReader(b),
				LineSeparator: []byte(s.lineSeparator)[0],
				DropCR:        s.dropCR,
				Comment:       s.comment,
				SkipLines:     s.skipLines,
				SkipBlanks:    s.skipBlanks,
				SkipComments:  s.skipComments,
				Limit:         s.limit,
				Trim:          s.trim,
			})
		case FormatProperties:
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
		case FormatTags:
			return tags.Read(&tags.ReadInput{
				Type:          s.objectType,
				Reader:        bytes.NewReader(b),
				LineSeparator: []byte(s.lineSeparator)[0],
				DropCR:        s.dropCR,
				Comment:       s.comment,
				SkipLines:     s.skipLines,
				SkipBlanks:    s.skipBlanks,
				SkipComments:  s.skipComments,
				Limit:         s.limit,
			})
		}
	case FormatHCL:
		ptr := reflect.New(s.objectType)
		ptr.Elem().Set(reflect.MakeMap(s.objectType))
		obj, err := hcl.Parse(string(b))
		if err != nil {
			return nil, errors.Wrap(err, "Error parsing hcl")
		}
		if err := hcl.DecodeObject(ptr.Interface(), obj); err != nil {
			return nil, errors.Wrap(err, "Error decoding hcl")
		}
		return ptr.Elem().Interface(), nil
	case FormatHCL2:
		file, diags := hclsyntax.ParseConfig(b, "<stdin>", hcl2.Pos{Byte: 0, Line: 1, Column: 1})
		if diags.HasErrors() {
			return nil, errors.Wrap(errors.New(diags.Error()), "Error parsing hcl2")
		}
		return &file.Body, nil
	}
	return nil, &ErrUnknownFormat{Name: s.format}
}

// Serialize serializes an object into a slice of byte and returns and error, if any.
func (s *Serializer) Serialize(object interface{}) ([]byte, error) {

	keySerializer := s.keySerializer
	if keySerializer == nil {
		keySerializer = stringify.NewStringer("", false, false, false)
	}

	valueSerializer := s.valueSerializer
	if valueSerializer == nil {
		valueSerializer = stringify.NewStringer("", false, false, false)
	}

	switch s.format {
	case FormatBSON:
		o, err := stringify.StringifyMapKeys(object, keySerializer)
		if err != nil {
			return make([]byte, 0), errors.Wrap(err, "error stringifying map keys")
		}
		return bson.Marshal(o)
	case FormatCSV, FormatTSV:
		separator, errSeparator := sv.FormatToSeparator(s.format)
		if errSeparator != nil {
			return make([]byte, 0), errSeparator
		}
		buf := new(bytes.Buffer)
		errWrite := sv.Write(&sv.WriteInput{
			Writer:          buf,
			Separator:       separator,
			Object:          object,
			KeySerializer:   keySerializer,
			ValueSerializer: valueSerializer,
			Sorted:          s.sorted,
			Reversed:        s.reversed,
			Header:          s.header,
			ExpandHeader:    s.expandHeader,
			Limit:           s.limit,
		})
		if errWrite != nil {
			return make([]byte, 0), errors.Wrap(errWrite, "error writing separated values")
		}
		return buf.Bytes(), nil
	case FormatGo:
		// TODO:
		// Pretty output disabled until https://github.com/kr/pretty/issues/45 is fixed
		// 	krpretty "github.com/kr/pretty"
		//if s.pretty {
		//	return []byte(krpretty.Sprint(object)), nil
		//}
		return []byte(fmt.Sprintf("%#v", object)), nil
	case FormatJSON:
		o, err := stringify.StringifyMapKeys(object, keySerializer)
		if err != nil {
			return make([]byte, 0), errors.Wrap(err, "error stringifying map keys")
		}
		return json.Marshal(o, s.pretty)
	case FormatJSONL:
		return jsonl.Marshal(object, s.lineSeparator, keySerializer, s.pretty, s.limit)
	case FormatProperties:
		buf := new(bytes.Buffer)
		err := properties.Write(&properties.WriteInput{
			Writer:            buf,
			LineSeparator:     s.lineSeparator,
			KeyValueSeparator: s.keyValueSeparator,
			Object:            object,
			KeySerializer:     keySerializer,
			ValueSerializer:   valueSerializer,
			Sorted:            s.sorted,
			Reversed:          s.reversed,
			EscapePrefix:      s.escapePrefix,
			EscapeSpace:       s.escapeSpace,
			EscapeColon:       s.escapeColon,
			EscapeNewLine:     s.escapeNewLine,
			EscapeEqual:       s.escapeEqual,
		})
		if err != nil {
			return make([]byte, 0), errors.Wrap(err, "error writing properties")
		}
		return buf.Bytes(), nil
	case FormatTags:
		buf := new(bytes.Buffer)
		err := tags.Write(&tags.WriteInput{
			Writer:          buf,
			LineSeparator:   s.lineSeparator,
			Object:          object,
			KeySerializer:   keySerializer,
			ValueSerializer: valueSerializer,
			Sorted:          s.sorted,
			Reversed:        s.reversed,
			Limit:           s.limit,
		})
		if err != nil {
			return make([]byte, 0), errors.Wrap(err, "error writing tags")
		}
		return buf.Bytes(), nil
	case FormatTOML, FormatYAML:
		o, err := stringify.StringifyMapKeys(object, keySerializer)
		if err != nil {
			return make([]byte, 0), errors.Wrap(err, "error stringifying map keys")
		}
		return MarshalFuncs[s.format](o)
	}
	return make([]byte, 0), &ErrUnknownFormat{Name: s.format}
}
