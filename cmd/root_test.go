package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestExecuteVersion(t *testing.T) {
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	// Restore the default stdout after call
	t.Cleanup(func() { rootCmd.SetOut(nil) })
	rootCmd.SetArgs([]string{"--version"})
	SetVersion("1.2.3")
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "1.2.3") {
		t.Fatalf("got %q, want it to contain 1.2.3", buf.String())
	}
}
