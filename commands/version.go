/*
Copyright © 2023 nukoneko-tarou
*/
package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version = "1.1.0"

type VersionCommand struct {
	Cmd *cobra.Command
}

func NewVersion() *VersionCommand {
	cmd := &VersionCommand{}
	cmd.Cmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of elves",
		Long:  `Print the version number of elves`,
		Run:   cmd.run,
	}

	return cmd
}

func (c *VersionCommand) run(cmd *cobra.Command, _ []string) {
	fmt.Printf("elves version %s\n", Version)
}
