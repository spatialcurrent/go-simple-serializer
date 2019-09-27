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
	"io/ioutil"
	"os"
	"reflect"
	"strings"

	"github.com/pkg/errors"
	"github.com/spatialcurrent/go-pipe/pkg/pipe"
	"github.com/spatialcurrent/go-simple-serializer/pkg/cli"
	"github.com/spatialcurrent/go-simple-serializer/pkg/cli/formats"
	"github.com/spatialcurrent/go-simple-serializer/pkg/cli/version"
	"github.com/spatialcurrent/go-simple-serializer/pkg/gob"
	"github.com/spatialcurrent/go-simple-serializer/pkg/gss"
	"github.com/spatialcurrent/go-simple-serializer/pkg/iterator"
	"github.com/spatialcurrent/go-simple-serializer/pkg/properties"
	"github.com/spatialcurrent/go-simple-serializer/pkg/serializer"
	"github.com/spatialcurrent/go-simple-serializer/pkg/writer"
	"github.com/spatialcurrent/go-stringify/pkg/stringify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var gitBranch string
var gitCommit string

var (
	mapStringInterfaceType      = reflect.TypeOf(map[string]interface{}{})
	sliceMapStringInterfaceType = reflect.TypeOf([]map[string]interface{}{})
)

func main() {

	// Register gob types
	gob.RegisterTypes()

	rootCommand := &cobra.Command{
		Use: "gss -i INPUT_FORMAT -o OUTPUT_FORMAT",
		DisableFlagsInUseLine: false,
		Short:         "gss",
		Long:          `gsss is a simple and fast program for serializing/deserializing data that supports following file formats: ` + strings.Join(gss.Formats, ", "),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			v := viper.New()

			if errorBind := v.BindPFlags(cmd.Flags()); errorBind != nil {
				return errorBind
			}
			v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
			v.AutomaticEnv()

			if errorConfig := cli.CheckConfig(v, serializer.Formats); errorConfig != nil {
				return errorConfig
			}

			inputFormat := v.GetString(cli.FlagInputFormat)

			inputHeader := stringify.StringSliceToInterfaceSlice(v.GetStringSlice(cli.FlagInputHeader))

			inputLineSeparator := v.GetString(cli.FlagInputLineSeparator)

			outputFormat := v.GetString(cli.FlagOutputFormat)

			outputHeader := stringify.StringSliceToInterfaceSlice(v.GetStringSlice(cli.FlagOutputHeader))

			outputSorted := v.GetBool(cli.FlagOutputSorted)
			outputReversed := v.GetBool(cli.FlagOutputReversed)

			outputKeyValueSeparator := v.GetString(cli.FlagOutputKeyValueSeparator)
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

			outputFit := v.GetBool(cli.FlagOutputFit)

			outputLimit := v.GetInt(cli.FlagOutputLimit)

			verbose := v.GetBool(cli.FlagVerbose)

			if verbose {
				err := properties.Write(&properties.WriteInput{
					Writer:            os.Stdout,
					LineSeparator:     "\n",
					KeyValueSeparator: ":",
					Object:            v.AllSettings(),
					KeySerializer:     stringify.NewDefaultStringer(),
					ValueSerializer:   stringify.NewDefaultStringer(),
					Sorted:            true,
					Reversed:          false,
					EscapePrefix:      "\\",
					EscapeSpace:       false,
					EscapeEqual:       true,
					EscapeColon:       false,
					EscapeNewLine:     true,
				})
				if err != nil {
					return errors.Wrap(err, "error writing viper settings")
				}
				fmt.Println("")
			}

			noStream := v.GetBool("no-stream")

			if (!noStream) && gss.CanStream(inputFormat, outputFormat, outputSorted) {

				if verbose {
					fmt.Println("Streaming: yes")
				}

				p := pipe.NewBuilder()

				var inputType reflect.Type
				if inputFormat == "gob" {
					inputType = mapStringInterfaceType
				}

				if verbose {
					fmt.Println("Streaming Type:", inputType)
				}

				it, errorIterator := iterator.NewIterator(&iterator.NewIteratorInput{
					Reader:            os.Stdin,
					Type:              inputType,
					Format:            inputFormat,
					Header:            inputHeader,
					ScannerBufferSize: v.GetInt(cli.FlagInputScannerBufferSize),
					SkipLines:         v.GetInt(cli.FlagInputSkipLines),
					SkipBlanks:        true,
					SkipComments:      true,
					Comment:           v.GetString(cli.FlagInputComment),
					Trim:              v.GetBool(cli.FlagInputTrim),
					LineSeparator:     inputLineSeparator,
					DropCR:            v.GetBool(cli.FlagInputDropCR),
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

				w, errorWriter := writer.NewWriter(&writer.NewWriterInput{
					Writer:            os.Stdout,
					Format:            outputFormat,
					FormatSpecifier:   v.GetString(cli.FlagOutputFormatSpecifier),
					Header:            outputHeader,
					KeySerializer:     outputKeySerializer,
					ValueSerializer:   outputValueSerializer,
					KeyValueSeparator: outputKeyValueSeparator,
					LineSeparator:     outputLineSeparator,
					Fit:               outputFit,
					Pretty:            outputPretty,
					Sorted:            outputSorted,
					Reversed:          outputReversed,
				})
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

			var inputType reflect.Type
			if inputFormat == "gob" {
				inputType = sliceMapStringInterfaceType
			}

			outputBytes, err := gss.Convert(&gss.ConvertInput{
				InputBytes:              inputBytes,
				InputFormat:             inputFormat,
				InputHeader:             inputHeader,
				InputComment:            v.GetString(cli.FlagInputComment),
				InputLazyQuotes:         v.GetBool(cli.FlagInputLazyQuotes),
				InputScannerBufferSize:  v.GetInt(cli.FlagInputScannerBufferSize),
				InputSkipLines:          v.GetInt(cli.FlagInputSkipLines),
				InputLimit:              v.GetInt(cli.FlagInputLimit),
				InputLineSeparator:      inputLineSeparator,
				InputDropCR:             v.GetBool(cli.FlagInputDropCR),
				InputEscapePrefix:       v.GetString(cli.FlagInputEscapePrefix),
				InputUnescapeSpace:      v.GetBool(cli.FlagInputUnescapeSpace),
				InputUnescapeNewLine:    v.GetBool(cli.FlagInputUnescapeNewLine),
				InputUnescapeEqual:      v.GetBool(cli.FlagInputUnescapeEqual),
				InputTrim:               v.GetBool(cli.FlagInputTrim),
				InputType:               inputType,
				OutputFormat:            outputFormat,
				OutputFormatSpecifier:   v.GetString(cli.FlagOutputFormatSpecifier),
				OutputFit:               outputFit,
				OutputHeader:            outputHeader,
				OutputLimit:             outputLimit,
				OutputPretty:            v.GetBool(cli.FlagOutputPretty),
				OutputSorted:            outputSorted,
				OutputReversed:          outputReversed,
				OutputKeySerializer:     outputKeySerializer,
				OutputValueSerializer:   outputValueSerializer,
				OutputLineSeparator:     outputLineSeparator,
				OutputKeyValueSeparator: outputKeyValueSeparator,
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
	cli.InitFlags(rootCommand.Flags())

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

	rootCommand.AddCommand(formats.NewCommand())

	rootCommand.AddCommand(version.NewCommand(&version.NewCommandInput{
		GitBranch: gitBranch,
		GitCommit: gitCommit,
	}))

	if err := rootCommand.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: "+err.Error()+"\n")
		fmt.Fprintf(os.Stderr, "Help: gss -h\n")
		os.Exit(1)
	}
}
