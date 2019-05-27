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
	"reflect"
	"strings"
)

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

import (
	"github.com/spatialcurrent/go-pipe/pkg/pipe"
	"github.com/spatialcurrent/go-stringify/pkg/stringify"
)

import (
	"github.com/spatialcurrent/go-simple-serializer/pkg/gss"
	"github.com/spatialcurrent/go-simple-serializer/pkg/iterator"
	"github.com/spatialcurrent/go-simple-serializer/pkg/jsonl"
	"github.com/spatialcurrent/go-simple-serializer/pkg/sv"
)

const (
	flagInputUri              string = "input-uri"
	flagInputCompression      string = "input-compression"
	flagInputFormat           string = "input-format"
	flagInputHeader           string = "input-header"
	flagInputLimit            string = "input-limit"
	flagInputComment          string = "input-comment"
	flagInputLazyQuotes       string = "input-lazy-quotes"
	flagInputTrim             string = "input-trim"
	flagInputReaderBufferSize string = "input-reader-buffer-size"
	flagInputSkipLines        string = "input-skip-lines"
	flagInputLineSeparator    string = "input-line-separator"
)

const (
	flagOutputUri               string = "output-uri"
	flagOutputCompression       string = "output-compression"
	flagOutputFormat            string = "output-format"
	flagOutputPretty            string = "output-pretty"
	flagOutputHeader            string = "output-header"
	flagOutputLimit             string = "output-limit"
	flagOutputAppend            string = "output-append"
	flagOutputOverwrite         string = "output-overwrite"
	flagOutputBufferMemory      string = "output-buffer-memory"
	flagOutputMkdirs            string = "output-mkdirs"
	flagOutputPassphrase        string = "output-passphrase"
	flagOutputSalt              string = "output-salt"
	flagOutputDecimal           string = "output-decimal"
	flagOutputNoDataValue       string = "output-no-data-value"
	flagOutputLineSeparator     string = "output-line-separator"
	flagOutputKeyValueSeparator string = "output-key-value-separator"
	flagOutputEscapePrefix      string = "output-escape-prefix"
	flagOutputEscapeSpace       string = "output-escape-space"
	flagOutputEscapeNewLine     string = "output-escape-new-line"
	flagOutputEscapeEqual       string = "output-escape-equal"
	flagOutputSorted            string = "output-sorted"
)

var gitBranch string
var gitCommit string

func buildValueSerializer(decimal bool, noDataValue string) func(object interface{}) (string, error) {
	if decimal {
		return stringify.DecimalValueStringer(noDataValue)
	}
	return stringify.DefaultValueStringer(noDataValue)
}

