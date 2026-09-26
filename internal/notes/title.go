package notes

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// stem returns the filename without its directory or extension:
// "docs/guide.md" becomes "guide". It is the fallback Title for a Note
// that has no H1.
func stem(path string) string {
	base := filepath.Base(path)
	ext := filepath.Ext(path)
	final := strings.TrimSuffix(base, ext)
	return final
}

// Title returns the link text for the Note at path: the text of its first
// H1 line (a line starting with "# "), trimmed of surrounding whitespace.
// If the file has no H1, or cannot be opened or read, Title falls back to
// stem(path).
//
// Errors are deliberately swallowed into the fallback: one unreadable Note
// should not abort building the whole Tree, and the stem is always a usable
// name. For the same reason scanner.Err is not checked; a read error simply
// ends the scan and yields the stem.
//
// Only the first H1 is considered, so reading stops as soon as one is found
// and the rest of the file is never loaded.
func Title(path string) string {
	file, err := os.Open(path)
	if err != nil {
		return stem(path)
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if rest, ok := strings.CutPrefix(line, "# "); ok {
			return strings.TrimSpace(rest)
		}

	}
	return stem(path)
}
