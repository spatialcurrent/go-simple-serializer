// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

// gss is the command line program for go-simple-serializer (GSS).
//
// Usage
//
// Use `gss help` to see full help documentation.
//
//	gss -i INPUT_FORMAT -o OUTPUT_FORMAT [flags]
//
// Examples
//
//	# convert .gitignore to JSON
//	cat .gitignore | gss -i csv --input-header path -o json
//
//	# extract version from CircleCI config
//	cat .circleci/config.yml | gss -i yaml -o json -c '#' | jq -r .version
//
//	# convert list of files to JSON Lines
//	find . -name '*.go' | gss -i csv --input-header path -o jsonl
package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"strings"

	"github.com/pkg/errors"
	"github.com/spatialcurrent/go-pipe/pkg/pipe"
	"github.com/spatialcurrent/go-simple-serializer/pkg/cli"
	"github.com/spatialcurrent/go-simple-serializer/pkg/gss"
	"github.com/spatialcurrent/go-simple-serializer/pkg/iterator"
	"github.com/spatialcurrent/go-simple-serializer/pkg/jsonl"
	"github.com/spatialcurrent/go-simple-serializer/pkg/serializer"
	"github.com/spatialcurrent/go-simple-serializer/pkg/sv"
	"github.com/spatialcurrent/go-stringify/pkg/stringify"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var gitBranch string
var gitCommit string

