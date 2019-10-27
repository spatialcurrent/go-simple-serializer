// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package yaml

// ParseKeyValue parses a stream of bytes into a key and value.
// Returns true if a key and value if found.  Returns false if could not parse.
func ParseKeyValue(line []byte) ([]byte, []byte, bool) {
	d := -1
	eol := -1
	for i, c := range line {
		if c == ':' {
			if i == 0 {
				return make([]byte, 0), make([]byte, 0), false
			}
			if line[i-1] != '\\' && (i+1 < len(line) && line[i+1] == ' ') {
				d = i
			}
		} else if c == '\n' {
			eol = i
			break
		}
	}
	if d != -1 {
		if eol != -1 {
			return line[0:d], line[d+2 : eol], true
		}
		return line[0:d], line[d+2:], true
	}
	return make([]byte, 0), make([]byte, 0), false
}
