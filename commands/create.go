/*
Copyright © 2023 nukoneko-tarou
*/
package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

type DirectoryJson struct {
	Type     string          `json:"type"`
	Name     string          `json:"name"`
	Contents []DirectoryJson `json:"contents,omitempty"`
}

type TreeJson []struct {
	Type        string          `json:"type"`
	Name        string          `json:"name,omitempty"`
	Contents    []DirectoryJson `json:"contents,omitempty"`
	Directories int             `json:"directories,omitempty"`
	Files       int             `json:"files,omitempty"`
}

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

	raw, err := os.ReadFile(args[0])
	if err != nil {
		return err
	}

	tree, err := c.perseJson(raw)
	if err != nil {
		return err
	}
	if len(tree) == 0 || len(tree[0].Contents) == 0 {
		return errors.New("invalid or empty json file")
	}

	if dr {
		subDirectoryName := "."
		if sn != "" {
			subDirectoryName = sn
		}
		fmt.Fprintln(cmd.OutOrStdout(), subDirectoryName)
		c.printDirectoryTree(cmd.OutOrStdout(), tree[0].Contents, "", gk)

		return nil
	}

	permissionString := "755"
	subDirectoryName := "./"

	if p != "" {
		permissionString = p
	}

	permission, err := strconv.ParseUint(permissionString, 8, 32)
	if err != nil {
		return fmt.Errorf("invalid permission value: %w", err)
	}

	if sn != "" {
		path := filepath.Join(".", sn)
		if err := os.Mkdir(path, os.FileMode(permission)); err != nil {
			return err
		}
		subDirectoryName = path
	}

	err = c.createDirectory(tree[0].Contents, subDirectoryName, os.FileMode(permission), gk)
	if err != nil {
		return err
	}

	return nil
}

func (c *Create) printDirectoryTree(w io.Writer, contents []DirectoryJson, prefix string, isCreateGitkeep bool) {
	for i, v := range contents {
		if v.Type == "directory" {
			connector := "├── "
			if i == len(contents)-1 {
				connector = "└── "
			}
			fmt.Fprintf(w, "%s%s%s\n", prefix, connector, v.Name)

			newPrefix := prefix + "│   "
			if i == len(contents)-1 {
				newPrefix = prefix + "    "
			}

			if isCreateGitkeep {
				gitkeepConnector := "├── "
				if len(v.Contents) == 0 {
					gitkeepConnector = "└── "
				}
				fmt.Fprintf(w, "%s%s.gitkeep\n", newPrefix, gitkeepConnector)
			}

			if len(v.Contents) != 0 {
				c.printDirectoryTree(w, v.Contents, newPrefix, isCreateGitkeep)
			}
		}
	}
}

func (c *Create) perseJson(raw []byte) (TreeJson, error) {
	var fc TreeJson
	if err := json.Unmarshal(raw, &fc); err != nil {
		return nil, err
	}

	return fc, nil
}

func (c *Create) createDirectory(contents []DirectoryJson, parentName string, permission os.FileMode, isCreateGitkeep bool) error {
	for _, v := range contents {
		if v.Type == "directory" {
			path := filepath.Join(parentName, v.Name)

			if err := os.Mkdir(path, permission); err != nil {
				return err
			}

			if isCreateGitkeep {
				if _, err := os.Create(filepath.Join(path, ".gitkeep")); err != nil {
					return err
				}
			}

			if len(v.Contents) != 0 {
				if err := c.createDirectory(v.Contents, path, permission, isCreateGitkeep); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
