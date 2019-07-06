// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package tags

import (
	"unicode/utf8" // utf8 is used to decode the first rune in the string
)

import (
	"github.com/spatialcurrent/go-simple-serializer/pkg/escaper"
)

// Unmarshal parses a slice of bytes into an object using a few simple type inference rules.
// This package is useful when your program needs to parse data,
// that you have no a priori awareness of its structure or type.
// If no input is given, then returns ErrEmptyInput.
// If the first rune is invalid, then returns ErrInvalidRune.
func Unmarshal(b []byte) (interface{}, error) {

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
			case '"':
				quotes++
			case '=':
				if len(key) == 0 {
					key = str
					str = ""
				}
			case ' ':
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
			case '"':
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
		obj[key] = e.Unescape(str)
	}

	return obj, nil
}
