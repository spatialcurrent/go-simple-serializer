// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package jsonl

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

import (
	"github.com/pkg/errors"
)

type Iterator struct {
	Scanner      *bufio.Scanner
	Comment      string
	Trim         bool
	SkipBlanks   bool
	SkipComments bool
	Limit        int
	Count        int
}

// Input for NewIterator function.
type NewIteratorInput struct {
	Reader       io.Reader
	SkipLines    int
	SkipBlanks   bool
	SkipComments bool
	Comment      string
	Trim         bool
	Limit        int
}

func NewIterator(input *NewIteratorInput) *Iterator {
	scanner := bufio.NewScanner(input.Reader)
	scanner.Split(bufio.ScanLines)
	for i := 0; i < input.SkipLines; i++ {
		if !scanner.Scan() {
			break
		}
	}
	return &Iterator{
		Scanner:      scanner,
		Comment:      input.Comment,
		Trim:         input.Trim,
		SkipBlanks:   input.SkipBlanks,
		SkipComments: input.SkipComments,
		Limit:        input.Limit,
		Count:        0,
	}
}

// Next reads from the underlying reader and returns the next object and error, if any.
// If a blank line is found and SkipBlanks is false, then returns (nil, nil).
// If a commented line is found and SkipComments is false, then returns (nil, nil).
// When finished, returns (nil, io.EOF).
func (it *Iterator) Next() (interface{}, error) {

	// If reached limit, return io.EOF
	if it.Limit > 0 && it.Count >= it.Limit {
		return nil, io.EOF
	}

	// Increment Counter
	it.Count++

	if it.Scanner.Scan() {
		line := it.Scanner.Text()
		if it.Trim {
			line = strings.TrimSpace(line)
		}
		if len(line) == 0 {
			if it.SkipBlanks {
				return it.Next()
			}
			return nil, nil
		}
		if len(it.Comment) > 0 && strings.HasPrefix(line, it.Comment) {
			if it.SkipComments {
				return it.Next()
			}
			return nil, nil
		}
		switch line {
		case "true":
			return true, nil
		case "false":
			return false, nil
		case "null":
			return nil, nil
		}
		switch line[0] {
		case '[':
			obj := make([]interface{}, 0)
			err := json.Unmarshal([]byte(line), &obj)
			if err != nil {
				return nil, errors.Wrap(err, fmt.Sprintf("error unmarshaling JSON line %q", line))
			}
			return obj, nil
		case '{':
			obj := map[string]interface{}{}
			err := json.Unmarshal([]byte(line), &obj)
			if err != nil {
				return nil, errors.Wrap(err, fmt.Sprintf("error unmarshaling JSON line %q", line))
			}
			return obj, nil
		case '"':
			obj := ""
			err := json.Unmarshal([]byte(line), &obj)
			if err != nil {
				return nil, errors.Wrap(err, fmt.Sprintf("error unmarshaling JSON line %q", line))
			}
			return obj, nil
		default:
			obj := 0.0
			err := json.Unmarshal([]byte(line), &obj)
			if err != nil {
				return nil, errors.Wrap(err, fmt.Sprintf("error unmarshaling JSON line %q", line))
			}
			return obj, nil
		}
	}
	return nil, io.EOF
}
