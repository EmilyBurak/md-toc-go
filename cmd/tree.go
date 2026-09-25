package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/EmilyBurak/md-toc-go/internal/notes"
	"github.com/EmilyBurak/md-toc-go/internal/tree"
	"github.com/spf13/cobra"
)

var target string
var dryRun bool
var check bool

var treeCmd = &cobra.Command{
	Use:   "tree [dir]",
	Args:  cobra.MaximumNArgs(1),
	Short: "Write a nested list of the repo's Markdown notes into a file",
	Long: `Build a nested Markdown list of every tracked .md file under dir
(default: the current directory) and write it into a target file
(default: README.md in dir).

The list goes between these markers, which the tool owns:

  <!-- md-toc-go:tree:start -->
  <!-- md-toc-go:tree:end -->

If the markers exist, only the text between them is replaced. If they do
not, the block is inserted after a leading H1, or at the top of the file
when there is none. A missing target file is created. Running the command
twice produces the same file, so it is safe in a pre-commit hook.

Link text is each note's first H1, falling back to its filename. A
directory that contains a README.md links to it instead of listing it.
The target file is never listed in its own tree.

--dry run prints the result to stdout instead of writing it. 
--check does not write and exits 1 if the target is
out of date, for CI purposes. These are mutually exclusive.`,
	Example: `  md-toc-go tree
  md-toc-go tree docs --target docs/INDEX.md
  md-toc-go tree --dry-run
  md-toc-go tree --check`,
	RunE: func(cmd *cobra.Command, args []string) error {
		dir := "."
		dst := target
		if len(args) == 1 {
			dir = args[0]
		}
		if dst == "" {
			dst = filepath.Join(dir, "README.md")
		}
		paths, err := notes.ListMarkdown(dir)
		if err != nil {
			return err
		}
		relTarget, err := filepath.Rel(dir, dst)
		if err != nil {
			return err
		}
		relTarget = filepath.ToSlash(relTarget)
		titleFn := func(p string) string {
			return notes.Title(filepath.Join(dir, p))
		}
		root := tree.Build(paths, relTarget, titleFn)
		block := tree.Render(root)
		contents, err := os.ReadFile(dst)
		if errors.Is(err, os.ErrNotExist) {
			contents = []byte{}
		} else if err != nil {
			return err
		}
		finalDoc := tree.Splice(string(contents), block)
		if check {
			if finalDoc == string(contents) {
				return nil
			}
			return fmt.Errorf("%s is out of date", dst)
		}
		if dryRun {
			_, err := os.Stdout.WriteString(finalDoc)
			if err != nil {
				return err
			}
		} else {
			err = os.WriteFile(dst, []byte(finalDoc), 0o644)
			if err != nil {
				return err
			}
		}
		return nil
	},
}

func init() {
	treeCmd.Flags().StringVar(&target, "target", "", "file to write the tree into (default: <dir>/README.md)")
	treeCmd.Flags().BoolVar(&dryRun, "dry-run", false, "print the resulting file to stdout instead of writing it")
	treeCmd.Flags().BoolVar(&check, "check", false, "exit 1 if the target is out of date, without writing")
	treeCmd.MarkFlagsMutuallyExclusive("check", "dry-run")
	rootCmd.AddCommand(treeCmd)
}
