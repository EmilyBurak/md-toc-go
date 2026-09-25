package tree

import "strings"

const (
	startMarker = "<!-- md-toc-go:tree:start -->"
	endMarker   = "<!-- md-toc-go:tree:end -->"
)

// Splice returns doc with block placed in the Block, the region between
// startMarker and endMarker that md-toc-go owns. Text outside the Block is
// never touched.
//
// If the markers are already present, the old Block contents are replaced
// with block. If they are absent, a new Block is created at insertPosition
// followed by a blank line. Splice is idempotent: running it again with the
// same block yields the same doc, which is what lets the CLI be re-run
// safely on every commit.
//
// block is expected to be newline-terminated (Render guarantees this), so
// no separator is added between block and endMarker.
func Splice(doc, block string) string {
	startIndex := strings.Index(doc, startMarker)
	if startIndex == -1 {
		// No Block yet: wrap block in fresh markers and insert it. The extra
		// "\n" after wrapped is the blank line that separates the Block from
		// whatever body text follows.
		pos := insertPosition(doc)
		wrapped := startMarker + "\n" + block + endMarker + "\n"
		return doc[:pos] + wrapped + "\n" + doc[pos:]
	}
	// Existing Block. keepUntil is the offset just past the start marker's
	// line ("+1" skips its "\n"); everything before it is preserved.
	keepUntil := startIndex + len(startMarker) + 1
	// strings.Index on the sub-slice gives an offset relative to keepUntil,
	// so add keepUntil back to make it an offset into doc. Searching from
	// keepUntil rather than startIndex also guarantees we don't match inside
	// the start marker itself.
	endIndex := keepUntil + strings.Index(doc[keepUntil:], endMarker)
	// Keep the start marker line, drop the stale contents, keep from the end
	// marker onward.
	return doc[:keepUntil] + block + doc[endIndex:]
}

// insertPosition returns the byte offset in doc where a new Block should be
// inserted when no markers exist.
//
// Rules:
//   - A leading H1 is "# " at the very start of doc: no leading blank lines,
//     and the space after "#" is required.
//   - With a leading H1, the offset is just past the heading line and the
//     single blank line after it (if any), so the Block sits between the
//     heading and the body.
//   - If the H1 is the only line and has no trailing newline, the offset is
//     len(doc); the Block is appended directly after it.
//   - No leading H1: 0, so the Block is prepended.
func insertPosition(doc string) int {
	if !strings.HasPrefix(doc, "# ") {
		return 0
	}

	newLineIndex := strings.Index(doc, "\n")
	if newLineIndex == -1 {
		// H1 with no newline at all: nothing follows it.
		return len(doc)
	}

	// Step past the heading's "\n".
	newLineIndex++

	// If the next byte is another "\n", that's the blank line after the
	// heading; step past it too so the Block lands below it. The length
	// guard covers a doc that ends right after the heading line.
	if newLineIndex < len(doc) && doc[newLineIndex] == '\n' {
		newLineIndex++
	}
	return newLineIndex
}
