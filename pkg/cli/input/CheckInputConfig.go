// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package input

import (
	"fmt"

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
			return fmt.Errorf("unescaping colon requires an escape prefix: %w", ErrMissingInputEscapePrefix)
		}
		if v.GetBool(FlagInputUnescapeEqual) {
			return fmt.Errorf("unescaping equal requires an escape prefix: %w", ErrMissingInputEscapePrefix)
		}
		if v.GetBool(FlagInputUnescapeSpace) {
			return fmt.Errorf("unescaping space requires an escape prefix: %w", ErrMissingInputEscapePrefix)
		}
		if v.GetBool(FlagInputUnescapeNewLine) {
			return fmt.Errorf("unescaping new line requires an escape prefix: %w", ErrMissingInputEscapePrefix)
		}
	}
	inputComment := v.GetString(FlagInputComment)
	if (inputFormat == "csv" || inputFormat == "tsv") && len(inputComment) > 1 {
		return &ErrInvalidInputComment{Value: inputComment}
	}
	return nil
}
