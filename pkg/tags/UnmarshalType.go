// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package tags

import (
	"reflect"
	"strings"
	"unicode/utf8"

	"github.com/spatialcurrent/go-simple-serializer/pkg/escaper" // utf8 is used to decode the first rune in the string
)

// UnmarshalType parses a slice of bytes into an object of a given type.
// If no input is given, then returns ErrEmptyInput.
// If the first rune in the slice of bytes is invalid, then returns ErrInvalidRune.
func UnmarshalType(b []byte, keyValueSeparator rune, outputType reflect.Type) (interface{}, error) {

	if len(b) == 0 {
		return nil, ErrEmptyInput
	}

	// If the kind of the output type is interface{}, then simply use Unmarshal.
	if outputType.Kind() == reflect.Interface {
		return Unmarshal(b, keyValueSeparator)
	}

	first, _ := utf8.DecodeRune(b)
	if first == utf8.RuneError {
		return nil, ErrInvalidRune
	}

	e := escaper.New().Prefix("\\").Sub("\"", "\n")

	switch outputType.Kind() {
	case reflect.Map:
		m := reflect.MakeMap(outputType)
		key := ""
		quotes := 0
		str := ""
		for i, c := range string(b) {
			if quotes == 0 {
				switch c {
				case quote:
					quotes++
				case keyValueSeparator:
					if len(key) == 0 {
						key = str
						str = ""
					}
				case space:
					if len(key) > 0 {
						m.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(e.Unescape(str)))
					}
					key = ""
					str = ""
				default:
					str += string(c)
				}
			} else if quotes == 1 {
				switch c {
				case quote:
					// if the previous character is an escape character
					if b[i-1] == '\\' {
						str += string(c)
					} else {
						quotes--
					}
				default:
					str += string(c)
				}
			}
		}

		if len(key) > 0 {
			m.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(strings.TrimSpace(e.Unescape(str))))
		}
		return m.Interface(), nil
	}

	return nil, &ErrInvalidKind{Value: outputType, Expected: []reflect.Kind{reflect.Map}}
}
