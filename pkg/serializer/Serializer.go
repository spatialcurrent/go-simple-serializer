// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

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
	"github.com/spatialcurrent/go-simple-serializer/pkg/tags"
	"github.com/spatialcurrent/go-simple-serializer/pkg/toml"
	"github.com/spatialcurrent/go-simple-serializer/pkg/yaml"
)

const (
	FormatBSON       = "bson"
	FormatCSV        = "csv"
	FormatTSV        = "tsv"
	FormatJSON       = "json"
	FormatJSONL      = "jsonl"
	FormatProperties = "properties"
	FormatTags       = "tags"
	FormatHCL        = "hcl"
	FormatHCL2       = "hcl2"
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

// Serializer is a struct for serializing/deserializing objects.  This is the workhorse of the gss package.
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

// New returns a new serializer with the given format.
func New(format string) *Serializer {
	return &Serializer{
		format: format,
	}
}

// Format sets the format of the serializer.
func (s *Serializer) Format(format string) *Serializer {
	s.format = format
	return s
}

// Header sets the header of the serializer.
func (s *Serializer) Header(header []string) *Serializer {
	s.header = header
	return s
}

// Comment sets the comment of the serializer.
func (s *Serializer) Comment(comment string) *Serializer {
	s.comment = comment
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

// ValueSerializer sets the function for serializing values as strings for the csv, tsv, and properties formats.
func (s *Serializer) ValueSerializer(valueSerializer func(object interface{}) (string, error)) *Serializer {
	s.valueSerializer = valueSerializer
	return s
}

// Deserialize deserializes the input slice of bytes into an object and returns an error, if any.
func (s *Serializer) Deserialize(b []byte) (interface{}, error) {
	switch s.format {
	case FormatBSON, FormatJSON, "toml", "yaml":
		if s.objectType != nil {
			return UnmarshalTypeFuncs[s.format](b, s.objectType)
		}
		return UnmarshalFuncs[s.format](b)
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
		return properties.Read(&properties.ReadInput{
			Type:          s.objectType,
			Reader:        bytes.NewReader(b),
			LineSeparator: []byte(s.lineSeparator)[0],
			DropCR:        s.dropCR,
			Comment:       s.comment,
		})
	}
	return nil, &ErrUnknownFormat{Name: s.format}
}

// Serialize serializes an object into a slice of byte and returns and error, if any.
func (s *Serializer) Serialize(object interface{}) ([]byte, error) {

	valueSerializer := s.valueSerializer
	if valueSerializer == nil {
		valueSerializer = stringify.DefaultValueStringer("")
	}

	switch s.format {
	case FormatBSON:
		return bson.Marshal(stringify.StringifyMapKeys(object))
	case FormatCSV, FormatTSV:
		return make([]byte, 0), errors.New("not implemented")
	case FormatJSON:
		return json.Marshal(object, s.pretty)
	case FormatJSONL:
		return jsonl.Marshal(object, s.lineSeparator, s.pretty)
	case FormatProperties:
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
		return buf.Bytes(), nil
	case FormatTags:
		buf := new(bytes.Buffer)
		err := tags.Write(&tags.WriteInput{
			Writer:          buf,
			LineSeparator:   s.lineSeparator,
			Object:          object,
			ValueSerializer: valueSerializer,
			Sorted:          s.sorted,
		})
		if err != nil {
			return make([]byte, 0), errors.Wrap(err, "error writing tags")
		}
		return buf.Bytes(), nil
	case "go", "toml", "yaml":
		return MarshalFuncs[s.format](object)
	}
	return make([]byte, 0), &ErrUnknownFormat{Name: s.format}
}
