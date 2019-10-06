// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package fit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFitSlice(t *testing.T) {
	out := Fit([]string{"a", "b", "c"})
	assert.Equal(t, []string{"a", "b", "c"}, out)
}

func TestFitMap(t *testing.T) {
	out := Fit(map[string]interface{}{"a": "x", "b": "y", "c": "z"})
	assert.Equal(t, map[string]string{"a": "x", "b": "y", "c": "z"}, out)
}

func TestFitSlice(t *testing.T) {
	out := Fit([]interface{}{0.1, 0.2, 0.3})
	assert.Equal(t, []float64{0.1, 0.2, 0.3}, out)
}
