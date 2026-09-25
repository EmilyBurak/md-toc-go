package tree

import (
	"path"
	"slices"
	"strings"
	"testing"
)

func stemTitleHelper(p string) string {
	isMd := strings.HasSuffix(p, ".md")
	if !isMd {
		return ""
	}
	title := path.Base(p)
	title = strings.TrimSuffix(title, ".md")
	return title
}

func names(nodes []*Node) (nameSlice []string) {
	for _, node := range nodes {
		nameSlice = append(nameSlice, node.name)
	}
	return nameSlice
}

func TestBuild(t *testing.T) {
	paths := []string{"a.md", "b.md"}
	// a := "a.md"
	// b := "b.md"
	got := Build(paths, "README.md", stemTitleHelper)
	gotChildren := got.children
	childrenNames := names(got.children)
	if len(gotChildren) == 0 {
		t.Fatalf("The children list is empty")
	}

	want := slices.Equal(childrenNames, paths) == true && gotChildren[0].directory != true && gotChildren[0].linkPath == "a.md" && gotChildren[0].title == "a"

	if want == false {
		t.Fatalf("the child list is wrong")
	}
}

func TestBuildNestsDirectories(t *testing.T) {
	got := Build([]string{"docs/guide.md"}, "README.md", stemTitleHelper)
	if len(got.children) != 1 {
		t.Fatalf("The children list has a %+v amount of children", len(got.children))
	}
	docs := got.children[0]
	if !docs.directory || docs.linkPath != "" || docs.name != "docs" {
		t.Fatalf("docs' directory is %+v & linkPath is %+q, name is %+v", docs.directory, docs.linkPath, docs.name)
	}
	if len(docs.children) != 1 {
		t.Fatalf("docs has %+v children", len(docs.children))
	}
	if docs.children[0].linkPath != "docs/guide.md" || docs.children[0].name != "guide.md" {
		t.Fatalf("docs' child's linkPath is %+v & name is %+q", docs.children[0].linkPath, docs.children[0].name)
	}
}

func TestBuildIndexNote(t *testing.T) {
	got := Build([]string{"docs/README.md", "docs/guide.md"}, "README.md", stemTitleHelper)
	if len(got.children) != 1 {
		t.Fatalf("root should have 1 child, has %+v", len(got.children))
	}
	docs := got.children[0]
	if docs.title != "README" {
		t.Fatalf("title should be 'README', got %+v", docs.title)
	}
	if docs.linkPath != "docs/README.md" {
		t.Fatalf("linkPath should be docs/README.md, is %+v", docs.linkPath)
	}
	if len(docs.children) != 1 {
		t.Fatalf("docs should have 1 child, has %+v", len(docs.children))
	}
	if docs.children[0].name != "guide.md" {
		t.Fatalf("doc's child should be 'guide.md', is %+v", docs.children[0].name)
	}
}

func TestBuildOrdering(t *testing.T) {
	root := Build([]string{"b.md", "Zed/x.md", "a.md", "apple/y.md"}, "README.md", stemTitleHelper)
	if len(root.children) != 4 {
		t.Fatalf("should have 4 children, have %+v", len(root.children))
	}
	got := names(root.children)
	want := []string{"apple", "Zed", "a.md", "b.md"}
	if !slices.Equal(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestChildSorting(t *testing.T) {
	root := Build([]string{"docs/b.md", "docs/a.md"}, "README.md", stemTitleHelper)
	if len(root.children) != 1 {
		t.Fatalf("root should have 1 child, have %v", len(root.children))
	}
	docs := root.children[0]
	got := names(docs.children)
	want := []string{"a.md", "b.md"}
	if !slices.Equal(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}
