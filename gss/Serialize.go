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

// Serialize serializes an object to its representation given by format.
func Serialize(input interface{}, format string) (string, error) {
	output := ""

	if format == "csv" || format == "tsv" {
		s := reflect.ValueOf(input)
		if s.Kind() != reflect.Slice {
			return "", errors.New("Input is not of kind slice.")
		}
		if s.Len() > 0 {
			first := s.Index(0).Type()
			if first.Kind() == reflect.Ptr {
				first = first.Elem()
			}
			//firstValue := reflect.ValueOf(first)
			rows := make([][]string, 0)
			header := make([]string, 0)
			fmt.Println("first.Kind():", first.Kind())
			//fmt.Println("firstValue.Kind():", firstValue.Kind())
			switch first.Kind() {
			case reflect.Map:
				mapKeys := s.Index(0).MapKeys()
				header = make([]string, 0, len(mapKeys))
				for _, key := range mapKeys {
					header = append(header, fmt.Sprint(key))
				}
				for i := 0; i < s.Len(); i++ {
					rows = append(rows, make([]string, len(mapKeys)))
					for j, key := range mapKeys {
						m := s.Index(i)
						if m.Kind() != reflect.Map {
							return "", errors.New("Row is not of kind map.")
						}
						rows[i][j] = fmt.Sprint(m.MapIndex(key).Interface())
					}
				}
			case reflect.Struct:
				header = make([]string, first.NumField())
				for i := 0; i < first.NumField(); i++ {
					if tag := first.Field(i).Tag.Get(format); len(tag) > 0 {
						header[i] = tag
					} else {
						header[i] = first.Field(i).Name
					}
				}
				for i := 0; i < s.Len(); i++ {
					rows = append(rows, make([]string, first.NumField()))
					for j := 0; j < len(header); j++ {
						ft := first.Field(j).Type
						switch {
						case ft == reflect.TypeOf(""):
							rows[i][j] = s.Index(i).Elem().Field(j).String()
						case ft == reflect.TypeOf([]string{}):
							slice_as_strings := s.Index(i).Elem().Field(j).Interface().([]string)
							rows[i][j] = strings.Join(slice_as_strings, ",")
						default:
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
			w.Write(header)
			w.WriteAll(rows)
			output = buf.String()
		}
	} else if format == "properties" {
		m := reflect.ValueOf(input)
		if m.Kind() == reflect.Map {
			output = ""
			for i, key := range m.MapKeys() {
				output += escapePropertyText(fmt.Sprint(key)) + "=" + escapePropertyText(fmt.Sprint(m.MapIndex(key).Interface()))
				if i < m.Len()-1 {
					output += "\n"
				}
			}
		} else {
			switch input.(type) {
			case string:
				output = input.(string)
			case int:
				output = strconv.Itoa(input.(int))
			case float64:
				output = strconv.FormatFloat(input.(float64), 'f', -1, 64)
			default:
				return "", errors.New("Input is not of kind map but " + fmt.Sprint(reflect.TypeOf(input)))
			}
		}
	} else if format == "bson" {
		text, err := bson.Marshal(StringifyMapKeys(input))
		if err != nil {
			return "", err
		}
		output = string(text)
	} else if format == "json" {
		text, err := json.Marshal(StringifyMapKeys(input))
		if err != nil {
			return "", err
		}
		output = string(text)
	} else if format == "jsonl" {
		s := reflect.ValueOf(input)
		if s.Kind() != reflect.Slice {
			return "", errors.New("Input is not of kind slice.")
		}
		for i := 0; i < s.Len(); i++ {
			text, err := json.Marshal(s.Index(i).Interface())
			if err != nil {
				return "", err
			}
			output += string(text)
			if i < s.Len()-1 {
				output += "\n"
			}
		}
	} else if format == "hcl" {
		return "", errors.New("Error cannot serialize to HCL")
	} else if format == "hcl2" {
		return "", errors.New("Error cannot serialize to HCL2")
	} else if format == "toml" {
		buf := new(bytes.Buffer)
		if err := toml.NewEncoder(buf).Encode(input); err != nil {
			return "", errors.Wrap(err, "Error encoding TOML")
		}
		output = buf.String()
	} else if format == "yaml" {
		y, err := yaml.Marshal(input)
		if err != nil {
			return "", err
		}
		output = string(y)
	}
	return output, nil
}
