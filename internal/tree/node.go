package tree

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"
)

// NodeType represents the type of entry in the tree JSON.
type NodeType string

const (
	NodeTypeDirectory NodeType = "directory"
	NodeTypeFile      NodeType = "file"
	NodeTypeReport    NodeType = "report"
)

// Node represents an entry in tree -J JSON output.
type Node struct {
	Type     NodeType `json:"type"`
	Name     string   `json:"name"`
	Contents []Node   `json:"contents,omitempty"`

	// Report metadata from tree -J
	Directories int `json:"directories,omitempty"`
	Files       int `json:"files,omitempty"`
}

// Tree represents the root array structure of tree -J JSON output.
type Tree []Node

// ParseJSON parses and validates the JSON content from an io.Reader.
func ParseJSON(r io.Reader) (*Tree, error) {
	var tree Tree
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&tree); err != nil {
		return nil, fmt.Errorf("failed to decode json: %w", err)
	}

	if len(tree) == 0 {
		return nil, errors.New("invalid or empty json file")
	}

	// Find the root directory entry (typically tree[0])
	rootNode := &tree[0]
	if len(rootNode.Contents) == 0 {
		return nil, errors.New("invalid or empty json file")
	}

	// Validate nodes recursively for path traversal or invalid names
	if err := validateNodes(rootNode.Contents); err != nil {
		return nil, err
	}

	return &tree, nil
}

// validateNodes recursively validates node names against path traversal.
func validateNodes(nodes []Node) error {
	for _, node := range nodes {
		if node.Type == NodeTypeDirectory || node.Type == NodeTypeFile {
			if err := validateName(node.Name); err != nil {
				return err
			}
			if len(node.Contents) > 0 {
				if err := validateNodes(node.Contents); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// validateName ensures a node name does not contain path traversals or separators.
func validateName(name string) error {
	if name == "" {
		return errors.New("empty entry name is not allowed")
	}
	if name == "." || name == ".." {
		return fmt.Errorf("invalid entry name %q: relative references not allowed", name)
	}
	if strings.ContainsAny(name, `/\`) {
		return fmt.Errorf("invalid entry name %q: path separators not allowed", name)
	}
	cleaned := filepath.Clean(name)
	if cleaned != name {
		return fmt.Errorf("invalid entry name %q", name)
	}
	return nil
}
