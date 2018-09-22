// =================================================================
//
// Copyright (C) 2018 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gss

import (
	"reflect"
)

// TestCase contains a case for unit tests.
type TestCase struct {
	String     string       // the object in serialized form
	Object     interface{}  // the object as a go object
	Format     string       // the serialization format
	Header     []string     // if format is a csv or tsv, the names of the columns
	Comment    string       // the line comment prefix
	LazyQuotes bool         // if format is csv or tsv, allow LazyQuotes.
	Limit      int          // if format is a csv, tsv, or jsonl, then limit the number of items processed.
	Type       reflect.Type // // the type of the object when in Go
}

var testCases = []TestCase{
	TestCase{
		String: "a,b\nx,y",
		Format: "csv",
		Type:   reflect.TypeOf([]map[string]interface{}{}),
		Object: []map[string]interface{}{map[string]interface{}{"a": "x", "b": "y"}},
	},
	TestCase{
		String: "a\tb\nx\ty",
		Format: "tsv",
		Type:   reflect.TypeOf([]map[string]interface{}{}),
		Object: []map[string]interface{}{map[string]interface{}{"a": "x", "b": "y"}},
	},
	TestCase{
		String: "{\"a\":\"x\",\"b\":\"y\"}",
		Format: "json",
		Type:   reflect.TypeOf(map[string]interface{}{}),
		Object: map[string]interface{}{"a": "x", "b": "y"},
	},
	TestCase{
		String: "[{\"a\":\"x\",\"b\":\"y\"}]",
		Format: "json",
		Type:   reflect.TypeOf([]map[string]interface{}{}),
		Object: []map[string]interface{}{map[string]interface{}{"a": "x", "b": "y"}},
	},
	TestCase{
		String: "{\"a\":\"x\"}\n{\"b\":\"y\"}",
		Format: "jsonl",
		Limit:  -1,
		Type:   reflect.TypeOf([]map[string]interface{}{}),
		Object: []map[string]interface{}{map[string]interface{}{"a": "x"}, map[string]interface{}{"b": "y"}},
	},
	TestCase{
		String: "a=x\nb=y",
		Format: "properties",
		Type:   reflect.TypeOf(map[string]interface{}{}),
		Object: map[string]interface{}{"a": "x", "b": "y"},
	},
	TestCase{
		String: "a: x\nb: \"y\"",
		Format: "yaml",
		Type:   reflect.TypeOf(map[string]interface{}{}),
		Object: map[string]interface{}{"a": "x", "b": "y"},
	},
}

var serializeTestCases = append(testCases, []TestCase{
	TestCase{
		String: "{\"a\":\"x\"}",
		Format: "jsonl",
		Limit:  1,
		Type:   reflect.TypeOf([]map[string]interface{}{}),
		Object: []map[string]interface{}{map[string]interface{}{"a": "x"}, map[string]interface{}{"ab": "y"}},
	},
}...)

var deserializeTestCases = append(testCases, []TestCase{
	TestCase{
		String: "{\"a\":\"x\"}\n{\"b\":\"y\"}",
		Format: "jsonl",
		Limit:  1,
		Type:   reflect.TypeOf([]map[string]interface{}{}),
		Object: []map[string]interface{}{map[string]interface{}{"a": "x"}},
	},
}...)
