# md-toc-go

<!-- md-toc-go:tree:start -->
- docs/
  - adr/
    - [Git decides which files are Notes](docs/adr/0001-git-decides-inclusion.md)
- [md-toc-go domain glossary](CONTEXT.md)
<!-- md-toc-go:tree:end -->

A Go project with opinionated defaults to generate Markdown tables of contents and file trees for git repositories.

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
