package notes

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// ListMarkdown returns the slash-separated, repo-relative paths of every
// Note under dir: the Markdown files git tracks there, in git's own order.
// Ignored and untracked files are excluded, which is the whole point of
// asking git rather than walking the filesystem (see
// docs/adr/0001-git-decides-inclusion.md).
//
// dir must be inside a git repository and git must be on PATH; either
// failing is returned as an error, not an empty list. An empty repo yields
// a nil slice and no error. git's own message is included.
//
// The -z flag makes git NUL-separate paths instead of quoting them, so
// names with spaces or non-ASCII bytes come back byte-exact.
func ListMarkdown(dir string) ([]string, error) {
	cmd := exec.Command("git", "ls-files", "-z", "--", "*.md")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		var exitErr *exec.ExitError
		// Output stores git's stderr, surface it for legibility
		matchedErr := errors.As(err, &exitErr)
		if matchedErr && len(exitErr.Stderr) > 0 {
			return nil, fmt.Errorf("git ls-files in %s: %w : %s", dir, err, strings.TrimSpace(string(exitErr.Stderr)))
		}
		return nil, fmt.Errorf("git ls-files in %s: %w", dir, err)
	}
	trimmedOut := strings.TrimSuffix(string(out), "\x00")
	if trimmedOut == "" {
		return nil, nil
	}
	return strings.Split(trimmedOut, "\x00"), nil
}
