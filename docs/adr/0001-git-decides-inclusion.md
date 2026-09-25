---
status: accepted
---

# Git decides which files are Notes

A Tree must exclude vendored, ignored and untracked Markdown (a `node_modules/` directory alone can hold hundreds of `.md` files). Rather than walking the filesystem and re-implementing `.gitignore` semantics, md-toc-go asks `git ls-files` for the tracked, non-ignored files and keeps the `.md` ones. The cost is a hard dependency on the `git` binary and a requirement that the Target live inside a git repository; that requirement is part of the CLI's contract, not a limitation to be patched around.

## Considered Options

- **Walk the filesystem and skip hidden directories.** No git dependency, but vendored Markdown leaks in unless users configure excludes.
- **Walk the filesystem and parse `.gitignore` with a library.** Works without git installed, but adds a dependency and a second, subtly different, implementation of ignore rules.
