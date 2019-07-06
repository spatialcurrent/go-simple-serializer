// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package inspector

import (
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

func TestGetKeys(t *testing.T) {
	in := map[string]interface{}{
		"a": "x",
		"b": "y",
		"c": "z",
	}
	keys := GetKeys(in, true, false)
	assert.Equal(t, []interface{}{"a", "b", "c"}, keys)
}

func TestGetKeysReversed(t *testing.T) {
	in := map[string]interface{}{
		"a": "x",
		"b": "y",
		"c": "z",
	}
	keys := GetKeys(in, true, true)
	assert.Equal(t, []interface{}{"c", "b", "a"}, keys)
}
