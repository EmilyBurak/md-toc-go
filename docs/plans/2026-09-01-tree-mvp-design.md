
# md-toc-go MVP: file tree generation

Archived 2026-09-27. This is the original plan from the first grilling session, kept verbatim. The tree feature described here is implemented; CONTEXT.md and docs/adr hold the decisions that outlived it.

## Context

`md-toc-go` is an empty Go project (README only) meant to generate Markdown tables of contents and file trees for git repositories, in the spirit of doctoc, tre and markdown-notes-tree. The MVP starts with the **tree** feature. A grilling session settled every design branch below; the TOC feature is out of scope but the CLI shape leaves room for it.

## Decisions (from the grilling session)

| Branch | Decision |
|---|---|
| Output shape | Nested Markdown bullet list of relative links, written into a Markdown file |
| Inclusion | Markdown files only (`.md`). The tree is a **notes index**, not a repo map |
| Exclusion | Delegated to `git ls-files` (tracked, non-ignored), filtered to `.md`. Requires git + a repo |
| Target | One target file, full nested tree. Default `README.md` in the scanned dir; `--target` overrides |
| Insertion | Between `<!-- md-toc-go:tree:start -->` / `<!-- md-toc-go:tree:end -->` markers. If markers are absent, insert at top, after a leading H1 if present. Create the file if missing. Idempotent |
| Link text | First `# H1` of the file; fallback filename without extension |
| Directories | Dir containing `README.md` renders as a link to that README, and the README is not listed as a child. Otherwise plain text `dir/` |
| Self-reference | Target file is excluded from the tree. No root node; root's children are top level |
| Ordering | Dirs first, then files; case-insensitive A–Z by path name (not title) |
| CLI | `md-toc-go tree [dir]` using spf13/cobra. Flags: `--target`, `--dry-run` (print to stdout, no write) |
| Seam | Pure core (paths + title lookup → Node tree → Markdown → spliced document). Git and file I/O in thin adapters |
| Module | `github.com/EmilyBurak/md-toc-go`, Go 1.25 |

## Ubiquitous language (goes in CONTEXT.md)

- **Note**: a Markdown file tracked by git.
- **Tree**: the nested index of Notes under a directory, rendered as a Markdown list.
- **Target**: the Markdown file the Tree is written into. Never appears in its own Tree.
- **Index Note**: a `README.md` inside a directory; the directory links to it and it is not listed separately.
- **Title**: link text of a Note: its first H1, else its filename stem.
- **Block**: the region of the Target between the start and end markers that md-toc-go owns.

## Layout

```
go.mod                          github.com/EmilyBurak/md-toc-go
main.go                         calls cmd.Execute()
cmd/root.go                     cobra root command
cmd/tree.go                     `tree` subcommand: wires adapters to core
internal/tree/node.go           Node type + Build(paths []string, title func(path) string) *Node
internal/tree/render.go         Render(*Node) string  (Markdown list)
internal/tree/block.go          Splice(doc, block string) string  (markers / insert-after-H1)
internal/tree/*_test.go         table tests on fixtures, no git
internal/notes/git.go           ListMarkdown(dir) ([]string, error)  via `git ls-files -z -- '*.md'`
internal/notes/title.go         Title(path) string  (first H1 else stem)
internal/notes/git_test.go      integration test: git init in t.TempDir(), t.Skip if git missing
CONTEXT.md                      glossary above
docs/adr/0001-git-decides-inclusion.md   why exclusion is delegated to git (hard to reverse: it is the CLI's contract that a repo is required)
```

## Implementation steps (TDD, one red-green cycle each)

1. `go mod init github.com/EmilyBurak/md-toc-go`; `go get github.com/spf13/cobra@latest` (check current API via Context7 first).
2. **Build**: from `[]string` of slash-separated relative paths, produce a `Node` tree. Drop the target path. Collapse `dir/README.md` into the dir node's link. Sort: dirs first, then files, case-insensitive by name.
3. **Render**: `Node` → Markdown. Two-space indentation per level. File: `- [Title](rel/path.md)`. Dir with Index Note: `- [Title](dir/README.md)`. Dir without: `- dir/`. Links are relative to the Target's directory.
4. **Splice**: replace Block if markers exist; else insert after first line matching `^# ` (plus blank line), else prepend. Preserve trailing newline. Same input twice yields same output.
5. **Adapters**: `git ls-files -z -- '*.md'` run in `dir`, error clearly if not a repo or git missing. `Title` reads until first `# ` line.
6. **CLI**: `tree [dir]`, default `.`; `--target` default `<dir>/README.md`; `--dry-run` prints the spliced document to stdout.
7. Write `CONTEXT.md` and the ADR.

## Learning-mode contribution points

Two spots where the user's choice shapes behaviour, prepared as TODO stubs for them to fill:
- `internal/tree/block.go` `insertPosition(doc string) int`: the "after leading H1" heuristic (what counts as a leading H1, blank-line handling).
- `internal/tree/node.go` sort comparator: the dirs-first, case-insensitive rule.

## Verification

- `go test ./...` passes; core tests need no git.
- `go run . tree --dry-run` on this repo prints a block containing nothing (only README.md exists and it is the target) — confirms self-exclusion.
- Create a temp repo with `docs/README.md`, `docs/Guide.md` (H1 "User Guide"), `notes/a.md` (no H1), run `go run . tree`, and confirm README.md gains:
  ```
  <!-- md-toc-go:tree:start -->
  - [Docs title](docs/README.md)
    - [User Guide](docs/Guide.md)
  - notes/
    - [a](notes/a.md)
  <!-- md-toc-go:tree:end -->
  ```
  after its H1. Run again and confirm the file is byte-identical.
