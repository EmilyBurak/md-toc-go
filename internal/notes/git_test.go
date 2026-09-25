package notes

import (
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"testing"
)

func TestListMarkdown(t *testing.T) {
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

	err = os.MkdirAll(filepath.Join(dir, "docs"), 0o755)
	if err != nil {
		t.Fatalf("error writing /docs: %v", err)
	}
	err = os.WriteFile(filepath.Join(dir, "docs/b.md"), []byte{}, 0o644)
	if err != nil {
		t.Fatalf("error writing docs/b.md: %v", err)
	}

	err = os.WriteFile(filepath.Join(dir, "notes.txt"), []byte{}, 0o644)
	if err != nil {
		t.Fatalf("error writing notes.txt: %v", err)
	}

	cmd = exec.Command("git", "add", "-A")
	cmd.Dir = dir
	err = cmd.Run()
	if err != nil {
		t.Fatalf("error running git add: %v", err)
	}

	got, err := ListMarkdown(dir)
	if err != nil {
		t.Fatalf("ListMarkdown error: %v", err)
	}
	want := []string{"a.md", "docs/b.md"}
	if !slices.Equal(want, got) {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestListMarkdownNotARepo(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	// A bare temp dir with no git init. Per ADR 0001 this is a hard error,
	// not an empty result.
	dir := t.TempDir()

	got, err := ListMarkdown(dir)
	if err == nil {
		t.Fatalf("expected an error outside a git repo, got nil (paths %q)", got)
	}
}

func TestListMarkdownEmptyRepo(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	dir := t.TempDir()
	cmd := exec.Command("git", "init", "-q")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("error running git init: %v", err)
	}

	got, err := ListMarkdown(dir)
	if err != nil {
		t.Fatalf("ListMarkdown error: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("expected no paths in an empty repo, got %q", got)
	}
}
