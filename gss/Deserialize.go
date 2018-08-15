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
func Deserialize(input string, format string, input_header []string, input_comment string, output interface{}) error {

	if format == "csv" || format == "tsv" {
		switch output_slice := output.(type) {
		case *[]map[string]string:
			reader := csv.NewReader(strings.NewReader(input))
			if format == "tsv" {
				reader.Comma = '\t'
			}
			if len(input_comment) > 1 {
				return errors.New("go's encoding/csv package only supports single character comment characters")
			} else if len(input_comment) == 1 {
				reader.Comment = []rune(input_comment)[0]
			}
			if len(input_header) == 0 {
				h, err := reader.Read()
				if err != nil {
					if err != io.EOF {
						return errors.Wrap(err, "Error reading header from input with format csv")
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
						return errors.Wrap(err, "Error reading row from input with format csv")
					}
				}
				*output_slice = append(*output_slice, RowToMapOfStrings(input_header, inRow))
			}
		case *[]map[string]interface{}:
			reader := csv.NewReader(strings.NewReader(input))
			if format == "tsv" {
				reader.Comma = '\t'
			}
			if len(input_comment) > 1 {
				return errors.New("go's encoding/csv package only supports single character comment characters")
			} else if len(input_comment) == 1 {
				reader.Comment = []rune(input_comment)[0]
			}
			if len(input_header) == 0 {
				h, err := reader.Read()
				if err != nil {
					if err != io.EOF {
						return errors.Wrap(err, "Error reading header from input with format csv")
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
						return errors.Wrap(err, "Error reading row from input with format csv")
					}
				}
				*output_slice = append(*output_slice, RowToMapOfInterfaces(input_header, inRow))
			}
		default:
			return errors.New("Cannot deserialize to type " + fmt.Sprint(reflect.ValueOf(output)))
		}
	} else if format == "properties" {
		m := reflect.ValueOf(output)
		if m.Kind() == reflect.Ptr {
			if m.Elem().Kind() != reflect.Map {
				return errors.New("Output is not of kind map.")
			}
			m = m.Elem()
		} else if m.Kind() != reflect.Map {
			return errors.New("Output is not of kind map.")
		}
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
						return errors.New("error deserializing properties for property " + property)
					}
					m.SetMapIndex(reflect.ValueOf(unescapePropertyText(strings.TrimSpace(propertyName))), reflect.ValueOf(unescapePropertyText(strings.TrimSpace(propertyValue))))
					property = ""
				}
			}
		}
	} else if format == "bson" {
		return bson.Unmarshal([]byte(input), output)
	} else if format == "json" {
		return json.Unmarshal([]byte(input), output)
	} else if format == "jsonl" {
		switch output.(type) {
		case *[]map[string]string:
			output_slice := output.(*[]map[string]string)
			scanner := bufio.NewScanner(strings.NewReader(input))
			scanner.Split(bufio.ScanLines)
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				if len(input_comment) == 0 || !strings.HasPrefix(line, input_comment) {
					obj := map[string]string{}
					err := json.Unmarshal([]byte(line), &obj)
					if err != nil {
						return errors.Wrap(err, "Error reading object from JSON line")
					}
					*output_slice = append(*output_slice, obj)
				}
			}
		case *[]map[string]interface{}:
			output_slice := output.(*[]map[string]interface{})
			scanner := bufio.NewScanner(strings.NewReader(input))
			scanner.Split(bufio.ScanLines)
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				if len(input_comment) == 0 || !strings.HasPrefix(line, input_comment) {
					obj := map[string]interface{}{}
					err := json.Unmarshal([]byte(line), &obj)
					if err != nil {
						return errors.Wrap(err, "Error reading object from JSON line")
					}
					*output_slice = append(*output_slice, obj)
				}
			}
		default:
			return errors.New("Cannot deserialize to type " + fmt.Sprint(reflect.ValueOf(output)))
		}
	} else if format == "hcl" {
		obj, err := hcl.Parse(input)
		if err != nil {
			return errors.Wrap(err, "Error parsing hcl")
		}
		if err := hcl.DecodeObject(output, obj); err != nil {
			return errors.Wrap(err, "Error decoding hcl")
		}
	} else if format == "hcl2" {
		file, diags := hclsyntax.ParseConfig([]byte(input), "<stdin>", hcl2.Pos{Byte: 0, Line: 1, Column: 1})
		if diags.HasErrors() {
			return errors.Wrap(errors.New(diags.Error()), "Error parsing hcl2")
		}
		output = &file.Body
	} else if format == "toml" {
		_, err := toml.Decode(input, output)
		return err
	} else if format == "yaml" {
		return yaml.Unmarshal([]byte(input), output)
	}
	return nil
}
