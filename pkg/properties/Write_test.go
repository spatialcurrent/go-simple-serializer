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

	"github.com/stretchr/testify/assert"

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
		KeySerializer:     stringify.NewStringer("", false, false, false),
		ValueSerializer:   stringify.NewStringer("", false, false, false),
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

func TestWriteDecimal(t *testing.T) {
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
		ValueSerializer:   stringify.NewStringer("", true, false, false),
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
		ValueSerializer:   stringify.NewStringer("", false, false, false),
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

func TestWriteLower(t *testing.T) {
	in := struct {
		A string
		B int
		C bool
	}{
		A: "foo",
		B: 1,
		C: true,
	}
	buf := new(bytes.Buffer)
	err := Write(&WriteInput{
		Writer:            buf,
		KeyValueSeparator: "=",
		LineSeparator:     "\n",
		Object:            in,
		KeySerializer:     stringify.NewStringer("", false, true, false),
		ValueSerializer:   stringify.NewStringer("", false, false, false),
		Sorted:            true,
		EscapePrefix:      "\\",
		EscapeSpace:       false,
		EscapeEqual:       true,
		EscapeColon:       false,
		EscapeNewLine:     false,
	})
	assert.NoError(t, err)
	out := buf.String()
	assert.Equal(t, "a=foo\nb=1\nc=true", out)
}