func buildWriter(outputWriter io.Writer, outputFormat string, outputHeader []interface{}, outputValueSerializer func(object interface{}) (string, error), outputLineSeparator string, outputPretty bool) (pipe.Writer, error) {
	if outputFormat == "csv" || outputFormat == "tsv" {
		separator, err := sv.FormatToSeparator(outputFormat)
		if err != nil {
			return nil, err
		}
		return sv.NewWriter(outputWriter, separator, outputHeader, outputValueSerializer), nil
	} else if outputFormat == "jsonl" {
		return jsonl.NewWriter(outputWriter, outputLineSeparator, outputPretty), nil
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

func initInputFlags(flag *pflag.FlagSet) {
	flag.StringP(flagInputFormat, "i", "", "The input format: "+strings.Join(gss.Formats, ", "))
	flag.StringSlice(flagInputHeader, []string{}, "The input header if the stdin input has no header.")
	flag.StringP(flagInputComment, "c", "", "The input comment character, e.g., #.  Commented lines are not sent to output.")
	flag.Bool(flagInputLazyQuotes, false, "allows lazy quotes for CSV and TSV")
	flag.Int(flagInputSkipLines, gss.NoSkip, "The number of lines to skip before processing")
	flag.IntP(flagInputLimit, "l", gss.NoLimit, "The input limit")
	flag.BoolP(flagInputTrim, "t", false, "trim input lines")
	flag.String(flagInputLineSeparator, "\n", "override line separator.  Used with properties and JSONL formats.")
}

func initOutputFlags(flag *pflag.FlagSet) {
	flag.StringP(flagOutputFormat, "o", "", "The output format: "+strings.Join(gss.Formats, ", "))
	flag.StringSlice(flagOutputHeader, []string{}, "The output header if the stdout output has no header.")
	flag.IntP(flagOutputLimit, "n", gss.NoLimit, "the output limit")
	flag.BoolP(flagOutputPretty, "p", false, "print pretty output")
	flag.BoolP(flagOutputSorted, "s", false, "sort output")
	flag.BoolP(flagOutputDecimal, "d", false, "when converting floats to strings use decimals rather than scientific notation")
	flag.StringP(flagOutputNoDataValue, "0", "", "no data value, e.g., used for missing values when converting JSON to CSV")
	flag.String(flagOutputLineSeparator, "\n", "override line separator.  Used with properties and JSONL formats.")
	flag.String(flagOutputKeyValueSeparator, "=", "override key value separator.  Used with properties format.")
	flag.String(flagOutputEscapePrefix, "", "override escape prefix.  Used with properties format.")
	flag.Bool(flagOutputEscapeSpace, false, "Escape space characters in output.  Used with properties format.")
	flag.Bool(flagOutputEscapeNewLine, false, "Escape new line characters in output.  Used with properties format.")
	flag.Bool(flagOutputEscapeEqual, false, "Escape equal characters in output.  Used with properties format.")
}

func initFlags(flag *pflag.FlagSet) {
	initInputFlags(flag)
	initOutputFlags(flag)

	flag.BoolP("async", "a", false, "async processing")
	flag.Bool("verbose", false, "Print debug info to stdout")
}

func CheckOutput(v *viper.Viper) error {
	if lineSepartor := v.GetString(flagOutputLineSeparator); len(lineSepartor) != 1 {
		return errors.New("line separator must be 1 character")
	}
	if len(v.GetString(flagOutputEscapePrefix)) == 0 {
		if v.GetBool(flagOutputEscapeSpace) {
			return errors.New("escaping space in output, but escape prefix is not set")
		}
		if v.GetBool(flagOutputEscapeNewLine) {
			return errors.New("escaping new line in output, but escape prefix is not set")
		}
		if v.GetBool(flagOutputEscapeEqual) {
			return errors.New("escaping equal in output, but escape prefix is not set")
		}
	}
	return nil
}

func checkConfig(v *viper.Viper) error {
	if err := CheckOutput(v); err != nil {
		return err
	}
	return nil
}

func main() {
	rootCommand := &cobra.Command{
		Use:                   "gss -i INPUT_FORMAT -o OUTPUT_FORMAT [--output-sorted]",
		DisableFlagsInUseLine: true,
		Short:                 "gss",
		Long:                  `gss is a simple program for serializing/deserializing data.`,
		RunE: func(cmd *cobra.Command, args []string) error {

			v := viper.New()

			if errorBind := v.BindPFlags(cmd.Flags()); errorBind != nil {
				return errorBind
			}
			v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
			v.AutomaticEnv()

			if errorConfig := checkConfig(v); errorConfig != nil {
				return errorConfig
			}

			inputFormat := v.GetString("input-format")

			if len(inputFormat) == 0 {
				return errors.New("input-format is required")
			}

			outputFormat := v.GetString("output-format")

			if len(outputFormat) == 0 {
				return errors.New("output-format is required")
			}

			outputHeader := stringify.StringSliceToInterfaceSlice(v.GetStringSlice("output-header"))

			outputSorted := v.GetBool("output-sorted")

			outputLineSeparator := v.GetString(flagOutputLineSeparator)

			outputValueSerializer := buildValueSerializer(v.GetBool("output-decimal"), v.GetString("output-no-data-value"))

			outputPretty := v.GetBool("output-pretty")

			if canStream(inputFormat, outputFormat, outputSorted) {

				p := pipe.NewBuilder()

				it, errorIterator := iterator.NewIterator(&iterator.NewIteratorInput{
					Reader:       os.Stdin,
					Type:         reflect.TypeOf([]map[string]interface{}{}),
					Format:       inputFormat,
					Header:       stringify.StringSliceToInterfaceSlice(v.GetStringSlice("input-header")),
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

				w, errorWriter := buildWriter(
					os.Stdout,
					outputFormat,
					outputHeader,
					outputValueSerializer,
					outputLineSeparator,
					outputPretty)
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
				InputBytes:              inputBytes,
				InputFormat:             inputFormat,
				InputHeader:             stringify.StringSliceToInterfaceSlice(v.GetStringSlice(flagInputHeader)),
				InputComment:            v.GetString(flagInputComment),
				InputLazyQuotes:         v.GetBool(flagInputLazyQuotes),
				InputSkipLines:          v.GetInt(flagInputSkipLines),
				InputLimit:              v.GetInt(flagInputLimit),
				InputLineSeparator:      v.GetString(flagInputLineSeparator),
				OutputFormat:            outputFormat,
				OutputHeader:            outputHeader,
				OutputLimit:             v.GetInt(flagOutputLimit),
				OutputPretty:            v.GetBool(flagOutputPretty),
				OutputSorted:            v.GetBool(flagOutputSorted),
				OutputValueSerializer:   outputValueSerializer,
				OutputLineSeparator:     v.GetString(flagOutputLineSeparator),
				OutputKeyValueSeparator: v.GetString(flagOutputKeyValueSeparator),
				OutputEscapePrefix:      v.GetString(flagOutputEscapePrefix),
				OutputEscapeSpace:       v.GetBool(flagOutputEscapeSpace),
				OutputEscapeNewLine:     v.GetBool(flagOutputEscapeNewLine),
				OutputEscapeEqual:       v.GetBool(flagOutputEscapeEqual),
				Async:                   v.GetBool("async"),
				Verbose:                 v.GetBool("verbose"),
			})
			if err != nil {
				return errors.Wrap(err, "error converting")
			}
			fmt.Println(outputString)
			return nil
		},
	}
	initFlags(rootCommand.Flags())

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