func buildWriter(outputWriter io.Writer, outputFormat string, outputHeader []interface{}, outputKeySerializer stringify.Stringer, outputValueSerializer stringify.Stringer, outputLineSeparator string, outputPretty bool, outputSorted bool, outputReversed bool) (pipe.Writer, error) {
	if outputFormat == "csv" || outputFormat == "tsv" {
		separator, err := sv.FormatToSeparator(outputFormat)
		if err != nil {
			return nil, err
		}
		w := sv.NewWriter(
			outputWriter,
			separator,
			outputHeader,
			outputKeySerializer,
			outputValueSerializer,
			outputSorted,
			outputReversed,
		)
		return w, nil
	} else if outputFormat == "jsonl" {
		return jsonl.NewWriter(outputWriter, outputLineSeparator, outputKeySerializer, outputPretty), nil
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

func initFlags(flag *pflag.FlagSet, formats []string) {
	cli.InitInputFlags(flag, formats)
	cli.InitOutputFlags(flag, formats)
}

func checkConfig(v *viper.Viper, formats []string) error {
	if err := cli.CheckInput(v, formats); err != nil {
		return err
	}
	if err := cli.CheckOutput(v, formats); err != nil {
		return err
	}
	return nil
}

func main() {

	formats := serializer.Formats

	rootCommand := &cobra.Command{
		Use:                   "gss -i INPUT_FORMAT -o OUTPUT_FORMAT",
		DisableFlagsInUseLine: false,
		Short:                 "gss",
		Long:                  `gss is a simple program for serializing/deserializing data.`,
		SilenceUsage:          true,
		SilenceErrors:         true,
		RunE: func(cmd *cobra.Command, args []string) error {

			v := viper.New()

			if errorBind := v.BindPFlags(cmd.Flags()); errorBind != nil {
				return errorBind
			}
			v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
			v.AutomaticEnv()

			if errorConfig := checkConfig(v, formats); errorConfig != nil {
				return errorConfig
			}

			inputFormat := v.GetString(cli.FlagInputFormat)

			inputHeader := stringify.StringSliceToInterfaceSlice(v.GetStringSlice(cli.FlagInputHeader))

			outputFormat := v.GetString(cli.FlagOutputFormat)

			outputHeader := stringify.StringSliceToInterfaceSlice(v.GetStringSlice(cli.FlagOutputHeader))

			outputSorted := v.GetBool(cli.FlagOutputSorted)
			outputReversed := v.GetBool(cli.FlagOutputReversed)

			outputLineSeparator := v.GetString(cli.FlagOutputLineSeparator)

			outputKeySerializer := stringify.NewStringer(
				"",
				v.GetBool(cli.FlagOutputDecimal),
				v.GetBool(cli.FlagOutputKeyLower),
				v.GetBool(cli.FlagOutputKeyUpper),
			)

			outputValueSerializer := stringify.NewStringer(
				v.GetString(cli.FlagOutputNoDataValue),
				v.GetBool(cli.FlagOutputDecimal),
				v.GetBool(cli.FlagOutputValueLower),
				v.GetBool(cli.FlagOutputValueUpper),
			)

			outputPretty := v.GetBool(cli.FlagOutputPretty)

			outputLimit := v.GetInt(cli.FlagOutputLimit)

			if canStream(inputFormat, outputFormat, outputSorted) {

				p := pipe.NewBuilder()

				it, errorIterator := iterator.NewIterator(&iterator.NewIteratorInput{
					Reader:       os.Stdin,
					Type:         reflect.TypeOf([]map[string]interface{}{}),
					Format:       inputFormat,
					Header:       inputHeader,
					SkipLines:    v.GetInt(cli.FlagInputSkipLines),
					SkipBlanks:   true,
					SkipComments: true,
					Comment:      v.GetString(cli.FlagInputComment),
					Trim:         v.GetBool(cli.FlagInputTrim),
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
					outputKeySerializer,
					outputValueSerializer,
					outputLineSeparator,
					outputPretty,
					outputSorted,
					outputReversed,
				)
				if errorWriter != nil {
					return errors.Wrap(errorWriter, "error building output writer")
				}
				p = p.Output(w)

				if outputLimit >= 0 {
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

			outputBytes, err := gss.Convert(&gss.ConvertInput{
				InputBytes:              inputBytes,
				InputFormat:             inputFormat,
				InputHeader:             inputHeader,
				InputComment:            v.GetString(cli.FlagInputComment),
				InputLazyQuotes:         v.GetBool(cli.FlagInputLazyQuotes),
				InputSkipLines:          v.GetInt(cli.FlagInputSkipLines),
				InputLimit:              v.GetInt(cli.FlagInputLimit),
				InputLineSeparator:      v.GetString(cli.FlagInputLineSeparator),
				InputEscapePrefix:       v.GetString(cli.FlagInputEscapePrefix),
				InputUnescapeSpace:      v.GetBool(cli.FlagInputUnescapeSpace),
				InputUnescapeNewLine:    v.GetBool(cli.FlagInputUnescapeNewLine),
				InputUnescapeEqual:      v.GetBool(cli.FlagInputUnescapeEqual),
				OutputFormat:            outputFormat,
				OutputHeader:            outputHeader,
				OutputLimit:             outputLimit,
				OutputPretty:            v.GetBool(cli.FlagOutputPretty),
				OutputSorted:            outputSorted,
				OutputReversed:          outputReversed,
				OutputKeySerializer:     outputKeySerializer,
				OutputValueSerializer:   outputValueSerializer,
				OutputLineSeparator:     v.GetString(cli.FlagOutputLineSeparator),
				OutputKeyValueSeparator: v.GetString(cli.FlagOutputKeyValueSeparator),
				OutputEscapePrefix:      v.GetString(cli.FlagOutputEscapePrefix),
				OutputEscapeSpace:       v.GetBool(cli.FlagOutputEscapeSpace),
				OutputEscapeNewLine:     v.GetBool(cli.FlagOutputEscapeNewLine),
				OutputEscapeEqual:       v.GetBool(cli.FlagOutputEscapeEqual),
			})
			if err != nil {
				return errors.Wrap(err, "error converting")
			}
			switch outputFormat {
			case serializer.FormatCSV, serializer.FormatJSONL, serializer.FormatProperties, serializer.FormatTags, serializer.FormatTOML, serializer.FormatTSV, serializer.FormatYAML:
				// do not include trailing new line, since it comes with the output
				fmt.Print(string(outputBytes))
			default:
				// print trailing new line for all others
				fmt.Println(string(outputBytes))
			}
			return nil
		},
	}
	initFlags(rootCommand.Flags(), formats)

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
		fmt.Fprintf(os.Stderr, "Error: "+err.Error()+"\n")
		fmt.Fprintf(os.Stderr, "Help: gss -h\n")
		os.Exit(1)
	}
}
