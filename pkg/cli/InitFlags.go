// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package cli

import (
	"github.com/spf13/pflag"

	"github.com/spatialcurrent/go-simple-serializer/pkg/cli/input"
	"github.com/spatialcurrent/go-simple-serializer/pkg/cli/output"
)

// Initialize cli flags
func InitFlags(flag *pflag.FlagSet) {

	input.InitInputFlags(flag)

	output.InitOutputFlags(flag)

	flag.Bool(FlagNoStream, false, "disable streaming")

	flag.BoolP(FlagVerbose, "v", false, "verbose output")
}
