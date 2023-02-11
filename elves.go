/*
Copyright © 2023 nukoneko-tarou
*/
package elves

import (
	"os"

	"github.com/nukoneko-tarou/elves/commands"
	"github.com/spf13/cobra"
)

func main() {
	cmd := cobra.Command{
		Use:   "elves",
		Short: "Tool to generate directories from json files",
		Long:  `Tool to generate directories from json files`,
	}

	cmd.AddCommand(commands.NewVersion().Cmd)
	cmd.AddCommand(commands.NewCreate().Cmd)

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
