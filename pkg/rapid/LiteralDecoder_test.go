// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package rapid

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLiteralDecoderBool(t *testing.T) {
	decoder := NewLiteralDecoder(bytes.NewReader([]byte{
		DefaultTypeBool, 0x1,
		DefaultTypeBool, 0x0,
		DefaultTypeInt8, 0x4,
	}))
	v := false
	err := decoder.Decode(&v)
	require.NoError(t, err)
	require.Equal(t, true, v)
	err = decoder.Decode(&v)
	require.NoError(t, err)
	require.Equal(t, false, v)
}

func TestLiteralDecoderInt(t *testing.T) {
	decoder := NewLiteralDecoder(bytes.NewReader([]byte{
		DefaultTypeInt8, 0x4,
		DefaultTypeInt64, 0x4, 0x4, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
	}))
	v := 0
	err := decoder.Decode(&v)
	require.NoError(t, err)
	require.Equal(t, 4, v)
	err = decoder.Decode(&v)
	require.NoError(t, err)
	require.Equal(t, 1028, v)
}

func TestLiteralDecoderString(t *testing.T) {
	decoder := NewLiteralDecoder(bytes.NewReader(flatten([][]byte{
		[]byte{DefaultTypeString},
		[]byte("hello world"),
		[]byte{byte(0)},
		[]byte{DefaultTypeString},
		[]byte("hello\\nworld"),
		[]byte{byte(0)},
	})))
	v := ""
	err := decoder.Decode(&v)
	require.NoError(t, err)
	require.Equal(t, "hello world", v)
	err = decoder.Decode(&v)
	require.NoError(t, err)
	require.Equal(t, "hello\nworld", v)
}
