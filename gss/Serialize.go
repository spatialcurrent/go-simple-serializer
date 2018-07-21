package gss

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"reflect"
	"strings"
)

import (
	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// Serialize serializes an object to its representation given by format.
func Serialize(input interface{}, format string) (string, error) {
	output := ""

	if format == "csv" {
		s := reflect.ValueOf(input)
		if s.Kind() != reflect.Slice {
			return "", errors.New("Input is not of kind slice.")
		}
		if s.Len() > 0 {
			first := s.Index(0).Type()
			if first.Kind() == reflect.Ptr {
				first = first.Elem()
			}
			rows := make([][]string, 0)
			header := make([]string, first.NumField())
			for i := 0; i < first.NumField(); i++ {
				if tag := first.Field(i).Tag.Get("csv"); len(tag) > 0 {
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
			buf := new(bytes.Buffer)
			w := csv.NewWriter(buf)
			w.Write(header)
			w.WriteAll(rows)
			output = buf.String()
		}
	} else if format == "json" {
		text, err := json.Marshal(input)
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
