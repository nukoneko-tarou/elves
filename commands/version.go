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
		RunE:  cmd.run,
	}

	return cmd
}

func (c *Version) run(cmd *cobra.Command, _ []string) error {
	fmt.Println("elves version 0.1")

	return nil
}
