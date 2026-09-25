/*
Copyright © 2023 nukoneko-tarou
*/
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/nukoneko-tarou/elves/commands"
	"github.com/spf13/cobra"
)

func main() {
	cmd := cobra.Command{
		Use:   "elves",
		Short: "Fast CLI to scaffold directory structures from tree -J JSON",
		Long: strings.TrimSpace(`
Elves is a high-performance CLI tool that scaffolds project directory structures
from JSON files compatible with UNIX 'tree -J' command output.`),
		Example: strings.Join([]string{
			"  elves create ./sample.json",
			"  elves create ./sample.json --dry-run",
			"  elves create ./sample.json --sub my-project --permission 755",
			"  elves version",
		}, "\n"),
		Version: commands.Version,
	}

	cmd.AddCommand(commands.NewVersion().Cmd)
	cmd.AddCommand(commands.NewCreate().Cmd)

	if err := cmd.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
