/*
Copyright © 2023 nukoneko-tarou
*/
package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

type Version struct {
	Cmd *cobra.Command
}

func NewVersion() *Version {
	cmd := &Version{}
	cmd.Cmd = &cobra.Command{
		Use:   "version",
		Short: "version",
		Long:  `version`,
		Run:   cmd.run,
	}

	return cmd
}

func (c *Version) run(cmd *cobra.Command, _ []string) {
	fmt.Println("elves version 1.1.0")
}
