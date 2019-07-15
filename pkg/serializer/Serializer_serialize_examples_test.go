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

import (
	"github.com/spatialcurrent/go-stringify/pkg/stringify"
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

func ExampleSerializer_Serialize_go() {
	in := map[string]interface{}{
		"a": "1",
		"b": "2",
		"c": "3",
	}
	s := New(FormatGo)
	out, err := s.Serialize(in)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(out))
	// Output: map[string]interface {}{"a":"1", "b":"2", "c":"3"}
}

func ExampleSerializer_Serialize_json() {
	in := map[interface{}]interface{}{
		"foo": "bar",
	}
	s := New(FormatJSON).Sorted(true)
	out, err := s.Serialize(in)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(out))
	// Output: {"foo":"bar"}
}

func ExampleSerializer_Serialize_jsonl() {
	in := []map[string]interface{}{
		map[string]interface{}{
			"a": "1",
			"b": "2",
			"c": "3",
		},
		map[string]interface{}{
			"a": "4",
			"b": "5",
			"c": "6",
		},
	}
	s := New(FormatJSONL).Limit(-1).LineSeparator("\n")
	out, err := s.Serialize(in)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(out))
	// Output: {"a":"1","b":"2","c":"3"}
	// {"a":"4","b":"5","c":"6"}
}

func ExampleSerializer_Serialize_properties() {
	in := map[string]interface{}{
		"a": "1",
		"b": "2",
		"c": "3",
	}
	s := New(FormatProperties).Sorted(true).LineSeparator("\n").KeyValueSeparator("=")
	out, err := s.Serialize(in)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(out))
	// Output: a=1
	//b=2
	//c=3

}

func ExampleSerializer_Serialize_tags() {
	in := map[interface{}]interface{}{
		"hello": "beautiful world",
	}
	s := New(FormatTags).
		KeyValueSeparator("=").
		LineSeparator("\n").
		ValueSerializer(stringify.NewStringer("", false, false, false)).
		Limit(-1).
		Sorted(true)
	out, err := s.Serialize(in)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(out))
	// Output: hello="beautiful world"
}

func ExampleSerializer_Serialize_toml() {
	in := map[string]interface{}{
		"a": "1",
		"b": "2",
		"c": "3",
	}
	s := New(FormatTOML)
	out, err := s.Serialize(in)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(out))
	// Output: a = "1"
	//b = "2"
	//c = "3"
}

func ExampleSerializer_Serialize_yaml() {
	in := map[string]interface{}{
		"a": "1",
		"b": "2",
		"c": "3",
	}
	s := New(FormatYAML)
	out, err := s.Serialize(in)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(out))
	// Output: a: "1"
	// b: "2"
	// c: "3"
}
