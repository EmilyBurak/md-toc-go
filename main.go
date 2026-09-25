// md-toc-go generates a Markdown file tree (and, later, tables of contents)
// for git repositories. All behaviour lives in cmd; this file only starts
// the CLI.
package main

import "github.com/EmilyBurak/md-toc-go/cmd"

func main() {
	cmd.Execute()
}
