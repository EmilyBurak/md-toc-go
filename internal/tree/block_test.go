package tree

import "testing"

func wrap(block string) string {
	return startMarker + "\n" + block + endMarker + "\n\n"
}

func TestSplice(t *testing.T) {
	doc := "# Title\n\n<!-- md-toc-go:tree:start -->\nold\n<!-- md-toc-go:tree:end -->\n\ntext\n"
	block := "- [a](a.md)\n"
	got := Splice(doc, block)
	want := "# Title\n\n<!-- md-toc-go:tree:start -->\n" + block + "<!-- md-toc-go:tree:end -->\n\ntext\n"

	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestSpliceMissingMarkers(t *testing.T) {
	doc := "# Title\n\nsome text\n"
	block := "- [a](a.md)\n"
	got := Splice(doc, block)
	want := "# Title\n\n" + wrap(block) + "some text\n"

	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}

}

func TestSpliceJustH1Text(t *testing.T) {
	doc := "# Title\n"
	block := "- [a](a.md)\n"
	got := Splice(doc, block)
	want := doc + wrap(block)

	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestSpliceNoH1(t *testing.T) {
	doc := "some text\n"
	block := "- [a](a.md)\n"
	got := Splice(doc, block)
	want := wrap(block) + doc
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestSpliceIdempotency(t *testing.T) {
	doc := "# Title\n\nsome text\n"
	block := "- [a](a.md)\n"

	once := Splice(doc, block)

	twice := Splice(once, block)

	if twice != once {
		t.Fatalf("got %q, want %q", twice, once)
	}
}
