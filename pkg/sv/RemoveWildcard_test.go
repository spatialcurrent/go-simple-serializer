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

func TestRemoveWildcard(t *testing.T) {
	in := []interface{}{
		"a",
		"*",
		"c",
	}
	out := RemoveWildcard(in)
	assert.Equal(t, []interface{}{"a", "c"}, out)
}
