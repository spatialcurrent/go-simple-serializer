// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package inspector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFieldNames(t *testing.T) {
	in := struct {
		A string
		B string
		C string
	}{
		A: "x",
		B: "y",
		C: "z",
	}
	fieldNames := GetFieldNames(in, true, false)
	assert.Equal(t, []string{"A", "B", "C"}, fieldNames)
}

func TestGetFieldNamesReversed(t *testing.T) {
	in := struct {
		A string
		B string
		C string
	}{
		A: "x",
		B: "y",
		C: "z",
	}
	fieldNames := GetFieldNames(in, true, true)
	assert.Equal(t, []string{"C", "B", "A"}, fieldNames)
}
