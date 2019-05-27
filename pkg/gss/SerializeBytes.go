// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gss

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

import (
	"github.com/pkg/errors"
)

import (
	"github.com/spatialcurrent/go-stringify/pkg/stringify"
)

import (
	"github.com/spatialcurrent/go-simple-serializer/pkg/inspector"
	"github.com/spatialcurrent/go-simple-serializer/pkg/serializer"
	sv "github.com/spatialcurrent/go-simple-serializer/pkg/sv"
)

/*
func serializeRow(header []interface{}, knownKeys map[string]struct{}, obj interface{}) ([]string, map[string]struct{}, []string, error) {
	m := reflect.ValueOf(obj)
	if m.Kind() != reflect.Map {
		return header, knownKeys, make([]string, 0), &ErrInvalidKind{Value: m.Kind(), Valid: []reflect.Kind{reflect.Map}}
	}

	newHeader, newKnownKeys := sv.ExpandHeader(header, knownKeys, m, false)
  row := sv.ToRow(obj, newHeader, valueSerializer)

	row := make([]string, len(newHeader))
	for j, key := range newHeader {
		if v := m.MapIndex(reflect.ValueOf(key)); v.IsValid() && !v.IsNil() {
			row[j] = fmt.Sprint(v.Interface())
		} else {
			row[j] = ""
		}
	}

	return newHeader, newKnownKeys, row, nil
}
*/

// SerializeBytes serializes an object to its representation given by format.
func SerializeBytes(input *SerializeInput) ([]byte, error) {

	object := input.Object
	format := input.Format
	limit := input.Limit
	valueSerializer := input.ValueSerializer
	if valueSerializer == nil {
		valueSerializer = stringify.DefaultValueStringer("")
	}

	switch format {
	case "bson", "json", "jsonl", "properties", "go", "toml", "yaml":
		s := serializer.New(format)
		if format == "json" || format == "jsonl" {
			s = s.Pretty(input.Pretty)
		}
		if format == "jsonl" || format == "properties" {
			s = s.LineSeparator(input.LineSeparator)
		}
		if format == "properties" {
			s = s.
				Sorted(input.Sorted).
				KeyValueSeparator(input.KeyValueSeparator).
				LineSeparator(input.LineSeparator).
				ValueSerializer(valueSerializer).
				EscapePrefix(input.EscapePrefix).
				EscapeSpace(input.EscapeSpace).
				EscapeColon(input.EscapeColon).
				EscapeNewLine(input.EscapeNewLine).
				EscapeEqual(input.EscapeEqual)
		}
		return s.Serialize(object)
	case "hcl", "hcl2":
		return make([]byte, 0), fmt.Errorf("cannot serialize to format %q", format)
	}

	if format == "csv" || format == "tsv" {

		separator, err := sv.FormatToSeparator(format)
		if err != nil {
			return make([]byte, 0), err
		}

		header := input.Header
		wildcard := false
		knownKeys := map[interface{}]struct{}{}
		if len(header) > 0 {
			for _, k := range header {
				if k == "*" {
					wildcard = true
				} else {
					knownKeys[k] = struct{}{}
				}
			}
			if input.Sorted {
				sort.Slice(header, func(i, j int) bool {
					return fmt.Sprint(header[i]) < fmt.Sprint(header[j])
				})
			}
		}

		rows := make([][]string, 0)

		s := reflect.ValueOf(object)
		switch s.Type().Kind() {
		case reflect.Map, reflect.Struct:
			if len(header) == 0 {
				header, knownKeys = sv.CreateHeaderAndKnownKeys(s, input.Sorted)
			} else if wildcard {
				newHeader, newKnownKeys := sv.ExpandHeader(header, knownKeys, s, input.Sorted)
				header = newHeader
				knownKeys = newKnownKeys
			}
			row, err := sv.ToRowFromValue(s, header, valueSerializer)
			if err != nil {
				return make([]byte, 0), errors.Wrap(err, "error serializing object to row")
			}
			rows = append(rows, row)
		case reflect.Array, reflect.Slice:
			if s.Len() > 0 {
				if len(header) == 0 {
					header, knownKeys = sv.CreateHeaderAndKnownKeys(s.Index(0), input.Sorted)
				}
				for i := 0; i < s.Len() && (limit < 0 || i <= limit); i++ {
					if wildcard {
						header, knownKeys = sv.ExpandHeader(header, knownKeys, s.Index(i), input.Sorted)
					}
					row, err := sv.ToRowFromValue(s.Index(i), header, valueSerializer)
					if err != nil {
						return make([]byte, 0), errors.Wrap(err, "error serializing object to row")
					}
					rows = append(rows, row)
				}
			}
			// If there are no records then just return an empty string
			return []byte(""), nil
		}
		// Write header and rows
		headerAsStrings := make([]string, 0, len(header))
		for _, x := range header {
			headerAsStrings = append(headerAsStrings, fmt.Sprint(x))
		}
		buf := new(bytes.Buffer)
		err = sv.Write(&sv.WriteInput{
			Writer:    buf,
			Separator: separator,
			Header:    headerAsStrings,
			Rows:      rows,
		})
		return buf.Bytes(), err
	} else if format == "text" {
		t := reflect.TypeOf(object)
		if t.Kind() == reflect.Map {
			if k := t.Key().Kind(); k != reflect.String {
				return nil, errors.Wrap(&ErrInvalidKind{Value: k, Valid: []reflect.Kind{reflect.String}}, "can only serialize a map with string keys")
			}
			m := reflect.ValueOf(object)
			keys := inspector.GetKeys(object, input.Sorted)
			output := ""
			for i, key := range keys {
				value, err := valueSerializer(m.MapIndex(reflect.ValueOf(key)).Interface())
				if err != nil {
					return nil, errors.Wrap(err, "error serializing value")
				}
				value = strings.Replace(value, "\"", "\\\"", -1)
				if strings.Contains(value, " ") {
					value = "\"" + value + "\""
				}
				output += fmt.Sprint(key) + "=" + value
				if i < m.Len()-1 {
					output += " "
				}
			}
			return []byte(output), nil
		}
		switch obj := object.(type) {
		case string:
			return []byte(obj), nil
		case int:
			return []byte(strconv.Itoa(obj)), nil
		case float64:
			return []byte(strconv.FormatFloat(obj, 'f', -1, 64)), nil
		}
		return make([]byte, 0), errors.Wrap(&ErrInvalidKind{Value: reflect.TypeOf(object).Kind(), Valid: []reflect.Kind{reflect.Map, reflect.String, reflect.Int, reflect.Float64}}, "object is not valid")
	}
	return make([]byte, 0), errors.Wrap(&ErrUnknownFormat{Name: format}, "could not serialize object")
}
