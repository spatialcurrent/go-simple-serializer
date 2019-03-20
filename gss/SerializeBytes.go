// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gss

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

import (
	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/yaml.v2"
)

var jsonPrefix = ""
var jsonIndent = "  "

func escapePropertyText(in string) string {
	out := in
	out = strings.Replace(fmt.Sprint(out), "\\", "\\\\", -1)
	out = strings.Replace(fmt.Sprint(out), " ", "\\ ", -1)
	return out
}

func marshalJson(object interface{}, pretty bool) ([]byte, error) {
	if pretty {
		return json.MarshalIndent(StringifyMapKeys(object), jsonPrefix, jsonIndent)
	}
	return json.Marshal(StringifyMapKeys(object))
}

// SerializeBytes serializes an object to its representation given by format.
func SerializeBytes(input *SerializeInput) ([]byte, error) {

	object := input.Object
	format := input.Format
	limit := input.Limit

	if format == "csv" || format == "tsv" {
		header := input.Header

		s := reflect.ValueOf(object)
		if s.Kind() != reflect.Array && s.Kind() != reflect.Slice {
			return make([]byte, 0), &ErrInvalidKind{Value: s.Kind(), Valid: []reflect.Kind{reflect.Array, reflect.Slice}}
		}
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
				for i := 0; i < s.Len() && (limit < 0 || i <= limit); i++ {
					rows = append(rows, make([]string, len(header)))
					for j, key := range header {
						m := reflect.ValueOf(s.Index(i).Interface())
						if m.Kind() != reflect.Map {
							return make([]byte, 0), &ErrInvalidKind{Value: m.Kind(), Valid: []reflect.Kind{reflect.Map}}
						}
						if v := m.MapIndex(reflect.ValueOf(key)); v.IsValid() && !v.IsNil() {
							rows[i][j] = fmt.Sprint(v.Interface())
						} else {
							rows[i][j] = ""
						}
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
								slice_as_strings := s.Index(i).Elem().FieldByName(header[j]).Interface().([]string)
								rows[i][j] = strings.Join(slice_as_strings, ",")
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
			w := csv.NewWriter(buf)
			if format == "tsv" {
				w.Comma = '\t'
			}
			w.Write(header)  // nolint: gosec
			w.WriteAll(rows) // nolint: gosec
			return buf.Bytes(), nil
		}
		// If there are no records then just return an empty string
		return []byte(""), nil
	} else if format == "properties" || format == "text" {
		t := reflect.TypeOf(object)
		if t.Kind() == reflect.Map {
			if k := t.Key().Kind(); k != reflect.String {
				return nil, errors.Wrap(&ErrInvalidKind{Value: k, Valid: []reflect.Kind{reflect.String}}, "can only serialize a map with string keys")
			}
			m := reflect.ValueOf(object)
			keys := make([]string, m.Len())
			for i, key := range m.MapKeys() {
				keys[i] = key.Interface().(string)
			}
			sort.Strings(keys)
			output := ""
			for i, key := range keys {
				if format == "properties" {
					output += escapePropertyText(key) + "=" + escapePropertyText(fmt.Sprint(m.MapIndex(reflect.ValueOf(key)).Interface()))
					if i < m.Len()-1 {
						output += "\n"
					}
				} else if format == "text" {
					value := strings.Replace(fmt.Sprint(m.MapIndex(reflect.ValueOf(key)).Interface()), "\"", "\\\"", -1)
					if strings.Contains(value, " ") {
						value = "\"" + value + "\""
					}
					output += key + "=" + value
					if i < m.Len()-1 {
						output += " "
					}
				}

			}
			return []byte(output), nil
		}
		switch object.(type) {
		case string:
			return []byte(object.(string)), nil
		case int:
			return []byte(strconv.Itoa(object.(int))), nil
		case float64:
			return []byte(strconv.FormatFloat(object.(float64), 'f', -1, 64)), nil
		}
		return make([]byte, 0), errors.Wrap(&ErrInvalidKind{Value: reflect.TypeOf(object).Kind(), Valid: []reflect.Kind{reflect.Map, reflect.String, reflect.Int, reflect.Float64}}, "object is not valid")
	} else if format == "bson" {
		return bson.Marshal(StringifyMapKeys(object))
	} else if format == "json" {
		return marshalJson(object, input.Pretty)
	} else if format == "jsonl" {
		s := reflect.ValueOf(object)
		if s.Kind() != reflect.Array && s.Kind() != reflect.Slice {
			return make([]byte, 0), errors.Wrap(&ErrInvalidKind{Value: reflect.TypeOf(object).Kind(), Valid: []reflect.Kind{reflect.Array, reflect.Slice}}, "object is not valid")
		}
		output := make([]byte, 0)
		for i := 0; i < s.Len() && (limit < 0 || i < limit); i++ {
			b, err := marshalJson(s.Index(i).Interface(), input.Pretty)
			if err != nil {
				return output, err
			}
			output = append(output, b...)
			if i < s.Len()-1 {
				output = append(output, []byte("\n")[0])
			}
		}
		return output, nil
	} else if format == "hcl" {
		return make([]byte, 0), errors.New("Error cannot serialize to HCL")
	} else if format == "hcl2" {
		return make([]byte, 0), errors.New("Error cannot serialize to HCL2")
	} else if format == "toml" {
		buf := new(bytes.Buffer)
		if err := toml.NewEncoder(buf).Encode(object); err != nil {
			return make([]byte, 0), errors.Wrap(err, "Error encoding TOML")
		}
		return buf.Bytes(), nil
	} else if format == "yaml" {
		return yaml.Marshal(object)
	} else if format == "golang" || format == "go" {
		return []byte(fmt.Sprint(object)), nil
	}
	return make([]byte, 0), errors.Wrap(&ErrUnknownFormat{Name: format}, "could not serialize object")
}
