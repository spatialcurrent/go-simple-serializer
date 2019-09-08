// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gss

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCanStreamCSVJSONL(t *testing.T) {
	assert.True(t, CanStream("csv", "jsonl"))
}
