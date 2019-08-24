// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package sv

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpandHeaderWithWildcard(t *testing.T) {
	in := map[string]interface{}{
		"a": "x",
		"b": "y",
		"c": "z",
	}
	header := []interface{}{"b", Wildcard}
	knownKeys := map[interface{}]struct{}{
		"b": struct{}{},
	}
	expectedKnownKeys := map[interface{}]struct{}{
		"a": struct{}{},
		"b": struct{}{},
		"c": struct{}{},
	}

	newHeader, newKnownKeys := ExpandHeaderWithWildcard(header, knownKeys, reflect.ValueOf(in), true, false)
	assert.Equal(t, []interface{}{"b", "a", "c", Wildcard}, newHeader)
	assert.Equal(t, expectedKnownKeys, newKnownKeys)

	in = map[string]interface{}{
		"d": "x",
		"e": "y",
		"f": "z",
	}
	header = newHeader
	knownKeys = newKnownKeys
	expectedKnownKeys = map[interface{}]struct{}{
		"a": struct{}{},
		"b": struct{}{},
		"c": struct{}{},
		"d": struct{}{},
		"e": struct{}{},
		"f": struct{}{},
	}
	newHeader, newKnownKeys = ExpandHeaderWithWildcard(header, knownKeys, reflect.ValueOf(in), true, false)
	assert.Equal(t, []interface{}{"b", "a", "c", "d", "e", "f", Wildcard}, newHeader)
	assert.Equal(t, expectedKnownKeys, newKnownKeys)
}
