// md-toc-go generates a Markdown file tree (and, later, tables of contents)
// for git repositories. All behaviour lives in cmd; this file only starts
// the CLI.
package main

import (
	"runtime/debug"

	"github.com/EmilyBurak/md-toc-go/cmd"
)

var version = "dev"

func main() {
	if version == "dev" {
		if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "(devel)" && info.Main.Version != "" {
			version = info.Main.Version
		}
	}
	cmd.SetVersion(version)
	cmd.Execute()
}
