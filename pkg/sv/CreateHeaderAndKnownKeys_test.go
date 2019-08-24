// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package sv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateHeaderAndKnownKeys(t *testing.T) {
	in := map[string]interface{}{
		"a": "x",
		"b": "y",
		"c": "z",
	}
	expectedKnownKeys := map[interface{}]struct{}{
		"a": struct{}{},
		"b": struct{}{},
		"c": struct{}{},
	}
	header, knownKeys := CreateHeaderAndKnownKeys(in, true, false)
	assert.Equal(t, []interface{}{"a", "b", "c"}, header)
	assert.Equal(t, expectedKnownKeys, knownKeys)
}

func TestCreateHeaderAndKnownKeysPointer(t *testing.T) {
	in := map[string]interface{}{
		"a": "x",
		"b": "y",
		"c": "z",
	}
	expectedKnownKeys := map[interface{}]struct{}{
		"a": struct{}{},
		"b": struct{}{},
		"c": struct{}{},
	}
	header, knownKeys := CreateHeaderAndKnownKeys(&in, true, false)
	assert.Equal(t, []interface{}{"a", "b", "c"}, header)
	assert.Equal(t, expectedKnownKeys, knownKeys)
}
