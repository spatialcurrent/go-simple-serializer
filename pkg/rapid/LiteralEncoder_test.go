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

func TestLiteralEncoderFalse(t *testing.T) {
	buf := new(bytes.Buffer)
	encoder := NewLiteralEncoder(buf)
	encoder.Encode(false)
	require.Equal(t, []byte{DefaultTypeBool, 0b0}, buf.Bytes())
}

func TestLiteralEncoderTrue(t *testing.T) {
	buf := new(bytes.Buffer)
	encoder := NewLiteralEncoder(buf)
	encoder.Encode(true)
	require.Equal(t, []byte{DefaultTypeBool, 0b1}, buf.Bytes())
}

func TestLiteralEncoderZero(t *testing.T) {
	buf := new(bytes.Buffer)
	encoder := NewLiteralEncoder(buf)
	encoder.Encode(0)
	require.Equal(t, []byte{DefaultTypeInt8, 0b0}, buf.Bytes())
}

func TestLiteralEncoderOne(t *testing.T) {
	buf := new(bytes.Buffer)
	encoder := NewLiteralEncoder(buf)
	encoder.Encode(1)
	require.Equal(t, []byte{DefaultTypeInt8, 0b1}, buf.Bytes())
}

func TestLiteralEncoder256(t *testing.T) {
	buf := new(bytes.Buffer)
	encoder := NewLiteralEncoder(buf)
	encoder.Encode(256)
	require.Equal(t, []byte{DefaultTypeInt64, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}, buf.Bytes())
}

func TestLiteralEncoderFloat64(t *testing.T) {
	buf := new(bytes.Buffer)
	encoder := NewLiteralEncoder(buf)
	encoder.Encode(10.5)
	require.Equal(t, []byte{DefaultTypeFloat64, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x25, 0x40}, buf.Bytes())
}

func TestLiteralEncoderString(t *testing.T) {
	buf := new(bytes.Buffer)
	encoder := NewLiteralEncoder(buf)
	encoder.Encode("hello world")
	require.Equal(t, append(append([]byte{DefaultTypeString}, []byte("hello world")...), byte(0)), buf.Bytes())
}

func TestLiteralEncoderSliceString(t *testing.T) {
	expected := flatten([][]byte{
		[]byte{DefaultTypeArray & DefaultTypeString, DefaultTypeInt8, 2},
		[]byte("hello"),
		[]byte{0},
		[]byte("world"),
		[]byte{0},
	})
	buf := new(bytes.Buffer)
	encoder := NewLiteralEncoder(buf)
	encoder.Encode([]string{"hello", "world"})
	require.Equal(
		t,
		expected,
		buf.Bytes(),
	)
}

func TestLiteralEncoderSliceInt64(t *testing.T) {
	expected := []byte{
		DefaultTypeArray & DefaultTypeInt64,
		DefaultTypeInt8,
		2,
	}
	buf := new(bytes.Buffer)
	encoder := NewLiteralEncoder(buf)
	encoder.Encode([]int{256, 512})
	require.Equal(
		t,
		expected,
		buf.Bytes(),
	)
}

func TestLiteralEncoderMapStringString(t *testing.T) {
	expected := flatten([][]byte{
		[]byte{DefaultTypeMap & DefaultTypeString, DefaultTypeString, DefaultTypeInt8, 1},
		[]byte("hello"),
		[]byte{0},
		[]byte("world"),
		[]byte{0},
	})
	buf := new(bytes.Buffer)
	encoder := NewLiteralEncoder(buf)
	encoder.Encode(map[string]string{"hello": "world"})
	require.Equal(t, expected, buf.Bytes())
}

func TestLiteralEncoderMapStringInterface(t *testing.T) {
	expected := flatten([][]byte{
		[]byte{DefaultTypeMap & DefaultTypeString, DefaultTypeInterface, DefaultTypeInt8, 1},
		[]byte{DefaultTypeString},
		[]byte("hello"),
		[]byte{0},
		[]byte{DefaultTypeString},
		[]byte("world"),
		[]byte{0},
	})
	buf := new(bytes.Buffer)
	encoder := NewLiteralEncoder(buf)
	encoder.Encode(map[string]interface{}{"hello": "world"})
	require.Equal(t, expected, buf.Bytes())
}
