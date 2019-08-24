// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package tags

import (
	"strings"
	"unicode/utf8"

	"github.com/spatialcurrent/go-simple-serializer/pkg/escaper" // utf8 is used to decode the first rune in the string
)

// Unmarshal parses a slice of bytes into an object using a few simple type inference rules.
// This package is useful when your program needs to parse data,
// that you have no a priori awareness of its structure or type.
// If no input is given, then returns ErrEmptyInput.
// If the first rune is invalid, then returns ErrInvalidRune.
func Unmarshal(b []byte, keyValueSeparator rune) (interface{}, error) {

	if len(b) == 0 {
		return nil, ErrEmptyInput
	}

	first, _ := utf8.DecodeRune(b)
	if first == utf8.RuneError {
		return nil, ErrInvalidRune
	}

	e := escaper.New().Prefix("\\").Sub("\"", "\n")

	obj := map[string]string{}

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
					obj[key] = e.Unescape(str)
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
		obj[key] = e.Unescape(strings.TrimSpace(str))
	}

	return obj, nil
}
