// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package properties

import (
	"bytes"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

import (
	"github.com/spatialcurrent/go-stringify/pkg/stringify"
)

func TestWrite(t *testing.T) {
	in := map[interface{}]interface{}{
		"foo": "bar",
		"a":   1,
		"b":   1234567890.123,
		"c":   nil,
		1:     "hello world",
	}

	buf := new(bytes.Buffer)
	err := Write(&WriteInput{
		Writer:            buf,
		KeyValueSeparator: "=",
		LineSeparator:     "\n",
		Object:            in,
		ValueSerializer:   stringify.DefaultValueStringer(""),
		Sorted:            true,
		EscapePrefix:      "\\",
		EscapeSpace:       true,
		EscapeEqual:       true,
		EscapeColon:       false,
		EscapeNewLine:     false,
	})
	assert.NoError(t, err)
	out := buf.String()
	assert.Equal(t, "1=hello\\ world\na=1\nb=1.234567890123e+09\nc=\nfoo=bar", out)
}

func TestWriteDecimalValueStringer(t *testing.T) {
	in := map[interface{}]interface{}{
		"foo": "bar",
		"a":   1,
		"b":   1234567890.123,
		"c":   nil,
		1:     "10",
	}

	buf := new(bytes.Buffer)
	err := Write(&WriteInput{
		Writer:            buf,
		KeyValueSeparator: "=",
		LineSeparator:     "\n",
		Object:            in,
		ValueSerializer:   stringify.DecimalValueStringer(""),
		Sorted:            true,
		EscapePrefix:      "\\",
		EscapeSpace:       false,
		EscapeEqual:       true,
		EscapeColon:       false,
		EscapeNewLine:     false,
	})
	assert.NoError(t, err)
	out := buf.String()
	assert.Equal(t, "1=10\na=1\nb=1234567890.123000\nc=\nfoo=bar", out)
}

func TestWriteStruct(t *testing.T) {
	in := struct {
		A string
		B int
	}{
		A: "foo",
		B: 1,
	}
	buf := new(bytes.Buffer)
	err := Write(&WriteInput{
		Writer:            buf,
		KeyValueSeparator: "=",
		LineSeparator:     "\n",
		Object:            in,
		ValueSerializer:   stringify.DefaultValueStringer(""),
		Sorted:            true,
		EscapePrefix:      "\\",
		EscapeSpace:       false,
		EscapeEqual:       true,
		EscapeColon:       false,
		EscapeNewLine:     false,
	})
	assert.NoError(t, err)
	out := buf.String()
	assert.Equal(t, "A=foo\nB=1", out)
}
