package commands

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"syscall"
	"testing"
)

func TestCreateCommand(t *testing.T) {
	// Test cases
	testCases := []struct {
		name          string
		args          []string
		jsonContent   string
		expectError   bool
		expectedPaths []string
		permission    os.FileMode
	}{
		{
			name: "Basic directory creation",
			args: []string{"test.json"},
			jsonContent: `[
				{
					"type": "directory",
					"name": "root",
					"contents": [
						{ "type": "directory", "name": "dir1" },
						{ "type": "directory", "name": "dir2" }
					]
				}
			]`,
			expectError: false,
			expectedPaths: []string{
				"dir1",
				"dir2",
			},
			permission: 0755,
		},
		{
			name: "With --sub flag",
			args: []string{"test.json", "--sub", "new-project"},
			jsonContent: `[
				{
					"type": "directory",
					"name": "root",
					"contents": [
						{ "type": "directory", "name": "sub-dir" }
					]
				}
			]`,
			expectError: false,
			expectedPaths: []string{
				"new-project",
				"new-project/sub-dir",
			},
			permission: 0755,
		},
		{
			name: "With --gitkeep flag",
			args: []string{"test.json", "--gitkeep"},
			jsonContent: `[
				{
					"type": "directory",
					"name": "root",
					"contents": [
						{ "type": "directory", "name": "empty-dir" }
					]
				}
			]`,
			expectError: false,
			expectedPaths: []string{
				"empty-dir",
				"empty-dir/.gitkeep",
			},
			permission: 0755,
		},
		{
			name: "With --permission flag",
			args: []string{"test.json", "--permission", "777"},
			jsonContent: `[
				{
					"type": "directory",
					"name": "root",
					"contents": [
						{ "type": "directory", "name": "perm-dir" }
					]
				}
			]`,
			expectError: false,
			expectedPaths: []string{
				"perm-dir",
			},
			permission: 0777,
		},
		{
			name:        "Invalid JSON",
			args:        []string{"test.json"},
			jsonContent: `invalid json`,
			expectError: true,
		},
		{
			name:        "File not found",
			args:        []string{"nonexistent.json"},
			jsonContent: ``,
			expectError: true,
		},
	}

	// Add a new test case for the dry-run feature
	testCases = append(testCases, struct {
		name          string
		args          []string
		jsonContent   string
		expectError   bool
		expectedPaths []string
		permission    os.FileMode
	}{
		name: "Dry run",
		args: []string{"test.json", "--dry-run"},
		jsonContent: `[
			{
				"type": "directory",
				"name": "root",
				"contents": [
					{ "type": "directory", "name": "dir1" },
					{ "type": "directory", "name": "dir2" }
				]
			}
		]`,
		expectError:   false,
		expectedPaths: []string{}, // No paths should be created
	})

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if runtime.GOOS != "windows" {
				originalUmask := syscall.Umask(0)
				defer syscall.Umask(originalUmask)
			}

			// Create a temporary directory for the test
			tempDir := t.TempDir()
			originalWd, err := os.Getwd()
			if err != nil {
				t.Fatalf("Failed to get current working directory: %v", err)
			}
			if err := os.Chdir(tempDir); err != nil {
				t.Fatalf("Failed to change directory: %v", err)
			}
			defer os.Chdir(originalWd)

			// Create a test json file if content is provided
			if tc.jsonContent != "" {
				if err := os.WriteFile("test.json", []byte(tc.jsonContent), 0644); err != nil {
					t.Fatalf("Failed to create test json file: %v", err)
				}
			}

			// Execute the command
			cmd := NewCreate().Cmd
			cmd.SetOut(new(bytes.Buffer)) // Redirect stdout
			cmd.SetErr(new(bytes.Buffer)) // Redirect stderr
			cmd.SetArgs(tc.args)
			err = cmd.Execute()

			// Check for expected error
			if tc.expectError {
				if err == nil {
					t.Errorf("Expected an error, but got none")
				}
				return // Stop further checks if an error is expected
			}
			if err != nil {
				t.Fatalf("Did not expect an error, but got: %v", err)
			}

			if tc.name == "Dry run" {
				// Check output for dry run
				expectedOutput := ".\n├── dir1\n└── dir2\n"
				out := cmd.OutOrStdout().(*bytes.Buffer).String()
				if out != expectedOutput {
					t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, out)
				}

				// Check that no directories were created
				for _, p := range []string{"dir1", "dir2"} {
					fullPath := filepath.Join(tempDir, p)
					if _, err := os.Stat(fullPath); !os.IsNotExist(err) {
						t.Errorf("Expected path not to exist, but it does: %s", p)
					}
				}
			} else {
				// Verify created paths and permissions
				for _, p := range tc.expectedPaths {
					fullPath := filepath.Join(tempDir, p)
					info, err := os.Stat(fullPath)
					if os.IsNotExist(err) {
						t.Errorf("Expected path to exist, but it doesn't: %s", p)
						continue
					}
					if err != nil {
						t.Errorf("Error stating path %s: %v", p, err)
						continue
					}

					// Check permission only for directories on non-windows systems
					if info.IsDir() && runtime.GOOS != "windows" {
						if info.Mode().Perm() != tc.permission {
							t.Errorf("Expected permission %v for %s, but got %v", tc.permission, p, info.Mode().Perm())
						}
					}
				}
			}

			// Clean up is handled by t.TempDir()
		})
	}
}