package tree

import (
	"strings"
	"testing"
)

func TestParseJSON(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid tree JSON with directories and files",
			input: `[
				{
					"type": "directory",
					"name": ".",
					"contents": [
						{"type": "file", "name": "README.md"},
						{
							"type": "directory",
							"name": "src",
							"contents": [
								{"type": "file", "name": "main.go"}
							]
						}
					]
				},
				{"type": "report", "directories": 1, "files": 2}
			]`,
			expectError: false,
		},
		{
			name:        "Empty input",
			input:       "",
			expectError: true,
			errorMsg:    "failed to decode json",
		},
		{
			name:        "Invalid JSON syntax",
			input:       "{not a valid json}",
			expectError: true,
			errorMsg:    "failed to decode json",
		},
		{
			name:        "Empty array",
			input:       "[]",
			expectError: true,
			errorMsg:    "invalid or empty json file",
		},
		{
			name: "Empty contents in root",
			input: `[
				{
					"type": "directory",
					"name": ".",
					"contents": []
				}
			]`,
			expectError: true,
			errorMsg:    "invalid or empty json file",
		},
		{
			name: "Path traversal with .. in entry name",
			input: `[
				{
					"type": "directory",
					"name": ".",
					"contents": [
						{"type": "directory", "name": ".."}
					]
				}
			]`,
			expectError: true,
			errorMsg:    "relative references not allowed",
		},
		{
			name: "Path separator in entry name",
			input: `[
				{
					"type": "directory",
					"name": ".",
					"contents": [
						{"type": "directory", "name": "sub/dir"}
					]
				}
			]`,
			expectError: true,
			errorMsg:    "path separators not allowed",
		},
		{
			name: "Windows path separator in entry name",
			input: `[
				{
					"type": "directory",
					"name": ".",
					"contents": [
						{"type": "directory", "name": "sub\\dir"}
					]
				}
			]`,
			expectError: true,
			errorMsg:    "path separators not allowed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tree, err := ParseJSON(strings.NewReader(tt.input))
			if tt.expectError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if tt.errorMsg != "" && !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error message to contain %q, got %q", tt.errorMsg, err.Error())
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tree == nil || len(*tree) == 0 {
				t.Fatal("expected non-empty tree")
			}
		})
	}
}
