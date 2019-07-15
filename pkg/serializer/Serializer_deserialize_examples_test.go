// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package serializer

import (
	"fmt"
)

/*
func ExampleSerializer_csv() {
	in := map[interface{}]interface{}{
		"a": "x",
		"b": "y",
		"c": "z",
	}
	s := New(FormatCSV).ValueSerializer(stringify.NewStringer("", false, false, false))
	out, err := s.Serialize(in)
	if err != nil {
		panic(err)
	}
	fmt.Println(out)
	// Output: "a,b,c\bx,y,z"
}

func ExampleSerializer_tsv() {
	in := map[interface{}]interface{}{
		"a": "x",
		"b": "y",
		"c": "z",
	}
	s := New(FormatTSV).ValueSerializer(stringify.NewStringer("", false, false, false))
	out, err := s.Serialize(in)
	if err != nil {
		panic(err)
	}
	fmt.Println(out)
}
*/

func ExampleSerializer_Deserialize_json() {
	in := `{"foo":"bar"}`
	s := New(FormatJSON)
	out, err := s.Deserialize([]byte(in))
	if err != nil {
		panic(err)
	}
	fmt.Println(out)
	// Output: map[foo:bar]
}

func ExampleSerializer_Deserialize_jsonl() {
	in := `
	{"a": "b"}
  {"c": "d"}
  {"e": "f"}
  false
  true
  "foo"
  "bar"
  `
	s := New(FormatJSONL).LineSeparator("\n").Trim(true)
	out, err := s.Deserialize([]byte(in))
	if err != nil {
		panic(err)
	}
	fmt.Println(out)
	// Output: [<nil> map[a:b] map[c:d] map[e:f] false true foo bar <nil>]
}

func ExampleSerializer_Deserialize_properties() {
	in := "a=x\nb=y\nc=z"
	s := New(FormatProperties).LineSeparator("\n")
	out, err := s.Deserialize([]byte(in))
	if err != nil {
		panic(err)
	}
	fmt.Println(out)
	// Output: map[a:x b:y c:z]
}

func ExampleSerializer_Deserialize_tags() {
	in := "hello=\"beautiful world\""
	s := New(FormatTags).KeyValueSeparator("=").LineSeparator("\n")
	out, err := s.Deserialize([]byte(in))
	if err != nil {
		panic(err)
	}
	fmt.Println(out)
	// Output: [map[hello:beautiful world]]
}

func ExampleSerializer_Deserialize_toml() {
	in := "a = 1.0\nb = 2.0\nc = 3.0\n"
	s := New(FormatTOML)
	out, err := s.Deserialize([]byte(in))
	if err != nil {
		panic(err)
	}
	fmt.Println(out)
	// Output: map[a:1 b:2 c:3]
}

func ExampleSerializer_Deserialize_yaml() {
	in := `a: x
b: "y"
c: z`
	s := New(FormatYAML)
	out, err := s.Deserialize([]byte(in))
	if err != nil {
		panic(err)
	}
	fmt.Println(out)
	// Output: map[a:x b:y c:z]
}
