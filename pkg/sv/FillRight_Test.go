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

func TestFillRight(t *testing.T) {
	in := []string{"a", "b", "c"}
	out := FillRight(in, 6)
	assert.Equal(t, []string{"a", "b", "c", "", "", ""}, out)
}
