/*
Copyright © 2023 nukoneko-tarou
*/
package commands

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/nukoneko-tarou/elves/internal/tree"
	"github.com/spf13/cobra"
)

type Create struct {
	Cmd *cobra.Command
}

func NewCreate() *Create {
	cmd := &Create{}
	cmd.Cmd = &cobra.Command{
		Use:   "create",
		Short: "Create a directory",
		Long:  `Creates a directory where the command is executed`,
		Args:  cobra.MinimumNArgs(1),
		RunE:  cmd.run,
		Example: strings.Join([]string{
			"elves create ./sample.json",
			"elves create ./sample.json --sub new-project --permission 777",
		}, "\n"),
	}

	cmd.Cmd.Flags().StringP("permission", "p", "", "permission")
	cmd.Cmd.Flags().StringP("sub", "s", "", "Create in subdirectory")
	cmd.Cmd.Flags().BoolP("gitkeep", "g", false, "Create .gitkeep")
	cmd.Cmd.Flags().BoolP("dry-run", "d", false, "Dry run")
	cmd.Cmd.Flags().IntP("concurrency", "c", tree.DefaultConcurrency, "Concurrent workers for directory creation")

	return cmd
}

func (c *Create) run(cmd *cobra.Command, args []string) error {
	p, err := cmd.Flags().GetString("permission")
	if err != nil {
		return err
	}
	sn, err := cmd.Flags().GetString("sub")
	if err != nil {
		return err
	}
	gk, err := cmd.Flags().GetBool("gitkeep")
	if err != nil {
		return err
	}
	dr, err := cmd.Flags().GetBool("dry-run")
	if err != nil {
		return err
	}
	concurrency, err := cmd.Flags().GetInt("concurrency")
	if err != nil {
		return err
	}

	file, err := os.Open(args[0])
	if err != nil {
		return err
	}
	defer file.Close()

	treeData, err := tree.ParseJSON(file)
	if err != nil {
		return err
	}

	rootNodes := (*treeData)[0].Contents

	if dr {
		rootName := "."
		if sn != "" {
			rootName = sn
		}
		tree.Render(cmd.OutOrStdout(), rootName, rootNodes, tree.RenderOptions{
			Gitkeep: gk,
		})
		return nil
	}

	permissionString := "755"
	if p != "" {
		permissionString = p
	}

	permission, err := strconv.ParseUint(permissionString, 8, 32)
	if err != nil {
		return fmt.Errorf("invalid permission value: %w", err)
	}

	return tree.Build(rootNodes, tree.BuildOptions{
		BaseDir:     ".",
		SubDir:      sn,
		Permission:  os.FileMode(permission),
		Gitkeep:     gk,
		Concurrency: concurrency,
	})
}
