// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package formats

import (
	"github.com/spf13/pflag"
)

// InitFormatsFlags initializes the flags for the formats command.
func InitFormatsFlags(flag *pflag.FlagSet) {
	flag.StringP(FlagFormat, "f", DefaultFormat, "output format")
}
