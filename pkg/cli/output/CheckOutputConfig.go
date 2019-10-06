// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package output

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// CheckOutputConfig checks the output configuration.
func CheckOutputConfig(v *viper.Viper, formats []string) error {
	outputFormat := v.GetString(FlagOutputFormat)
	if len(outputFormat) == 0 {
		return &ErrMissingOutputFormat{Expected: formats}
	}
	if !stringSliceContains(formats, outputFormat) {
		return &ErrInvalidOutputFormat{Value: outputFormat, Expected: formats}
	}
	if ls := v.GetString(FlagOutputLineSeparator); len(ls) != 1 {
		if len(ls) == 0 {
			return ErrMissingOutputLineSeparator
		} else {
			return &ErrInvalidOutputLineSeparator{Value: ls}
		}
	}
	if kvs := v.GetString(FlagOutputKeyValueSeparator); len(kvs) != 1 {
		if len(kvs) == 0 {
			return ErrMissingOutputKeyValueSeparator
		} else {
			return &ErrInvalidOutputKeyValueSeparator{Value: kvs}
		}
	}
	if len(v.GetString(FlagOutputEscapePrefix)) == 0 {
		if v.GetBool(FlagOutputEscapeColon) {
			return errors.Wrap(ErrMissingOutputEscapePrefix, "escaping colon requires an escape prefix")
		}
		if v.GetBool(FlagOutputEscapeEqual) {
			return errors.Wrap(ErrMissingOutputEscapePrefix, "escaping equal requires an escape prefix")
		}
		if v.GetBool(FlagOutputEscapeSpace) {
			return errors.Wrap(ErrMissingOutputEscapePrefix, "escaping space requires an escape prefix")
		}
		if v.GetBool(FlagOutputEscapeNewLine) {
			return errors.Wrap(ErrMissingOutputEscapePrefix, "escaping new line requires an escape prefix")
		}
	}
	if v.GetBool(FlagOutputKeyLower) && v.GetBool(FlagOutputKeyUpper) {
		return errors.New("cannot lower case and upper case keys at the same time")
	}
	return nil
}
