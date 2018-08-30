// =================================================================
//
// Copyright (C) 2018 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

// gss is the command line program for go-simple-serializer (GSS).
//
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

import (
	"github.com/pkg/errors"
)

import (
	"github.com/spatialcurrent/go-simple-serializer/gss"
)

func printUsage() {
	fmt.Println("Usage: gss -i INPUT_FORMAT -o OUTPUT_FORMAT [-h HEADER] [-c COMMENT]")
}

func main() {

	var input_format string
	var input_header_text string
	var input_comment string

	var output_format string

	var version bool
	var verbose bool
	var help bool

	flag.StringVar(&input_format, "i", "", "The input format: "+strings.Join(gss.Formats, ", "))
	flag.StringVar(&input_header_text, "h", "", "The input header if the stdin input has no header.")
	flag.StringVar(&input_comment, "c", "", "The input comment character, e.g., #.  Commented lines are not sent to output.")
	flag.StringVar(&output_format, "o", "", "The output format: "+strings.Join(gss.Formats, ", "))
	flag.BoolVar(&version, "version", false, "Prints version to stdout")
	flag.BoolVar(&verbose, "verbose", false, "Print debug info to stdout")
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
		fmt.Println(gss.VERSION)
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

	input_header := make([]string, 0)
	if len(input_header_text) > 0 {
		input_header = strings.Split(input_header_text, ",")
	}

	output_string, err := gss.Convert(string(input_bytes), input_format, input_header, input_comment, output_format, verbose)
	if err != nil {
		fmt.Println(errors.Wrap(err, "Error converting"))
		os.Exit(1)
	}
	fmt.Println(output_string)
}
