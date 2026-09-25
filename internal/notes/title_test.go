package notes

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTitle(t *testing.T) {
	dir := t.TempDir()
	text := []byte("# User Guide\n\ntext\n")
	path := filepath.Join(dir, "guide.md")
	err := os.WriteFile(path, text, 0o644)

	if err != nil {
		t.Fatalf("set up error when writing file")
	}

	got := Title(path)
	want := "User Guide"

	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestTitleNoH1(t *testing.T) {
	dir := t.TempDir()
	text := []byte("User Guide\n\ntext\n")
	path := filepath.Join(dir, "guide.md")
	err := os.WriteFile(path, text, 0o644)

	if err != nil {
		t.Fatalf("set up error %v", err)
	}

	got := Title(path)
	want := "guide"

	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}
