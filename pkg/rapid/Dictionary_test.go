// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package rapid

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDictionary(t *testing.T) {
	d := NewDictionary()
	indicies := d.AddChain("hello", "world")
	require.Equal(t, []int{0, 1}, indicies)
	chain, ok := d.GetChain(indicies)
	require.True(t, ok)
	require.Equal(t, []interface{}{"hello", "world"}, chain)
	indicies = d.AddChain("hello", "beautiful", "world")
	require.Equal(t, []int{0, 2, 1}, indicies)
	chain, ok = d.GetChain(indicies)
	require.True(t, ok)
	require.Equal(t, []interface{}{"hello", "beautiful", "world"}, chain)
	indicies = d.AddChain("hello", "beautiful", "world", 1)
	require.Equal(t, []int{0, 2, 1, 3}, indicies)
	chain, ok = d.GetChain(indicies)
	require.True(t, ok)
	require.Equal(t, []interface{}{"hello", "beautiful", "world", 1}, chain)
	indicies = d.AddChain("hello", 1, "beautiful", "world", 1)
	require.Equal(t, []int{0, 3, 2, 1, 3}, indicies)
	chain, ok = d.GetChain(indicies)
	require.True(t, ok)
	require.Equal(t, []interface{}{"hello", 1, "beautiful", "world", 1}, chain)
}
