// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package toml

import (
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

func TestMarshalNil(t *testing.T) {
	b, err := Marshal(nil)
	assert.Equal(t, err, ErrNilObject)
	assert.Equal(t, string(b), "")
}

func TestMarshalMap(t *testing.T) {
	b, err := Marshal(map[string]interface{}{"a": 1.0, "b": 2.0, "c": 3.0})
	assert.NoError(t, err)
	assert.Equal(t, "a = 1.0\nb = 2.0\nc = 3.0\n", string(b))
}
