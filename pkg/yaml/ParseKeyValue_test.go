// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package yaml

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseKeyValue(t *testing.T) {
	key, value, ok := ParseKeyValue([]byte("a: b\n"))
	assert.Equal(t, "a", string(key))
	assert.Equal(t, "b", string(value))
	assert.True(t, ok)
}
