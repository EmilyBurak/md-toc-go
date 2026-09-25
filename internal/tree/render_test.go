package tree

import "testing"

func TestRender(t *testing.T) {
	root := Build([]string{"a.md"}, "README.md", stemTitleHelper)
	got := Render(root)
	want := "- [a](a.md)\n"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestRenderDirIndexNote(t *testing.T) {
	root := Build([]string{"docs/README.md", "docs/guide.md"}, "README.md", stemTitleHelper)
	got := Render(root)
	want := "- [README](docs/README.md)\n  - [guide](docs/guide.md)\n"

	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestRenderPlainDir(t *testing.T) {
	root := Build([]string{"docs/guide.md"}, "README.md", stemTitleHelper)
	got := Render(root)
	want := "- docs/\n  - [guide](docs/guide.md)\n"

	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}
