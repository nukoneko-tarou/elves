package tree

import (
	"fmt"
	"io"
)

// RenderOptions defines rendering configuration.
type RenderOptions struct {
	Gitkeep bool
}

// Render writes the ASCII tree representation of the nodes to the provided writer.
func Render(w io.Writer, rootName string, nodes []Node, opts RenderOptions) {
	fmt.Fprintln(w, rootName)
	renderNodes(w, nodes, "", opts.Gitkeep)
}

func renderNodes(w io.Writer, nodes []Node, prefix string, gitkeep bool) {
	// Filter for directory nodes to match the directory-creation behavior
	var dirNodes []Node
	for _, n := range nodes {
		if n.Type == NodeTypeDirectory {
			dirNodes = append(dirNodes, n)
		}
	}

	for i, v := range dirNodes {
		isLast := i == len(dirNodes)-1
		connector := "├── "
		newPrefix := prefix + "│   "
		if isLast {
			connector = "└── "
			newPrefix = prefix + "    "
		}

		fmt.Fprintf(w, "%s%s%s\n", prefix, connector, v.Name)

		if gitkeep {
			gitkeepConnector := "├── "
			if len(v.Contents) == 0 {
				gitkeepConnector = "└── "
			}
			fmt.Fprintf(w, "%s%s.gitkeep\n", newPrefix, gitkeepConnector)
		}

		if len(v.Contents) > 0 {
			renderNodes(w, v.Contents, newPrefix, gitkeep)
		}
	}
}
