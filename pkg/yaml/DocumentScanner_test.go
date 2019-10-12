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
	in := []byte("hello\n---\nbeautiful\n---\nworld\n---\n")
	d := NewDocumentScanner(bytes.NewReader(in), true)

	require.True(t, d.Scan())
	require.Equal(t, "hello", d.Text())

	require.True(t, d.Scan())
	require.Equal(t, "beautiful", d.Text())

	require.True(t, d.Scan())
	require.Equal(t, "world", d.Text())

	require.False(t, d.Scan())
}
