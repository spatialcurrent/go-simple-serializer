// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package iterator

import (
	"fmt"
	"io"
	"strings"
)

// This example shows how to use an iterate through an input of JSON Lines.
func ExampleIterator_jsonl() {
	text := `
  {"a": "b"}
  {"c": "d"}
  {"e": "f"}
  false
  true
  "foo"
  "bar"
  `

	it, err := NewIterator(&NewIteratorInput{
		Reader:            strings.NewReader(text),
		Format:            "jsonl",
		SkipLines:         0,
		Comment:           "",
		Trim:              true,
		SkipBlanks:        false,
		SkipComments:      false,
		KeyValueSeparator: "=",
		LineSeparator:     []byte("\n")[0],
		DropCR:            true,
	})
	if err != nil {
		panic(err)
	}

	for {
		obj, err := it.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		fmt.Println(obj)
	}
	// Output: <nil>
	//map[a:b]
	//map[c:d]
	//map[e:f]
	//false
	//true
	//foo
	//bar
	//<nil>
}

// This example shows how to use an iterate through an input of a lines of tags.
func ExampleIterator_tags() {
	text := `
  a=b x=y
  c=d y=z
  e=f h=i
  `

	it, err := NewIterator(&NewIteratorInput{
		Reader:            strings.NewReader(text),
		Format:            "tags",
		SkipLines:         0,
		Comment:           "",
		Trim:              true,
		SkipBlanks:        false,
		SkipComments:      false,
		KeyValueSeparator: "=",
		LineSeparator:     []byte("\n")[0],
		DropCR:            true,
	})
	if err != nil {
		panic(err)
	}

	for {
		obj, err := it.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		fmt.Println(obj)
	}
	// Output: <nil>
	//map[a:b x:y]
	//map[c:d y:z]
	//map[e:f h:i]
	//<nil>
}
