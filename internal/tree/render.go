package tree

import (
	"fmt"
	"strings"
)

// renderInto writes n's children to b as Markdown list items, one per line,
// then recurses into each child one level deeper. depth controls the
// indent: two spaces per level, which is what Markdown needs for a nested
// bullet.
//
// A directory with no Index Note has no link target, so it prints as plain
// text "name/". Everything else (Notes, and directories whose Index Note
// gave them a linkPath and title) prints as "[title](linkPath)".
func renderInto(b *strings.Builder, n *Node, depth int) {
	indent := strings.Repeat("  ", depth)
	for _, c := range n.children {
		if c.directory && c.linkPath == "" {
			fmt.Fprintf(b, "%s- %s/\n", indent, c.name)
		} else {
			fmt.Fprintf(b, "%s- [%s](%s)\n", indent, c.title, c.linkPath)
		}
		renderInto(b, c, depth+1)
	}
}

// Render turns a Tree into the Markdown list that goes inside the Block.
// The root node itself is never printed; its children are the top-level
// items. Each line ends with "\n", so the result is empty for a Tree with
// no Notes and always newline-terminated otherwise.
//
// Links are emitted exactly as stored in linkPath, so the caller is
// responsible for having built the Tree relative to the Target's directory.
func Render(root *Node) string {
	builder := strings.Builder{}
	renderInto(&builder, root, 0)
	return builder.String()
}
