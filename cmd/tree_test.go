package cmd

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func newRepo(t *testing.T) string {
	t.Helper()
	_, err := exec.LookPath("git")
	if err != nil {
		t.Skip("git not installed")
	}
	dir := t.TempDir()
	cmd := exec.Command("git", "init", "-q")
	cmd.Dir = dir
	err = cmd.Run()
	if err != nil {
		t.Fatalf("error running git init: %v", err)
	}
	err = os.WriteFile(filepath.Join(dir, "a.md"), []byte{}, 0o644)
	if err != nil {
		t.Fatalf("error writing a.md: %v", err)
	}

	cmd = exec.Command("git", "add", "-A")
	cmd.Dir = dir
	err = cmd.Run()
	if err != nil {
		t.Fatalf("error running git add: %v", err)
	}
	return dir
}

func run(t *testing.T, args ...string) error {
	// Reset the globals on each run
	target = ""
	dryRun = false
	check = false
	rootCmd.SetArgs(append([]string{"tree"}, args...))
	return rootCmd.Execute()
}

func TestTreeCheckStale(t *testing.T) {
	dir := newRepo(t)
	err := run(t, dir, "--check")
	// Check should fail if stale
	if err == nil {
		t.Fatalf("--check on stale repo should return an error")
	}
	_, err = os.Stat(filepath.Join(dir, "README.md"))
	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("README exists")
	}
}

func TestTreeCheckCurrent(t *testing.T) {
	dir := newRepo(t)
	err := run(t, dir)
	if err != nil {
		t.Fatalf("error running tree %v: %v", dir, err)
	}
	err = run(t, dir, "--check")
	if err != nil {
		t.Fatalf("error running current tree check: %v", err)
	}
}
