// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package jsonl

import (
	"bufio"
	//"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/pkg/errors"
	"github.com/spatialcurrent/go-pipe/pkg/pipe"
	"github.com/spatialcurrent/go-stringify/pkg/stringify"
)

// Following guidance from "How to write benchmarks in Go"
//	- https://dave.cheney.net/2013/06/30/how-to-write-benchmarks-in-go

func loadBenchmarkData(path string) []interface{} {
	input, err := os.Open(path)
	if err != nil {
		panic(errors.Wrap(err, "error opening benchmark input file"))
	}
	defer input.Close()
	data := make([]interface{}, 0)
	err = pipe.NewBuilder().
		Input(NewIterator(&NewIteratorInput{
			Reader:        bufio.NewReader(input),
			SkipLines:     0,
			Comment:       "",
			Trim:          true,
			SkipBlanks:    false,
			SkipComments:  false,
			LineSeparator: []byte("\n")[0],
			DropCR:        true,
		})).
		OutputF(func(x interface{}) error {
			data = append(data, x)
			return nil
		}).
		Run()
	if err != nil {
		panic(errors.Wrap(err, "error reading in data"))
	}
	return data
}

func benchmarkWriter(batch int, b *testing.B) {
	inputPath := os.Getenv("GSS_BENCHMARK_JSONL_INPUT_FILE")
	if len(inputPath) == 0 {
		panic(errors.New("missing benchmark input file: GSS_BENCHMARK_JSONL_INPUT_FILE is not set"))
	}
	outputDir := os.Getenv("GSS_BENCHMARK_JSONL_OUTPUT_DIR")
	if len(outputDir) == 0 {
		panic(errors.New("missing benchmark input file: GSS_BENCHMARK_JSONL_OUTPUT_DIR is not set"))
	}
	//
	data := loadBenchmarkData(inputPath)
	//
	output, err := ioutil.TempFile(outputDir, fmt.Sprintf("benchmark_output_%d_%d_*.jsonl", batch, b.N))
	if err != nil {
		panic(errors.Wrap(err, "error creating temporary output file for running benchmarks"))
	}
	if os.Getenv("GSS_BENCHMARK_JSONL_OUTPUT_KEEP") != "1" {
		defer os.Remove(output.Name())
	}
	//
	writer := NewWriter(output, "\n", stringify.NewDecimalStringer(), false)
	//
	b.ResetTimer()
	//
	for n := 0; n < b.N; n++ {
		cursor := 0
		if batch == 1 {
			err := writer.WriteObject(data[cursor])
			if err != nil {
				panic(err)
			}
		} else {
			err := writer.WriteObjects(data[cursor : cursor+batch])
			if err != nil {
				panic(err)
			}
		}
		cursor += batch
	}
}

func BenchmarkWriter1(b *testing.B) {
	benchmarkWriter(1, b)
}

func BenchmarkWriter64(b *testing.B) {
	benchmarkWriter(64, b)
}

func BenchmarkWriter256(b *testing.B) {
	benchmarkWriter(256, b)
}

func BenchmarkWriter1024(b *testing.B) {
	benchmarkWriter(1024, b)
}

func BenchmarkWriter4096(b *testing.B) {
	benchmarkWriter(4096, b)
}
