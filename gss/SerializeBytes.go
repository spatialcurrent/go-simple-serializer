// =================================================================
//
// Copyright (C) 2018 Spatial Current, Inc. - All Rights Reserved
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

func escapePropertyText(in string) string {
	out := in
	out = strings.Replace(fmt.Sprint(out), "\\", "\\\\", -1)
	out = strings.Replace(fmt.Sprint(out), " ", "\\ ", -1)
	return out
}

// SerializeBytes serializes an object to its representation given by format.
func SerializeBytes(input interface{}, format string, header []string, limit int) ([]byte, error) {

	if format == "csv" || format == "tsv" {
		s := reflect.ValueOf(input)
		if s.Kind() != reflect.Slice {
			return make([]byte, 0), errors.New("Input is not of kind slice.")
		}
		if s.Len() > 0 {
			first := s.Index(0).Type()
			if first.Kind() == reflect.Ptr {
				first = first.Elem()
			}
			rows := make([][]string, 0)
			switch first.Kind() {
			case reflect.Map:
				mapKeys := s.Index(0).MapKeys()
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
						m := s.Index(i)
						if m.Kind() != reflect.Map {
							return make([]byte, 0), errors.New("Row is not of kind map.")
						}
						rows[i][j] = fmt.Sprint(m.MapIndex(reflect.ValueOf(key)).Interface())
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
	} else if format == "properties" {
		t := reflect.TypeOf(input)
		if t.Kind() == reflect.Map {
			if t.Key().Kind() != reflect.String {
				return nil, errors.New("can only serialize a map with string keys")
			}
			m := reflect.ValueOf(input)
			keys := make([]string, m.Len())
			for i, key := range m.MapKeys() {
				keys[i] = key.Interface().(string)
			}
			sort.Strings(keys)
			output := ""
			for i, key := range keys {
				output += escapePropertyText(key) + "=" + escapePropertyText(fmt.Sprint(m.MapIndex(reflect.ValueOf(key)).Interface()))
				if i < m.Len()-1 {
					output += "\n"
				}
			}
			return []byte(output), nil
		} else {
			switch input.(type) {
			case string:
				return []byte(input.(string)), nil
			case int:
				return []byte(strconv.Itoa(input.(int))), nil
			case float64:
				return []byte(strconv.FormatFloat(input.(float64), 'f', -1, 64)), nil
			}
			return make([]byte, 0), errors.New("Input is not of kind map but " + fmt.Sprint(reflect.TypeOf(input)))
		}
	} else if format == "bson" {
		return bson.Marshal(StringifyMapKeys(input))
	} else if format == "json" {
		return json.Marshal(StringifyMapKeys(input))
	} else if format == "jsonl" {
		s := reflect.ValueOf(input)
		if s.Kind() != reflect.Slice {
			return make([]byte, 0), errors.New("Input is not of kind slice.")
		}
		output := make([]byte, 0)
		for i := 0; i < s.Len() && (limit < 0 || i < limit); i++ {
			b, err := json.Marshal(s.Index(i).Interface())
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
		if err := toml.NewEncoder(buf).Encode(input); err != nil {
			return make([]byte, 0), errors.Wrap(err, "Error encoding TOML")
		}
		return buf.Bytes(), nil
	} else if format == "yaml" {
		return yaml.Marshal(input)
	}
	return make([]byte, 0), errors.Wrap(&ErrUnknownFormat{Name: format}, "could not serialize object")
}
