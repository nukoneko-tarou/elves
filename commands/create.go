/*
Copyright © 2023 nukoneko-tarou
*/
package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
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

	return cmd
}

func (c *Create) run(cmd *cobra.Command, args []string) error {
	permission := "755"
	subDirectoryName := "./"

	p, _ := cmd.Flags().GetString("permission")
	sn, _ := cmd.Flags().GetString("sub")

	if p != "" {
		permission = p
	}

	if sn != "" {
		path := "./" + sn
		perm32, _ := strconv.ParseUint(permission, 8, 32)
		err := os.Mkdir(path, os.FileMode(perm32))
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		subDirectoryName = path
	}

	raw, err := ioutil.ReadFile(args[0])
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	tree, err := c.perseJson(raw)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	err = c.createDirectory(tree[0].Contents, subDirectoryName, permission)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return nil
}

func (c *Create) perseJson(raw []byte) (TreeJson, error) {
	var fc TreeJson
	if err := json.Unmarshal(raw, &fc); err != nil {
		return nil, err
	}

	return fc, nil
}

func (c *Create) createDirectory(contents []DirectoryJson, parentName string, permission string) error {
	for _, v := range contents {
		if v.Type == "directory" {
			path := parentName + "/" + v.Name
			perm32, _ := strconv.ParseUint(permission, 8, 32)

			if err := os.Mkdir(path, os.FileMode(perm32)); err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			if len(v.Contents) != 0 {
				c.createDirectory(v.Contents, path, permission)
			}
		}
	}

	return nil
}
