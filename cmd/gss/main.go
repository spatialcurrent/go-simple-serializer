// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

// gss is the command line program for go-simple-serializer (GSS).
//
package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

import (
	"github.com/spatialcurrent/go-pipe/pkg/pipe"
	"github.com/spatialcurrent/go-simple-serializer/gss"
	"github.com/spatialcurrent/go-simple-serializer/pkg/iterator"
	"github.com/spatialcurrent/go-simple-serializer/pkg/jsonl"
	"github.com/spatialcurrent/go-simple-serializer/pkg/sv"
)

var gitTag string
var gitBranch string
var gitCommit string

func buildValueSerializer(decimal bool, noDataValue string) func(object interface{}) (string, error) {
	if decimal {
		return sv.DecimalValueSerializer(noDataValue)
	}
	return sv.DefaultValueSerializer(noDataValue)
}

func buildWriter(outputWriter io.Writer, outputFormat string, outputHeader []interface{}, outputValueSerializer func(object interface{}) (string, error)) (pipe.Writer, error) {
	if outputFormat == "csv" || outputFormat == "tsv" {
		separator, err := sv.FormatToSeparator(outputFormat)
		if err != nil {
			return nil, err
		}
		return sv.NewWriter(outputWriter, separator, outputHeader, outputValueSerializer), nil
	} else if outputFormat == "jsonl" {
		return jsonl.NewWriter(outputWriter), nil
	}
	return nil, fmt.Errorf("invalid format %q", outputFormat)
}

func canStream(inputFormat string, outputFormat string, outputSorted bool) bool {
	if !outputSorted {
		if inputFormat == "csv" || inputFormat == "tsv" {
			if outputFormat == "csv" || outputFormat == "tsv" || outputFormat == "jsonl" {
				return true
			}
		} else if inputFormat == "jsonl" {
			if outputFormat == "jsonl" {
				return true
			}
		}
	}
	return false
}

