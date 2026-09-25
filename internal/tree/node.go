package tree

import (
	"cmp"
	"slices"
	"strings"
)

// The Index Note from CONTEXT.md
const indexNoteName = "README.md"

// Node is one entry in the Tree: either a Note (a Markdown file) or a
// directory that contains Notes. The root returned by Build is an unnamed
// directory whose children are the top-level entries.
//
// Fields are unexported on purpose: only this package builds and renders
// trees, so callers never need to reach inside.
type Node struct {
	name      string  // last path segment, e.g. "guide.md" or "docs"
	linkPath  string  // full slash-separated path for the Markdown link; "" for a plain directory
	title     string  // link text from the title mapper; "" for a plain directory
	directory bool    // true for directories, false for Notes
	children  []*Node // nested entries; nil for Notes
}

// byDirsThenName checks if exactly one node is a directory, if so it wins,
// otherwise case-insensitive compare --> case sensitive compare on equality
func byDirsThenName(a, b *Node) int {
	if a.directory == b.directory {
		compareVal := cmp.Compare(strings.ToLower(a.name), strings.ToLower(b.name))
		if compareVal == 0 {
			compareVal = cmp.Compare(a.name, b.name)
		}
		return compareVal
	}
	if a.directory {
		return -1
	}
	return 1
}

// Build turns a flat list of slash-separated relative paths into a nested
// Tree. It is pure: no filesystem or git access. The caller supplies the
// paths (from git) and a mapper that turns a path into its Title.
//
// Algorithm: for each path, split on "/" and walk down from the root,
// finding or creating a directory node for every segment except the last.
// The last segment is the Note itself, attached wherever the walk ended.
// A path with no "/" has an empty directory slice, so the walk is skipped
// and the Note lands directly on the root.
//
// The Target path is skipped so the Tree never links to the file it is
// written into.
func Build(paths []string, targetPath string, pathTitleMapper func(p string) string) *Node {
	rootNode := &Node{}
	for _, p := range paths {
		if p == targetPath {
			continue
		}
		segments := strings.Split(p, "/")
		name := segments[len(segments)-1]
		// cur is a cursor that descends one level per directory segment.
		cur := rootNode
		for _, segment := range segments[:len(segments)-1] {
			cur = cur.childDir(segment)
		}
		// if last segment is README.md & cursor is not on root, it's an Index Note
		if name == indexNoteName && cur != rootNode {
			cur.linkPath = p
			cur.title = pathTitleMapper(p)
			continue
		}
		cur.children = append(cur.children, &Node{name: name, linkPath: p, title: pathTitleMapper(p)})
	}
	rootNode.sort()
	return rootNode
}

// childDir returns the child directory of n called name, creating and
// attaching it if it does not exist yet. The match checks both name and
// directory so a Note and a directory with the same name (notes.md next to
// notes/) stay distinct.
//
// Pointer receiver is required: appending to n.children must mutate the
// real node, not a copy. Returning the new node (not n) is what lets the
// caller's cursor actually descend.
func (n *Node) childDir(name string) *Node {
	for _, child := range n.children {
		if child.name == name && child.directory {
			return child
		}
	}
	newNode := &Node{name: name, directory: true}
	n.children = append(n.children, newNode)
	return newNode
}

func (n *Node) sort() {
	slices.SortFunc(n.children, byDirsThenName)

	for _, c := range n.children {
		// Recurses through children down
		c.sort()
	}
}
