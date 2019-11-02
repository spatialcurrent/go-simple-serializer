// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package mapper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsEmptyString(t *testing.T) {
	assert.True(t, IsEmpty(""))
}
