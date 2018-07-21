// =================================================================
//
// Copyright (C) 2018 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

import (
	"github.com/pkg/errors"
)

import (
	"github.com/spatialcurrent/go-simple-serializer/gss"
)

var GO_GSS_VERSION = "0.0.1"

func printUsage() {
	fmt.Println("Usage: gss -i INPUT_FORMAT -o OUTPUT_FORMAT")
}

func main() {

	var input_format string
	var output_format string

	var version bool
	var help bool

	flag.StringVar(&input_format, "i", "", "The input format: csv, hcl, hcl2, json, jsonl, toml, yaml")
	flag.StringVar(&output_format, "o", "", "The output format: csv, hcl, hcl2, json, jsonl, toml, yaml")
	flag.BoolVar(&version, "version", false, "Prints version to stdout.")
	flag.BoolVar(&help, "help", false, "Print help.")

	flag.Parse()

	if help {
		printUsage()
		fmt.Println("Options:")
		flag.PrintDefaults()
		os.Exit(0)
	} else if len(os.Args) == 1 {
		fmt.Println("Error: Provided no arguments.")
		fmt.Println("Run \"gss -help\" for more information.")
		os.Exit(0)
	} else if len(os.Args) == 2 && os.Args[1] == "help" {
		printUsage()
		fmt.Println("Options:")
		flag.PrintDefaults()
		os.Exit(0)
	}

	if version {
		fmt.Println(GO_GSS_VERSION)
		os.Exit(0)
	}

	if len(input_format) == 0 {
		fmt.Println("Error: Provided no -input_format.")
		fmt.Println("Run \"gss -help\" for more information.")
		os.Exit(1)
	}

	if len(output_format) == 0 {
		fmt.Println("Error: Provided no -output_format.")
		fmt.Println("Run \"gss -help\" for more information.")
		os.Exit(1)
	}

	input_bytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Println(errors.Wrap(err, "Error reading from stdin"))
		os.Exit(1)
	}

	output_string, err := gss.Convert(string(input_bytes), input_format, output_format)
	if err != nil {
		fmt.Println(errors.Wrap(err, "Error converting"))
		os.Exit(1)
	}
	fmt.Println(output_string)
}
