// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package properties

import (
	"reflect"
	"strings"
	"unicode"

	"github.com/pkg/errors"
	"github.com/spatialcurrent/go-simple-serializer/pkg/escaper"
	"github.com/spatialcurrent/go-simple-serializer/pkg/scanner"
)

// Read parses properties from the given reader and returns a map of the properties, and error if any.
func Read(input *ReadInput) (interface{}, error) {

	inputType := reflect.TypeOf(map[string]string{})
	if input.Type != nil {
		inputType = input.Type
	}

	// Initialize Escaper
	e := escaper.New()
	if len(input.EscapePrefix) > 0 {
		e = e.Prefix(input.EscapePrefix)
		if input.UnescapeSpace {
			e = e.Sub(" ")
		}
		if input.UnescapeEqual {
			e = e.Sub("=")
		}
		if input.UnescapeColon {
			e = e.Sub(":")
		}
		if input.UnescapeNewLine {
			e = e.Sub("\n")
		}
	}

	m := reflect.MakeMap(inputType)
	s := scanner.New(input.Reader, input.LineSeparator, input.DropCR)
	property := ""
	for s.Scan() {
		line := s.Text()
		if input.Trim {
			line = strings.TrimSpace(line)
		}
		if len(line) > 0 && (len(input.Comment) == 0 || !strings.HasPrefix(line, input.Comment)) {
			// If the line ends with a backslash and input.UnescapeNewLine is set to true.
			if line[len(line)-1] == '\\' && input.UnescapeNewLine {
				property += strings.TrimLeftFunc(line, unicode.IsSpace) // include backslash since we unescape later.
			} else {
				property += strings.TrimLeftFunc(line, unicode.IsSpace)
				propertyName := ""
				propertyValue := ""
				for i, c := range property {
					split := false
					if c == '=' {
						split = (!input.UnescapeEqual) || (i == 0) || (property[i-1] != '\\')
					} else if c == ':' {
						split = (!input.UnescapeColon) || (i == 0) || (property[i-1] != '\\')
					} else if c == ' ' {
						split = (!input.UnescapeSpace) || (i == 0) || (property[i-1] != '\\')
					}
					if split {
						propertyName = property[0:i]
						propertyValue = property[i+1:]
						break
					}
				}
				if len(propertyName) == 0 {
					return nil, errors.New("error deserializing properties for property " + property)
				}
				m.SetMapIndex(
					reflect.ValueOf(e.Unescape(strings.TrimSpace(propertyName))),
					reflect.ValueOf(e.Unescape(strings.TrimSpace(propertyValue))),
				)
				property = ""
			}
		}
	}
	return m.Interface(), nil
}
