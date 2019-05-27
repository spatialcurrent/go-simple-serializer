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

func TestGetUnknownKeys(t *testing.T) {
	in := map[string]interface{}{
		"a": "x",
		"b": "y",
		"c": "z",
	}
	keys := GetUnknownKeys(in, map[interface{}]struct{}{"b": struct{}{}}, true)
	assert.Equal(t, []interface{}{"a", "c"}, keys)
}
