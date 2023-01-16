/*
Copyright © 2023 nukoneko-tarou
*/
package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

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
		Short: "create",
		Long:  `create`,
		RunE:  cmd.run,
	}

	return cmd
}

func (c *Create) run(cmd *cobra.Command, _ []string) error {
	tree, err := c.perseJson()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	err = c.createDirectory(tree[0].Contents)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return nil
}

func (c *Create) perseJson() (TreeJson, error) {
	raw, err := ioutil.ReadFile("./sample.json")
	if err != nil {
		return nil, err
	}

	var fc TreeJson

	if err := json.Unmarshal(raw, &fc); err != nil {
		return nil, err
	}

	return fc, nil
}

func (c *Create) createDirectory(contents []DirectoryJson) error {
	for _, v := range contents {
		if v.Type == "directory" {
			if err := os.Mkdir(v.Name, 0777); err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		}
	}

	return nil
}
