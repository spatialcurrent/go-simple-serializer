// =================================================================
//
// Copyright (C) 2018 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gss

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"
	"unicode"
)

import (
	"github.com/BurntSushi/toml"
	"github.com/hashicorp/hcl"
	hcl2 "github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hcl/hclsyntax"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/yaml.v2"
)

func unescapePropertyText(in string) string {
	out := in
	out = strings.Replace(fmt.Sprint(out), "\\ ", " ", -1)
	out = strings.Replace(fmt.Sprint(out), "\\\\", "\\", -1)
	return out
}

// Deserialize reads in an object from a string given format
func Deserialize(input string, format string, input_header []string, input_comment string, output_type reflect.Type, verbose bool) (interface{}, error) {

	if format == "csv" || format == "tsv" {
		output := reflect.MakeSlice(output_type, 0, 0)
		reader := csv.NewReader(strings.NewReader(input))
		if format == "tsv" {
			reader.Comma = '\t'
		}
		if len(input_comment) > 1 {
			return nil, errors.New("go's encoding/csv package only supports single character comment characters")
		} else if len(input_comment) == 1 {
			reader.Comment = []rune(input_comment)[0]
		}
		if len(input_header) == 0 {
			h, err := reader.Read()
			if err != nil {
				if err != io.EOF {
					return nil, errors.Wrap(err, "Error reading header from input with format csv")
				}
			}
			input_header = h
		}
		for {
			inRow, err := reader.Read()
			if err != nil {
				if err == io.EOF {
					break
				} else {
					return nil, errors.Wrap(err, "Error reading row from input with format csv")
				}
			}
			m := reflect.MakeMap(reflect.TypeOf(output.Elem()))
			for i, h := range input_header {
				m.SetMapIndex(reflect.ValueOf(strings.ToLower(h)), reflect.ValueOf(inRow[i]))
			}
			output = reflect.Append(output, m)
		}
		return output.Interface(), nil
	} else if format == "properties" {
		m := reflect.MakeMap(output_type)
		if len(input_comment) == 0 {
			input_comment = "#"
		}
		scanner := bufio.NewScanner(strings.NewReader(input))
		scanner.Split(bufio.ScanLines)
		property := ""
		for scanner.Scan() {
			line := scanner.Text()
			if len(line) > 0 && !strings.HasPrefix(line, input_comment) {
				if line[len(line)-1] == '\\' {
					property += strings.TrimLeftFunc(line[0:len(line)-1], unicode.IsSpace)
				} else {
					property += strings.TrimLeftFunc(line, unicode.IsSpace)
					propertyName := ""
					propertyValue := ""
					for i, c := range property {
						if c == '=' || c == ':' {
							propertyName = property[0:i]
							propertyValue = property[i+1:]
							break
						}
					}
					if len(propertyName) == 0 {
						return nil, errors.New("error deserializing properties for property " + property)
					}
					m.SetMapIndex(reflect.ValueOf(unescapePropertyText(strings.TrimSpace(propertyName))), reflect.ValueOf(unescapePropertyText(strings.TrimSpace(propertyValue))))
					property = ""
				}
			}
		}
		return m.Interface(), nil
	} else if format == "bson" {
		if output_type.Kind() == reflect.Map {
			ptr := reflect.New(output_type)
			ptr.Elem().Set(reflect.MakeMap(output_type))
			err := bson.Unmarshal([]byte(input), ptr.Interface())
			return ptr.Elem().Interface(), err
		} else {
			return nil, errors.New("Invalid output type for bson " + fmt.Sprint(output_type))
		}
	} else if format == "json" {
		if output_type.Kind() == reflect.Map {
			ptr := reflect.New(output_type)
			ptr.Elem().Set(reflect.MakeMap(output_type))
			err := json.Unmarshal([]byte(input), ptr.Interface())
			return ptr.Elem().Interface(), err
		} else if output_type.Kind() == reflect.Slice {
			ptr := reflect.New(output_type)
			ptr.Elem().Set(reflect.MakeSlice(output_type, 0, 0))
			err := json.Unmarshal([]byte(input), ptr.Interface())
			return ptr.Elem().Interface(), err
		} else {
			return nil, errors.New("Invalid output type for json " + fmt.Sprint(output_type))
		}
	} else if format == "jsonl" {
		output := reflect.MakeSlice(output_type, 0, 0)
		scanner := bufio.NewScanner(strings.NewReader(input))
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if len(input_comment) == 0 || !strings.HasPrefix(line, input_comment) {
				ptr := reflect.New(output_type.Elem())
				ptr.Elem().Set(reflect.MakeMap(output_type.Elem()))
				err := json.Unmarshal([]byte(line), ptr.Interface())
				if err != nil {
					return nil, errors.Wrap(err, "Error reading object from JSON line")
				}
				output = reflect.Append(output, ptr.Elem())
			}
		}
		return output.Interface(), nil
	} else if format == "hcl" {
		ptr := reflect.New(output_type)
		ptr.Elem().Set(reflect.MakeMap(output_type))
		obj, err := hcl.Parse(input)
		if err != nil {
			return nil, errors.Wrap(err, "Error parsing hcl")
		}
		if err := hcl.DecodeObject(ptr.Interface(), obj); err != nil {
			return nil, errors.Wrap(err, "Error decoding hcl")
		}
		return ptr.Elem().Interface(), nil
	} else if format == "hcl2" {
		file, diags := hclsyntax.ParseConfig([]byte(input), "<stdin>", hcl2.Pos{Byte: 0, Line: 1, Column: 1})
		if diags.HasErrors() {
			return nil, errors.Wrap(errors.New(diags.Error()), "Error parsing hcl2")
		}
		return &file.Body, nil
	} else if format == "toml" {
		if output_type.Kind() == reflect.Map {
			ptr := reflect.New(output_type)
			ptr.Elem().Set(reflect.MakeMap(output_type))
			_, err := toml.Decode(input, ptr.Interface())
			return ptr.Elem().Interface(), err
		} else {
			return nil, errors.New("Invalid output type for toml " + fmt.Sprint(output_type))
		}
	} else if format == "yaml" {
		if output_type.Kind() == reflect.Map {
			ptr := reflect.New(output_type)
			ptr.Elem().Set(reflect.MakeMap(output_type))
			err := yaml.Unmarshal([]byte(input), ptr.Interface())
			return ptr.Elem().Interface(), err
		} else if output_type.Kind() == reflect.Slice {
			ptr := reflect.New(output_type)
			ptr.Elem().Set(reflect.MakeSlice(output_type, 0, 0))
			err := yaml.Unmarshal([]byte(input), ptr.Interface())
			return StringifyMapKeys(ptr.Elem().Interface()), err
		} else {
			return nil, errors.New("Invalid output type for yaml " + fmt.Sprint(output_type))
		}
	}

	return nil, nil
}
