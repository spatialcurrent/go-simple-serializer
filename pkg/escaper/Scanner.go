// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package escaper

import (
	"bufio"
	"bytes"
	"io"

	"github.com/spatialcurrent/go-simple-serializer/pkg/splitter"
)

// Scanner provides an easy API to scan over escaped lines.
type Scanner struct {
	*bufio.Scanner
	prefix []byte
	subs   [][]byte
	err    error
}

func NewScanner(r io.Reader, separator byte, dropCR bool, prefix []byte, subs [][]byte) *Scanner {
	scanner := bufio.NewScanner(r)
	scanner.Split(splitter.ScanLines(separator, dropCR))
	return &Scanner{
		Scanner: scanner,
		prefix:  prefix,
		subs:    subs,
		err:     nil,
	}
}

func (s *Scanner) unescape(in []byte) []byte {
	out := in
	if len(s.prefix) > 0 {
		for _, sub := range s.subs {
			out = bytes.Replace(out, append(s.prefix, sub...), sub, -1)
		}
		out = bytes.Replace(out, append(s.prefix, s.prefix...), s.prefix, -1) // unescape the prefix itself
	}
	return out
}

func (s *Scanner) Split(split bufio.SplitFunc) {
	s.Scanner.Split(split)
}

func (s *Scanner) Bytes() []byte {
	return s.unescape(s.Scanner.Bytes())
}

func (s *Scanner) Text() string {
	return string(s.unescape(s.Scanner.Bytes()))
}
