// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package escaper

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScanner(t *testing.T) {
	in := `{""properties"":{""shelter"":""yes"",""tourism"":""picnic_site""},""type"":""Feature""}`
	scanner := NewScanner(
		strings.NewReader(in),
		'\n',
		true,
		[]byte("\""),
		[][]byte{
			[]byte("\""),
		},
	)
	valid := scanner.Scan()
	assert.True(t, valid)
	assert.NoError(t, scanner.Err())
	assert.Equal(t, `{"properties":{"shelter":"yes","tourism":"picnic_site"},"type":"Feature"}`, scanner.Text())
	valid = scanner.Scan()
	assert.False(t, valid)
	assert.NoError(t, scanner.Err())
}
