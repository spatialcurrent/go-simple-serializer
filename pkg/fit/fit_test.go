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

func TestFitSliceString(t *testing.T) {
	out := Fit([]string{"a", "b", "c"})
	assert.Equal(t, []string{"a", "b", "c"}, out)
}

func TestFitSliceFloat64(t *testing.T) {
	out := Fit([]interface{}{0.1, 0.2, 0.3})
	assert.Equal(t, []float64{0.1, 0.2, 0.3}, out)
}

func TestFitMapStringString(t *testing.T) {
	out := Fit(map[string]interface{}{"a": "x", "b": "y", "c": "z"})
	assert.Equal(t, map[string]string{"a": "x", "b": "y", "c": "z"}, out)
}

func TestFitMapStringInterfaceMapStringString(t *testing.T) {
	out := Fit(map[string]interface{}{"@type": "null"})
	assert.Equal(t, map[string]string{"@type": "null"}, out)
}

func TestFitMapStringInterface(t *testing.T) {
	out := Fit(map[interface{}]interface{}{"a": true, "b": 1, "c": "z"})
	assert.Equal(t, map[string]interface{}{"a": true, "b": 1, "c": "z"}, out)
}

func TestFitMapInterfaceString(t *testing.T) {
	out := Fit(map[interface{}]interface{}{true: "x", 4: "y", "c": "z"})
	assert.Equal(t, map[interface{}]string{true: "x", 4: "y", "c": "z"}, out)
}
