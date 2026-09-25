package tree

import (
	"fmt"
	"testing"
)

// generateBenchTree generates a tree with depth and breadth.
func generateBenchTree(depth, breadth int) []Node {
	if depth <= 0 {
		return nil
	}
	var nodes []Node
	for i := 0; i < breadth; i++ {
		node := Node{
			Type: NodeTypeDirectory,
			Name: fmt.Sprintf("dir_%d_%d", depth, i),
		}
		if depth > 1 {
			node.Contents = generateBenchTree(depth-1, breadth)
		}
		nodes = append(nodes, node)
	}
	return nodes
}

func BenchmarkBuildNoGitkeep(b *testing.B) {
	tree := generateBenchTree(4, 4)

	concurrencies := []int{1, 2, 4}
	for _, c := range concurrencies {
		b.Run(fmt.Sprintf("concurrency=%d", c), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				tempDir := b.TempDir()
				opts := BuildOptions{
					BaseDir:     tempDir,
					Permission:  0755,
					Gitkeep:     false,
					Concurrency: c,
				}
				if err := Build(tree, opts); err != nil {
					b.Fatalf("build failed: %v", err)
				}
			}
		})
	}
}

func BenchmarkBuildWithGitkeep(b *testing.B) {
	tree := generateBenchTree(4, 4)

	concurrencies := []int{1, 2, 4}
	for _, c := range concurrencies {
		b.Run(fmt.Sprintf("concurrency=%d", c), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				tempDir := b.TempDir()
				opts := BuildOptions{
					BaseDir:     tempDir,
					Permission:  0755,
					Gitkeep:     true,
					Concurrency: c,
				}
				if err := Build(tree, opts); err != nil {
					b.Fatalf("build failed: %v", err)
				}
			}
		})
	}
}
