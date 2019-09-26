// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

type NewCommandInput struct {
	GitBranch string
	GitCommit string
}

// NewCommand returns a new instance of the version command.
func NewCommand(input *NewCommandInput) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "print version information to stdout",
		Long:  "print version information to stdout",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Branch: " + input.GitBranch)
			fmt.Println("Commit: " + input.GitCommit)
		},
	}
	return cmd
}
