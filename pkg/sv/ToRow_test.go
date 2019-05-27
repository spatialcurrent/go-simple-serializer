// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package sv

import (
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

import (
	"github.com/spatialcurrent/go-stringify/pkg/stringify"
)

func TestToRow(t *testing.T) {
	in := map[string]interface{}{
		"a": "x",
		"b": "y",
		"c": "z",
	}
	header := []interface{}{"a", "b", "c"}
	row, err := ToRow(in, header, stringify.DefaultValueStringer(""))
	assert.NoError(t, err)
	assert.Equal(t, []string{"x", "y", "z"}, row)
}