func main() {
	rootCommand := &cobra.Command{
		Use:   "gss",
		Short: "gss",
		Long:  `gss is a simple program for serializing/deserializing data.`,
		RunE: func(cmd *cobra.Command, args []string) error {

			v := viper.New()
			err := v.BindPFlags(cmd.Flags())
			if err != nil {
				return err
			}
			v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
			v.AutomaticEnv()

			inputFormat := v.GetString("input-format")

			if len(inputFormat) == 0 {
				return errors.New("input-format is required")
			}

			outputFormat := v.GetString("output-format")

			if len(outputFormat) == 0 {
				return errors.New("output-format is required")
			}

			outputHeader := v.GetStringSlice("output-header")

			outputSorted := v.GetBool("output-sorted")

			outputValueSerializer := buildValueSerializer(v.GetBool("output-decimal"), v.GetString("output-no-data-value"))

			if canStream(inputFormat, outputFormat, outputSorted) {

				p := pipe.NewBuilder()

				it, errorIterator := iterator.NewIterator(&iterator.NewIteratorInput{
					Reader:       os.Stdin,
					Format:       inputFormat,
					SkipLines:    v.GetInt("input-skip-lines"),
					SkipBlanks:   true,
					SkipComments: true,
					Comment:      v.GetString("input-comment"),
					Trim:         v.GetBool("input-trim"),
				})
				if it == nil {
					return errors.New(fmt.Sprintf("error building input iterator with format %q", inputFormat))
				}
				if errorIterator != nil {
					return errors.Wrap(errorIterator, "error creating input iterator")
				}
				p = p.Input(it)

				if inputLimit := v.GetInt("input-limit"); inputLimit >= 0 {
					p = p.InputLimit(inputLimit)
				}

				outputColumns := make([]interface{}, 0, len(outputHeader))
				for _, str := range outputHeader {
					outputColumns = append(outputColumns, str)
				}

				w, errorWriter := buildWriter(
					os.Stdout,
					outputFormat,
					outputColumns,
					outputValueSerializer)
				if errorWriter != nil {
					return errors.Wrap(errorWriter, "error building output writer")
				}
				p = p.Output(w)

				if outputLimit := v.GetInt("output-limit"); outputLimit >= 0 {
					p = p.OutputLimit(outputLimit)
				}

				errorRun := p.Run()
				if errorRun != nil {
					return errors.Wrap(errorRun, "error piping data")
				}
				return nil
			}

			inputBytes, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				return errors.Wrap(err, "error reading from stdin")
			}

			outputString, err := gss.Convert(&gss.ConvertInput{
				InputBytes:            inputBytes,
				InputFormat:           inputFormat,
				InputHeader:           v.GetStringSlice("input-header"),
				InputComment:          v.GetString("input-comment"),
				InputLazyQuotes:       v.GetBool("input-lazy-quotes"),
				InputSkipLines:        v.GetInt("input-skip-lines"),
				InputLimit:            v.GetInt("input-limit"),
				OutputFormat:          outputFormat,
				OutputHeader:          v.GetStringSlice("output-header"),
				OutputLimit:           v.GetInt("output-limit"),
				OutputPretty:          v.GetBool("output-pretty"),
				OutputSorted:          v.GetBool("output-sorted"),
				OutputValueSerializer: outputValueSerializer,
				Async:                 v.GetBool("async"),
				Verbose:               v.GetBool("verbose"),
			})
			if err != nil {
				return errors.Wrap(err, "error converting")
			}
			fmt.Println(outputString)
			return nil
		},
	}
	flags := rootCommand.Flags()
	flags.StringP("input-format", "i", "", "The input format: "+strings.Join(gss.Formats, ", "))
	flags.StringSlice("input-header", []string{}, "The input header if the stdin input has no header.")
	flags.StringP("input-comment", "c", "", "The input comment character, e.g., #.  Commented lines are not sent to output.")
	flags.Bool("input-lazy-quotes", false, "allows lazy quotes for CSV and TSV")
	flags.Int("input-skip-lines", gss.NoSkip, "The number of lines to skip before processing")
	flags.IntP("input-limit", "l", gss.NoLimit, "The input limit")
	flags.BoolP("input-trim", "t", false, "trim input lines")
	flags.StringP("output-format", "o", "", "The output format: "+strings.Join(gss.Formats, ", "))
	flags.StringSlice("output-header", []string{}, "The output header if the stdout output has no header.")
	flags.IntP("output-limit", "n", gss.NoLimit, "the output limit")
	flags.BoolP("output-pretty", "p", false, "print pretty output")
	flags.BoolP("output-sorted", "s", false, "sort output")
	flags.BoolP("output-decimal", "d", false, "when converting floats to strings use decimals rather than scientific notation")
	flags.StringP("output-no-data-value", "0", "", "no data value, e.g., used for missing values when converting JSON to CSV")
	flags.BoolP("async", "a", false, "async processing")
	flags.Bool("verbose", false, "Print debug info to stdout")

	completionCommandLong := ""
	if _, err := os.Stat("/etc/bash_completion.d/"); !os.IsNotExist(err) {
		completionCommandLong = "To install completion scripts run:\ngss completion > /etc/bash_completion.d/gss"
	} else {
		if _, err := os.Stat("/usr/local/etc/bash_completion.d/"); !os.IsNotExist(err) {
			completionCommandLong = "To install completion scripts run:\ngss completion > /usr/local/etc/bash_completion.d/gss"
		} else {
			completionCommandLong = "To install completion scripts run:\ngss completion > .../bash_completion.d/gss"
		}
	}

	completionCommand := &cobra.Command{
		Use:   "completion",
		Short: "Generates bash completion scripts",
		Long:  completionCommandLong,
		RunE: func(cmd *cobra.Command, args []string) error {
			return rootCommand.GenBashCompletion(os.Stdout)
		},
	}
	rootCommand.AddCommand(completionCommand)

	version := &cobra.Command{
		Use:   "version",
		Short: "print version information to stdout",
		Long:  "print version information to stdout",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(gitTag) > 0 {
				fmt.Println("Tag: " + gitTag)
			}
			if len(gitBranch) > 0 {
				fmt.Println("Branch: " + gitBranch)
			}
			if len(gitCommit) > 0 {
				fmt.Println("Commit: " + gitCommit)
			}
			return nil
		},
	}

	rootCommand.AddCommand(version)

	if err := rootCommand.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, err.Error()+"\n")
		os.Exit(1)
	}
}
