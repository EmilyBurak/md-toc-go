package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd has no Run of its own; invoking md-toc-go bare prints help.
// SilenceUsage keeps runtime errors (not a git repo, unreadable target)
// from being followed by the full usage text.
var rootCmd = &cobra.Command{
	Use:   "md-toc-go",
	Short: "Generate Markdown file trees for git repositories",
	Long: `md-toc-go writes an index of the Markdown notes in a git repository
into a Markdown file, between markers it owns, so it can be re-run safely.

Only files tracked by git are included; ignored and untracked files never
appear. Run it from inside a repository.`,
	SilenceUsage: true,
}

// Execute runs the CLI. Cobra prints the error itself, so this only sets
// the exit status.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
