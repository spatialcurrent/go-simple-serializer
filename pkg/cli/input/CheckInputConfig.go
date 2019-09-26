// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package input

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// CheckInputConfig checks the output configuration.
func CheckInputConfig(v *viper.Viper, formats []string) error {
	inputFormat := v.GetString(FlagInputFormat)
	if len(inputFormat) == 0 {
		return &ErrMissingInputFormat{Expected: formats}
	}
	if !stringSliceContains(formats, inputFormat) {
		return &ErrInvalidInputFormat{Value: inputFormat, Expected: formats}
	}
	if ls := v.GetString(FlagInputLineSeparator); len(ls) != 1 {
		if len(ls) == 0 {
			return ErrMissingInputLineSeparator
		} else {
			return &ErrInvalidInputLineSeparator{Value: ls}
		}
	}
	if kvs := v.GetString(FlagInputKeyValueSeparator); len(kvs) != 1 {
		if len(kvs) == 0 {
			return ErrMissingInputKeyValueSeparator
		} else {
			return &ErrInvalidInputKeyValueSeparator{Value: kvs}
		}
	}
	if len(v.GetString(FlagInputEscapePrefix)) == 0 {
		if v.GetBool(FlagInputUnescapeColon) {
			return errors.Wrap(ErrMissingInputEscapePrefix, "unescaping colon requires an escape prefix")
		}
		if v.GetBool(FlagInputUnescapeEqual) {
			return errors.Wrap(ErrMissingInputEscapePrefix, "unescaping equal requires an escape prefix")
		}
		if v.GetBool(FlagInputUnescapeSpace) {
			return errors.Wrap(ErrMissingInputEscapePrefix, "unescaping space requires an escape prefix")
		}
		if v.GetBool(FlagInputUnescapeNewLine) {
			return errors.Wrap(ErrMissingInputEscapePrefix, "unescaping new line requires an escape prefix")
		}
	}
	inputComment := v.GetString(FlagInputComment)
	if (inputFormat == "csv" || inputFormat == "tsv") && len(inputComment) > 1 {
		return &ErrInvalidInputComment{Value: inputComment}
	}
	return nil
}
