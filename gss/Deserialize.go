package gss

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"
)

import (
	"github.com/BurntSushi/toml"
	"github.com/hashicorp/hcl"
	hcl2 "github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hcl/hclsyntax"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// Deserialize reads in an object from a string given format
func Deserialize(input string, format string, output interface{}) error {

	if format == "csv" {
		switch output.(type) {
		case *[]map[string]string:
			output_slice := output.(*[]map[string]string)
			reader := csv.NewReader(strings.NewReader(input))
			header, err := reader.Read()
			if err != nil {
				if err != io.EOF {
					return errors.Wrap(err, "Error reading header from input with format csv")
				}
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

				inObj := map[string]string{}
				for i, h := range header {
					inObj[strings.ToLower(h)] = inRow[i]
				}
				*output_slice = append(*output_slice, inObj)
			}
		case *[]map[string]interface{}:
			output_slice := output.(*[]map[string]interface{})
			reader := csv.NewReader(strings.NewReader(input))
			header, err := reader.Read()
			if err != nil {
				if err != io.EOF {
					return errors.Wrap(err, "Error reading header from input with format csv")
				}
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

				inObj := map[string]interface{}{}
				for i, h := range header {
					inObj[strings.ToLower(h)] = inRow[i]
				}
				*output_slice = append(*output_slice, inObj)
			}
		default:
			return errors.New("Cannot deserialize to type " + fmt.Sprint(reflect.ValueOf(output)))
		}
	} else if format == "json" {
		return json.Unmarshal([]byte(input), output)
	} else if format == "jsonl" {
		switch output.(type) {
		case *[]map[string]string:
			output_slice := output.(*[]map[string]string)
			scanner := bufio.NewScanner(strings.NewReader(input))
			scanner.Split(bufio.ScanLines)
			for scanner.Scan() {
				obj := map[string]string{}
				err := json.Unmarshal([]byte(scanner.Text()), &obj)
				if err != nil {
					return errors.Wrap(err, "Error reading object from JSON line")
				}
				*output_slice = append(*output_slice, obj)
			}
		case *[]map[string]interface{}:
			output_slice := output.(*[]map[string]interface{})
			scanner := bufio.NewScanner(strings.NewReader(input))
			scanner.Split(bufio.ScanLines)
			for scanner.Scan() {
				obj := map[string]interface{}{}
				err := json.Unmarshal([]byte(scanner.Text()), &obj)
				if err != nil {
					return errors.Wrap(err, "Error reading object from JSON line")
				}
				*output_slice = append(*output_slice, obj)
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
		return yaml.Unmarshal([]byte(input), &output)
	}
	return nil
}
