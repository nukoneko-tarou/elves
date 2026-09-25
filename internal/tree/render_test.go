package tree

import (
	"bytes"
	"testing"
)

func TestRender(t *testing.T) {
	tests := []struct {
		name     string
		rootName string
		nodes    []Node
		opts     RenderOptions
		expected string
	}{
		{
			name:     "Flat directories",
			rootName: ".",
			nodes: []Node{
				{Type: NodeTypeDirectory, Name: "dir1"},
				{Type: NodeTypeDirectory, Name: "dir2"},
			},
			opts:     RenderOptions{Gitkeep: false},
			expected: ".\n├── dir1\n└── dir2\n",
		},
		{
			name:     "Nested directories with subdirectory root",
			rootName: "my-project",
			nodes: []Node{
				{
					Type: NodeTypeDirectory,
					Name: "src",
					Contents: []Node{
						{Type: NodeTypeDirectory, Name: "pkg"},
					},
				},
				{Type: NodeTypeDirectory, Name: "docs"},
			},
			opts: RenderOptions{Gitkeep: false},
			expected: "my-project\n" +
				"├── src\n" +
				"│   └── pkg\n" +
				"└── docs\n",
		},
		{
			name:     "With gitkeep flag",
			rootName: ".",
			nodes: []Node{
				{Type: NodeTypeDirectory, Name: "empty-dir"},
			},
			opts: RenderOptions{Gitkeep: true},
			expected: ".\n" +
				"└── empty-dir\n" +
				"    └── .gitkeep\n",
		},
		{
			name:     "Ignore file nodes in tree render",
			rootName: ".",
			nodes: []Node{
				{Type: NodeTypeFile, Name: "README.md"},
				{Type: NodeTypeDirectory, Name: "cmd"},
			},
			opts:     RenderOptions{Gitkeep: false},
			expected: ".\n└── cmd\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var buf bytes.Buffer
			Render(&buf, tt.rootName, tt.nodes, tt.opts)
			if buf.String() != tt.expected {
				t.Errorf("expected:\n%s\ngot:\n%s", tt.expected, buf.String())
			}
		})
	}
}
