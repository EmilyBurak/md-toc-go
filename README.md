# md-toc-go

<!-- md-toc-go:tree:start -->
- docs/
  - adr/
    - [Git decides which files are Notes](docs/adr/0001-git-decides-inclusion.md)
  - plans/
    - [Table of Contents design](docs/plans/2026-09-25-toc-design.md)
- [md-toc-go domain glossary](CONTEXT.md)
<!-- md-toc-go:tree:end -->

A Go project with opinionated defaults(see [docs/adr](docs/adr/) for design decisions) to generate Markdown tables of contents and file trees for git repositories.

Inspired by [doctoc](https://github.com/thlorenz/doctoc), [tre](https://github.com/dduan/tre) and [markdown-notes-tree](https://github.com/mistermicheels/markdown-notes-tree)

## Install

Requires Go 1.25 or newer.

```sh
go install github.com/EmilyBurak/md-toc-go@latest
```

This puts `md-toc-go` in `$(go env GOPATH)/bin` (usually `~/go/bin`). Make sure that directory is on your `PATH`.

To build from a local clone instead:

```sh
git clone https://github.com/EmilyBurak/md-toc-go.git
cd md-toc-go
go install .
```

## Usage

Run inside a git repository. Only Markdown files tracked by git are included.

```sh
md-toc-go tree                 # write the tree into ./README.md
md-toc-go tree docs            # scan docs/, write into docs/README.md
md-toc-go tree --target INDEX.md
md-toc-go tree --dry-run       # print the result instead of writing it
md-toc-go tree --check         # don't write, exit 1 if target is out of date
md-toc-go --version
```

The tree is written between these markers, which the tool owns and replaces on every run:

```markdown
<!-- md-toc-go:tree:start -->
<!-- md-toc-go:tree:end -->
```

If the markers are missing they are inserted after the file's leading H1, or at the top when there is none. The block at the top of this README is generated this way.

## Development

This project is a learning exercise in test-driven development with an AI pair. [Claude Code](https://claude.com/claude-code) is used as a guide rather than a code generator: it walks through each red-green cycle and reviews the result, while the source and tests are written by hand. The few exceptions are documentation edits such as the Install and Usage sections above.

Design happens before code. Each feature starts with a grilling session that settles the open questions, and the answers are recorded before any test is written:

- [CONTEXT.md](CONTEXT.md): the domain glossary. Code and docs use these terms and avoid the listed alternatives.
- [docs/adr](docs/adr/): architecture decision records for choices that are hard to reverse.
- [docs/plans](docs/plans/): agreed designs for features not yet built.

These artifacts and the workflow that produces them come from [Matt Pocock's engineering skills](https://github.com/mattpocock/skills) for Claude Code.
