// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package rapid

import (
	"bytes"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func compare(t *testing.T, a string, b string) bool {
	aLines := strings.Split(a, "\n")
	sort.Strings(aLines)
	bLines := strings.Split(b, "\n")
	sort.Strings(bLines)
	require.Equal(t, aLines, bLines)
	return false
}

func TestEncoder(t *testing.T) {
	buf := new(bytes.Buffer)
	encoder := NewEncoder(buf, '\n')
	encoder.Encode(map[string]string{"hello": "world"})
	require.Equal(t, "K:S:hello\nK:S:world\nP:0:1\n---\n", buf.String())
	encoder.Encode(map[string]string{
		"hello": "world",
		"foo":   "bar",
	})
	t.Log("\n" + buf.String())
	//compare(t, "K:S:hello\nK:S:world\nP:0:1\n---\nK:S:hello\nK:S:world\nP:0:1\nK:S:foo\nK:S:bar\nP:2:3\n---\n", buf.String())
}

func TestEncodeGeoJSON(t *testing.T) {
	buf := new(bytes.Buffer)
	encoder := NewEncoder(buf, '\n')
	obj := map[string]interface{}{
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
	encoder.Encode(obj)
	t.Log("\n" + buf.String())
	//compare(t, "K:S:geometry\nK:S:type\nP:0:1:0\nK:S:Feature\nP:1:2\nK:S:properties\nK:S:name\nK:S:Hello World\nP:3:4:5\n---\n", buf.String())
}
