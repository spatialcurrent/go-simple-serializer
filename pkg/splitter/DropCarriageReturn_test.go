// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package splitter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDropCarriageReturnWithoutCR(t *testing.T) {
	in := []byte("hello")
	out := DropCarriageReturn(in)
	assert.Equal(t, []byte("hello"), out)
}

func TestDropCarriageReturnWithCR(t *testing.T) {
	in := []byte("hello\r")
	out := DropCarriageReturn(in)
	assert.Equal(t, []byte("hello"), out)
}
