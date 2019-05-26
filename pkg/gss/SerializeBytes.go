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
	"gopkg.in/mgo.v2/bson"
)

import (
	"github.com/spatialcurrent/go-stringify/pkg/stringify"
)

import (
	"github.com/spatialcurrent/go-simple-serializer/pkg/inspector"
	json "github.com/spatialcurrent/go-simple-serializer/pkg/json"
	jsonl "github.com/spatialcurrent/go-simple-serializer/pkg/jsonl"
	properties "github.com/spatialcurrent/go-simple-serializer/pkg/properties"
	sv "github.com/spatialcurrent/go-simple-serializer/pkg/sv"
	toml "github.com/spatialcurrent/go-simple-serializer/pkg/toml"
	yaml "github.com/spatialcurrent/go-simple-serializer/pkg/yaml"
)

func unknownKeys(obj reflect.Value, knownKeys map[string]struct{}) []string {
	unknownKeys := make([]string, 0)
	for _, k := range obj.MapKeys() {
		str := fmt.Sprint(k.Interface())
		if _, exists := knownKeys[str]; !exists {
			unknownKeys = append(unknownKeys, str)
		}
	}
	return unknownKeys
}

func serializeRow(header []string, knownKeys map[string]struct{}, obj interface{}) ([]string, map[string]struct{}, []string, error) {
	m := reflect.ValueOf(obj)
	if m.Kind() != reflect.Map {
		return header, knownKeys, make([]string, 0), &ErrInvalidKind{Value: m.Kind(), Valid: []reflect.Kind{reflect.Map}}
	}

	newHeader := make([]string, 0, len(header))
	newKnownKeys := map[string]struct{}{}
	for _, k := range header {
		if k == "*" {
			for _, unknownKey := range unknownKeys(m, knownKeys) {
				newHeader = append(newHeader, unknownKey)
				newKnownKeys[unknownKey] = struct{}{}
			}
			newHeader = append(newHeader, k)
		} else {
			newHeader = append(newHeader, k)
			newKnownKeys[k] = struct{}{}
		}
	}

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

// SerializeBytes serializes an object to its representation given by format.
func SerializeBytes(input *SerializeInput) ([]byte, error) {

	object := input.Object
	format := input.Format
	limit := input.Limit
	valueSerializer := input.ValueSerializer
	if valueSerializer == nil {
		valueSerializer = stringify.DefaultValueStringer("")
	}

	if format == "csv" || format == "tsv" {

		separator, err := sv.FormatToSeparator(format)
		if err != nil {
			return make([]byte, 0), err
		}

		header := input.Header
		wildcard := false
		knownKeys := map[string]struct{}{}
		if len(header) > 0 {
			for _, k := range header {
				if k == "*" {
					wildcard = true
				} else {
					knownKeys[k] = struct{}{}
				}
			}
			if input.Sorted {
				sort.Strings(header)
			}
		}

		s := reflect.ValueOf(object)
		if s.Kind() != reflect.Map && s.Kind() != reflect.Array && s.Kind() != reflect.Slice {
			return make([]byte, 0), &ErrInvalidKind{Value: s.Kind(), Valid: []reflect.Kind{reflect.Array, reflect.Slice}}
		}

		switch s.Kind() {
		case reflect.Map:

			rows := make([][]string, 0)
			if len(header) > 0 {
				if wildcard {
					newHeader, _, row, err := serializeRow(header, knownKeys, s.Interface())
					if err != nil {
						return make([]byte, 0), errors.Wrap(err, "error serializing row")
					}
					header = newHeader
					rows = append(rows, row)
				} else {
					row, err := ToRowS(header, s, valueSerializer)
					if err != nil {
						return make([]byte, 0), errors.Wrap(err, "error serializing row")
					}
					rows = append(rows, row)
				}
			} else {
				keys := inspector.GetKeysFromValue(s, input.Sorted)
				header = ToStringSlice(keys) // string representations of keys
				row, err := ToRowI(keys, s, valueSerializer)
				if err != nil {
					return make([]byte, 0), errors.Wrap(err, "error serializing row")
				}
				rows = append(rows, row)
			}

			buf := new(bytes.Buffer)
			err := sv.Write(&sv.WriteInput{
				Writer:    buf,
				Separator: separator,
				Header:    header,
				Rows:      rows,
			})
			return buf.Bytes(), err
		case reflect.Array, reflect.Slice:
			if s.Len() > 0 {
				first := reflect.TypeOf(s.Index(0).Interface())
				if first.Kind() == reflect.Ptr {
					first = first.Elem()
				}
				rows := make([][]string, 0)
				switch first.Kind() {
				case reflect.Map:
					mapKeys := reflect.ValueOf(s.Index(0).Interface()).MapKeys()
					if len(header) == 0 {
						header = make([]string, 0, len(mapKeys))
						for _, key := range mapKeys {
							header = append(header, fmt.Sprint(key))
						}
						sort.Strings(header)
					}
					if wildcard {
						for i := 0; i < s.Len() && (limit < 0 || i <= limit); i++ {
							newHeader, newKnownKeys, row, err := serializeRow(header, knownKeys, s.Index(i).Interface())
							if err != nil {
								return make([]byte, 0), errors.Wrap(err, "error serializing row")
							}
							header = newHeader
							knownKeys = newKnownKeys
							rows = append(rows, row)
						}
					} else {
						for i := 0; i < s.Len() && (limit < 0 || i <= limit); i++ {
							m := reflect.ValueOf(s.Index(i).Interface())
							if m.Kind() != reflect.Map {
								return make([]byte, 0), &ErrInvalidKind{Value: m.Kind(), Valid: []reflect.Kind{reflect.Map}}
							}
							row := make([]string, len(header))
							for j, key := range header {
								if v := m.MapIndex(reflect.ValueOf(key)); v.IsValid() && !v.IsNil() {
									str, err := valueSerializer(v.Interface())
									if err != nil {
										return make([]byte, 0), errors.Wrap(err, "error serializing value")
									}
									row[j] = str
								} else {
									str, err := valueSerializer(nil)
									if err != nil {
										return make([]byte, 0), errors.Wrap(err, "error serializing value")
									}
									row[j] = str
								}
							}
							rows = append(rows, row)
						}
					}
				case reflect.Struct:
					if len(header) == 0 {
						header = make([]string, first.NumField())
						for i := 0; i < first.NumField(); i++ {
							if tag := first.Field(i).Tag.Get(format); len(tag) > 0 {
								header[i] = tag
							} else {
								header[i] = first.Field(i).Name
							}
						}
						sort.Strings(header)
					}
					for i := 0; i < s.Len() && (limit < 0 || i <= limit); i++ {
						rows = append(rows, make([]string, first.NumField()))
						for j := 0; j < len(header); j++ {
							f, ok := first.FieldByName(header[j])
							if ok {
								switch f.Type {
								case reflect.TypeOf(""):
									rows[i][j] = s.Index(i).Elem().FieldByName(header[j]).String()
								case reflect.TypeOf([]string{}):
									rows[i][j] = strings.Join(s.Index(i).Elem().FieldByName(header[j]).Interface().([]string), ",")
								default:
									rows[i][j] = ""
								}
							} else {
								rows[i][j] = ""
							}
						}
					}
				}
				buf := new(bytes.Buffer)
				err := sv.Write(&sv.WriteInput{
					Writer:    buf,
					Separator: separator,
					Header:    header,
					Rows:      rows,
				})
				return buf.Bytes(), err
			}
			// If there are no records then just return an empty string
			return []byte(""), nil
		}
		return make([]byte, 0), &ErrInvalidKind{Value: s.Kind(), Valid: []reflect.Kind{reflect.Map, reflect.Array, reflect.Slice}}
	} else if format == "properties" {
		buf := new(bytes.Buffer)
		err := properties.Write(&properties.WriteInput{
			Writer:            buf,
			LineSeparator:     input.LineSeparator,
			KeyValueSeparator: input.KeyValueSeparator,
			Object:            object,
			ValueSerializer:   valueSerializer,
			Sorted:            input.Sorted,
			EscapePrefix:      input.EscapePrefix,
			EscapeSpace:       input.EscapeSpace,
			EscapeColon:       false,
			EscapeNewLine:     input.EscapeNewLine,
			EscapeEqual:       input.EscapeEqual,
		})
		if err != nil {
			return make([]byte, 0), errors.Wrap(err, "error writing properties")
		}
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
	} else if format == "bson" {
		return bson.Marshal(stringify.StringifyMapKeys(object))
	} else if format == "json" {
		return json.Marshal(object, input.Pretty)
	} else if format == "jsonl" {
		buf := new(bytes.Buffer)
		errorWrite := jsonl.Write(&jsonl.WriteInput{
			Writer:        buf,
			LineSeparator: []byte(input.LineSeparator)[0],
			Object:        object,
		})
		if errorWrite != nil {
			return nil, errors.Wrap(errorWrite, "error writing json lines")
		}
		return buf.Bytes(), nil
	} else if format == "hcl" {
		return make([]byte, 0), errors.New("Error cannot serialize to HCL")
	} else if format == "hcl2" {
		return make([]byte, 0), errors.New("Error cannot serialize to HCL2")
	} else if format == "toml" {
		return toml.Marshal(object)
	} else if format == "yaml" {
		return yaml.Marshal(object)
	} else if format == "golang" || format == "go" {
		return []byte(fmt.Sprint(object)), nil
	}
	return make([]byte, 0), errors.Wrap(&ErrUnknownFormat{Name: format}, "could not serialize object")
}
