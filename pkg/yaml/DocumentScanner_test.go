// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package yaml

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDocumentScanner(t *testing.T) {
	in := []byte("hello\nperson\n---\nbeautiful\n---\nworld\n---\n")
	d := NewDocumentScanner(bytes.NewReader(in), true)

	require.True(t, d.Scan())
	require.Equal(t, "hello\nperson\n", d.Text())

	require.True(t, d.Scan())
	require.Equal(t, "beautiful\n", d.Text())

	require.True(t, d.Scan())
	require.Equal(t, "world\n", d.Text())

	require.False(t, d.Scan())
}
