// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package formats

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/spatialcurrent/go-simple-serializer/pkg/gss"
	"github.com/spatialcurrent/go-simple-serializer/pkg/serializer"
)

// NewCommand returns a new instance of the formats command.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   CliUse,
		Short: CliShort,
		Long:  CliLong,
		RunE: func(cmd *cobra.Command, args []string) error {
			v := viper.New()

			err := v.BindPFlags(cmd.Flags())
			if err != nil {
				return errors.Wrap(err, "error binding flags")
			}
			v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
			v.AutomaticEnv() // set environment variables to overwrite config

			err = CheckFormatsConfig(v)
			if err != nil {
				return errors.Wrap(err, "error with configuration")
			}

			f := v.GetString(FlagFormat)

			b, err := serializer.New(f).LineSeparator("\n").Serialize(gss.Formats)
			if err != nil {
				return errors.Wrapf(err, "error serializing formats with format %q", f)
			}

			fmt.Print(string(b))

			return nil
		},
	}
	InitFormatsFlags(cmd.Flags())
	return cmd
}
