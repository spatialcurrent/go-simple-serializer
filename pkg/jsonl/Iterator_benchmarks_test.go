// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package jsonl

import (
	"bufio"
	"encoding/json"
	"os"
	"testing"

	"github.com/pkg/errors"
)

// Following guidance from "How to write benchmarks in Go"
//	- https://dave.cheney.net/2013/06/30/how-to-write-benchmarks-in-go

var (
	results = make([]interface{}, 0)
)

func benchmarkIterator(limit int, b *testing.B) {
	path := os.Getenv("GSS_BENCHMARK_JSONL_INPUT_FILE")
	if len(path) == 0 {
		panic(errors.New("missing benchmark input file: GSS_BENCHMARK_JSONL_INPUT_FILE is not set"))
	}
	input, err := os.Open(path)
	if err != nil {
		panic(errors.Wrap(err, "error opening benchmark input file"))
	}
	defer input.Close()
	it := NewIterator(&NewIteratorInput{
		Reader:        bufio.NewReader(input),
		SkipLines:     0,
		Comment:       "",
		Trim:          true,
		SkipBlanks:    false,
		SkipComments:  false,
		LineSeparator: []byte("\n")[0],
		DropCR:        true,
	})
	out := make([]interface{}, 0)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		out = make([]interface{}, 0)
		for i := 0; i < limit; i++ {
			obj, err := it.Next()
			if err != nil {
				panic(err)
			}
			out = append(out, obj)
		}
	}
	results = out
}

func benchmarkStandardLibaryDecoder(limit int, b *testing.B) {
	path := os.Getenv("GSS_BENCHMARK_JSONL_INPUT_FILE")
	if len(path) == 0 {
		panic(errors.New("missing benchmark input file: GSS_BENCHMARK_JSONL_INPUT_FILE is not set"))
	}
	input, err := os.Open(path)
	if err != nil {
		panic(errors.Wrap(err, "error opening benchmark input file"))
	}
	defer input.Close()

	decoder := json.NewDecoder(bufio.NewReader(input))

	out := make([]interface{}, 0)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		out = make([]interface{}, 0)
		for i := 0; i < limit; i++ {
			obj := map[string]interface{}{}
			err := decoder.Decode(&obj)
			if err != nil {
				panic(err)
			}
			out = append(out, obj)
		}
	}
	// defeat compiler optimizations
	results = out
}

func BenchmarkIterator1(b *testing.B) {
	benchmarkIterator(1, b)
}

func BenchmarkIterator1024(b *testing.B) {
	benchmarkIterator(1024, b)
}

func BenchmarkIterator4096(b *testing.B) {
	benchmarkIterator(4096, b)
}

func BenchmarkStandardLibraryDecoder1(b *testing.B) {
	benchmarkStandardLibaryDecoder(1, b)
}

func BenchmarkStandardLibraryDecoder1024(b *testing.B) {
	benchmarkStandardLibaryDecoder(1024, b)
}

func BenchmarkStandardLibraryDecoder4096(b *testing.B) {
	benchmarkStandardLibaryDecoder(4096, b)
}
