// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package rapid

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecoder(t *testing.T) {
	str := "K:S:hello\nK:S:world\nP:0:1\n---\nK:S:hello\nK:S:world\nP:0:1\nK:S:foo\nK:S:bar\nP:2:3\n---\n"
	decoder := NewDecoder(strings.NewReader(str), '\n', true)
	// decode first object
	m := map[string]interface{}{}
	err := decoder.Decode(&m)
	require.NoError(t, err)
	require.Equal(t, map[string]interface{}{"hello": "world"}, m)
	// decode second object
	n := map[string]interface{}{}
	err = decoder.Decode(&n)
	require.NoError(t, err)
	require.Equal(t, map[string]interface{}{"hello": "world", "foo": "bar"}, n)
}

func TestDecodeGeoJSON(t *testing.T) {
	str := `K:S:type
K:S:Feature
P:0:1
K:S:id
K:I:1
P:2:3
K:S:properties
P:4:2:3
K:S:name
K:S:Hello World
P:4:5:6
K:S:geometry
K:S:Point
P:7:0:8
K:S:coordinates
K:A:2
K:I:0
K:F:9.055337
P:7:9:10:11:12
K:I:1
K:F:50.005633
P:7:9:10:13:14`
	expected := map[string]interface{}{
		"type": "Feature",
		"id":   1,
		"properties": map[string]interface{}{
			"id":   1,
			"name": "Hello World",
		},
		"geometry": map[string]interface{}{
			"type":        "Point",
			"coordinates": []interface{}{9.0553371, 50.0056328},
		},
	}
	decoder := NewDecoder(strings.NewReader(str), '\n', true)
	// decode first object
	m := map[string]interface{}{}
	err := decoder.Decode(&m)
	require.NoError(t, err)
	require.Equal(t, expected, m)
}
