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
	"github.com/spatialcurrent/go-simple-serializer/gss"
)

var gitTag string
var gitBranch string
var gitCommit string

func main() {
	root := &cobra.Command{
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

			inputBytes, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				return errors.Wrap(err, "error reading from stdin")
			}

			outputString, err := gss.Convert(&gss.ConvertInput{
				InputBytes:      inputBytes,
				InputFormat:     inputFormat,
				InputHeader:     v.GetStringSlice("input-header"),
				InputComment:    v.GetString("input-comment"),
				InputLazyQuotes: v.GetBool("input-lazy-quotes"),
				InputSkipLines:  v.GetInt("input-skip-lines"),
				InputLimit:      v.GetInt("input-limit"),
				OutputFormat:    outputFormat,
				OutputHeader:    v.GetStringSlice("output-header"),
				OutputLimit:     v.GetInt("output-limit"),
				OutputPretty:    v.GetBool("output-pretty"),
				Async:           v.GetBool("async"),
				Verbose:         v.GetBool("verbose"),
			})
			if err != nil {
				return errors.Wrap(err, "error converting")
			}
			fmt.Println(outputString)
			return nil
		},
	}
	flags := root.Flags()
	flags.StringP("input-format", "i", "", "The input format: "+strings.Join(gss.Formats, ", "))
	flags.StringSlice("input-header", []string{}, "The input header if the stdin input has no header.")
	flags.StringP("input-comment", "c", "", "The input comment character, e.g., #.  Commented lines are not sent to output.")
	flags.Bool("input-lazy-quotes", false, "allows lazy quotes for CSV and TSV")
	flags.Int("input-skip-lines", gss.NoSkip, "The number of lines to skip before processing")
	flags.IntP("input-limit", "l", gss.NoLimit, "The input limit")
	flags.StringP("output-format", "o", "", "The output format: "+strings.Join(gss.Formats, ", "))
	flags.StringSlice("output-header", []string{}, "The output header if the stdout output has no header.")
	flags.Int("output-limit", gss.NoLimit, "the output limit")
	flags.BoolP("output-pretty", "p", false, "print pretty output")
	flags.BoolP("async", "a", false, "async processing")
	flags.Bool("verbose", false, "Print debug info to stdout")

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

	root.AddCommand(version)

	if err := root.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, err.Error()+"\n")
		os.Exit(1)
	}
}
