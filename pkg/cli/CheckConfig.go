// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

// Package cli manages loading and testing configuration for serializing and deserializing objects from the command line.
package cli

import (
	"github.com/pkg/errors"
	"github.com/spatialcurrent/go-simple-serializer/pkg/cli/input"
	"github.com/spatialcurrent/go-simple-serializer/pkg/cli/output"
	"github.com/spf13/viper"
)

// CheckConfig checks the configuration.
func CheckConfig(v *viper.Viper, formats []string) error {
	err := input.CheckInputConfig(v, formats)
	if err != nil {
		return errors.Wrap(err, "error with input configuration")
	}
	err = output.CheckOutputConfig(v, formats)
	if err != nil {
		return errors.Wrap(err, "error with output configuration")
	}
	return nil
}
