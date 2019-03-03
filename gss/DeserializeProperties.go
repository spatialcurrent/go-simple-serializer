// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gss

import (
	"bufio"
	"github.com/pkg/errors"
	"reflect"
	"strings"
	"unicode"
)

// DeserializeProperties deserializes a properties string into a Go instance.
//  - https://en.wikipedia.org/wiki/.properties
func DeserializeProperties(input string, input_comment string, output_type reflect.Type) (interface{}, error) {
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
}
