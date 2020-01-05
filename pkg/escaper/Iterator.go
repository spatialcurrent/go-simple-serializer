// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package escaper

import (
	"bytes"
	"io"
)

// Iterator provides an easy API to iterate over a escaped string.
type Iterator struct {
	reader io.ByteReader
	prefix []byte
	subs   [][]byte
	delim  []byte
	err    error
}

// NewIterator returns a new iterator the reads escaped bytes from the given ByteReader.
//  The iterator uses the given escape prefix, subsitution bytes and delimiter.
func NewIterator(r io.ByteReader, prefix []byte, subs [][]byte, delim []byte) *Iterator {
	return &Iterator{
		reader: r,
		prefix: prefix,
		subs:   subs,
		delim:  delim,
		err:    nil,
	}
}

func (e *Iterator) unescape(in []byte) []byte {
	out := in
	if len(e.prefix) > 0 {
		for _, sub := range e.subs {
			out = bytes.Replace(out, append(e.prefix, sub...), sub, -1)
		}
		out = bytes.Replace(out, append(e.prefix, e.prefix...), e.prefix, -1) // unescape the prefix itself
	}
	return out
}

// Next returns the next token until the end of the input reader is reached.
// Once the input reader is exhausted, Next returns io.EOF.
func (it *Iterator) Next() ([]byte, error) {
	buf := make([]byte, 0)
	if it.err != nil {
		return buf, it.err
	}
	for {
		// read the next byte
		b, err := it.reader.ReadByte()
		if err != nil {
			it.err = err
			if err == io.EOF {
				break
			}
			return buf, err
		}
		// append byte to buffer
		buf = append(buf, b)
		// if buffer ends with delim
		if bytes.HasSuffix(buf, it.delim) {
			if !bytes.HasSuffix(buf[0:len(buf)-len(it.delim)], it.prefix) {
				return it.unescape(buf[0 : len(buf)-1]), nil
			}
		}
	}
	return it.unescape(buf), nil
}

// Reset sets the iterator's reader to the given reader and sets the stored error to nil.
func (it *Iterator) Reset(r io.ByteReader) {
	it.reader = r
	it.err = nil
}
