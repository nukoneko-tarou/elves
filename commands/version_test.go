package commands

import (
	"bytes"
	"fmt"
	"testing"
)

func TestVersionCommand(t *testing.T) {
	cmd := NewVersion().Cmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error executing version command: %v", err)
	}

	expected := fmt.Sprintf("elves version %s\n", Version)
	if buf.String() != expected {
		t.Errorf("expected %q, got %q", expected, buf.String())
	}
}
