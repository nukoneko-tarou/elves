package tree

import (
	"os"
	"path/filepath"
)

// BuildOptions defines configuration for directory tree construction.
type BuildOptions struct {
	BaseDir    string
	SubDir     string
	Permission os.FileMode
	Gitkeep    bool
}

// Build creates the directory tree on the filesystem according to the options.
func Build(nodes []Node, opts BuildOptions) error {
	baseDir := opts.BaseDir
	if baseDir == "" {
		baseDir = "."
	}

	targetDir := baseDir
	if opts.SubDir != "" {
		targetDir = filepath.Join(baseDir, opts.SubDir)
		if err := os.Mkdir(targetDir, opts.Permission); err != nil {
			return err
		}
	}

	return createDirectories(nodes, targetDir, opts.Permission, opts.Gitkeep)
}

func createDirectories(nodes []Node, parentDir string, permission os.FileMode, gitkeep bool) error {
	for _, n := range nodes {
		if n.Type == NodeTypeDirectory {
			path := filepath.Join(parentDir, n.Name)
			if err := os.Mkdir(path, permission); err != nil {
				return err
			}

			if gitkeep {
				f, err := os.Create(filepath.Join(path, ".gitkeep"))
				if err != nil {
					return err
				}
				if err := f.Close(); err != nil {
					return err
				}
			}

			if len(n.Contents) > 0 {
				if err := createDirectories(n.Contents, path, permission, gitkeep); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
