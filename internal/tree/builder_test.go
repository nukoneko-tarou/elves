package tree

import (
	"os"
	"path/filepath"
	"runtime"
	"syscall"
	"testing"
)

func TestBuild(t *testing.T) {
	if runtime.GOOS != "windows" {
		originalUmask := syscall.Umask(0)
		defer syscall.Umask(originalUmask)
	}

	nodes := []Node{
		{
			Type: NodeTypeDirectory,
			Name: "app",
			Contents: []Node{
				{Type: NodeTypeDirectory, Name: "controllers"},
			},
		},
		{
			Type: NodeTypeFile,
			Name: "ignored.txt",
		},
	}

	t.Run("Basic build without subdir", func(t *testing.T) {
		tempDir := t.TempDir()
		opts := BuildOptions{
			BaseDir:    tempDir,
			Permission: 0755,
			Gitkeep:    false,
		}

		if err := Build(nodes, opts); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		appDir := filepath.Join(tempDir, "app")
		if info, err := os.Stat(appDir); err != nil || !info.IsDir() {
			t.Fatalf("expected app directory to exist")
		}

		controllersDir := filepath.Join(appDir, "controllers")
		if info, err := os.Stat(controllersDir); err != nil || !info.IsDir() {
			t.Fatalf("expected controllers directory to exist")
		}

		ignoredFile := filepath.Join(tempDir, "ignored.txt")
		if _, err := os.Stat(ignoredFile); !os.IsNotExist(err) {
			t.Fatalf("expected ignored.txt to not exist")
		}
	})

	t.Run("Build with SubDir and Gitkeep", func(t *testing.T) {
		tempDir := t.TempDir()
		opts := BuildOptions{
			BaseDir:    tempDir,
			SubDir:     "my-project",
			Permission: 0755,
			Gitkeep:    true,
		}

		if err := Build(nodes, opts); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		gitkeepApp := filepath.Join(tempDir, "my-project", "app", ".gitkeep")
		if info, err := os.Stat(gitkeepApp); err != nil || info.IsDir() {
			t.Fatalf("expected .gitkeep in app directory")
		}

		gitkeepControllers := filepath.Join(tempDir, "my-project", "app", "controllers", ".gitkeep")
		if info, err := os.Stat(gitkeepControllers); err != nil || info.IsDir() {
			t.Fatalf("expected .gitkeep in controllers directory")
		}
	})

	t.Run("Custom permission", func(t *testing.T) {
		if runtime.GOOS == "windows" {
			t.Skip("skip permission check on windows")
		}
		tempDir := t.TempDir()
		opts := BuildOptions{
			BaseDir:    tempDir,
			Permission: 0777,
			Gitkeep:    false,
		}

		if err := Build(nodes, opts); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		appDir := filepath.Join(tempDir, "app")
		info, err := os.Stat(appDir)
		if err != nil {
			t.Fatalf("stat failed: %v", err)
		}
		if info.Mode().Perm() != 0777 {
			t.Errorf("expected permission 0777, got %v", info.Mode().Perm())
		}
	})
}
