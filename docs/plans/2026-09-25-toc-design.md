# Table of Contents design

Date: 2026-09-25. Status: agreed, not yet implemented.

## Goal

A second subcommand, `md-toc-go toc`, that writes a nested list of a Note's
headings into that same Note, between markers the tool owns. It must be
idempotent and usable from the pre-commit hook, and it must be able to
update every applicable Note in a repository in one pass.

## Decisions

| Question | Decision |
|---|---|
| Which headings | ATX `##` to `######`. The H1 is the Title, not a Heading. Headings inside fenced code are ignored. Setext headings are not supported. |
| Anchors | GitHub's slug rule. A `--mode` flag exists with `github` as the only value, so other hosts can be added without changing the interface. |
| Markers | Own pair, `<!-- md-toc-go:toc:start -->` / `<!-- md-toc-go:toc:end -->`. Inserted after the leading H1 when absent, same rule as the Tree. A file may hold both a Tree Block and a Contents Block. |
| Target | In place only. The TOC is read from and written into the same file. No `--target`. |
| Nesting | Relative depth. Each Heading is one level deeper than the nearest shallower Heading before it, so H2 → H4 → H3 renders at depths 0, 1, 1. |
| Repo-wide pass | `--all [dir]` updates every Applicable Note: a Note that already contains the Contents markers. Opt-in by marker, so the first run never injects a TOC into files nobody asked about. |

## Language (add to CONTEXT.md)

- **Heading**: an ATX heading of level 2 to 6 in a Note, outside fenced code. The H1 is the Note's Title, not a Heading.
- **Anchor**: the URL fragment a Heading resolves to on GitHub, derived from the Heading text by the GitHub slug rule, with `-1`, `-2` suffixes for repeats.
- **Contents**: the nested list of Anchors for every Heading in a Note, written into a Block of that same Note.
- **Applicable Note**: a Note that already carries the Contents markers. Only these are touched by `--all`.

## ADR

`docs/adr/0002-github-anchors.md`: anchors follow GitHub's rule. Hard to
reverse because generated links break silently on a host with a different
rule. `--mode` is the committed escape hatch.

## CLI

```
md-toc-go toc [file]        one Note, default README.md; inserts markers if missing
md-toc-go toc --all [dir]   every Applicable Note under dir, default .
```

Flags: `--all`, `--dry-run`, `--check`, `--mode github`. `--check` and
`--dry-run` are mutually exclusive, as in `tree`.

| Mode | Behaviour |
|---|---|
| write | Rewrite each file whose content changed. Unchanged files are not written. |
| `--dry-run` | Print each spliced document to stdout. With `--all`, precede each with `==> path <==`. |
| `--check` | Collect every stale path, print one per line to stderr, exit 1 if any. Do not stop at the first. |
| error | Unknown `--mode`, unreadable file, missing file in single-file form, or not a repo with `--all`: stop, exit 1. |

`--all` lists Notes with `notes.ListMarkdown`, so ADR 0001 applies unchanged.

## Data flow (per file, pure)

```
doc string
  → tree.Strip(doc, tocMarkers)          remove the old Block before parsing
  → toc.Parse(stripped) []Heading        {Level int; Text string}, ATX 2–6, fence-aware
  → toc.Build(headings, slugger) []Entry {Text, Anchor string; Depth int}, relative depth
  → toc.Render(entries) string           "- [Text](#anchor)\n", two spaces per depth
  → tree.Splice(doc, tocMarkers, block)  same as tree
```

`tree.Node` has unexported fields and models directories, so the TOC does
not reuse it. Headings with relative depth are a flat sequence; a slice of
entries and a small renderer are simpler than a tree.

### Slugger

```go
type Slugger interface { Slug(text string) string }
func New(mode string) (Slugger, error)   // error for unknown mode
```

Stateful: one instance per document, because duplicate suffixes depend on
what was already handed out. GitHub rule: strip inline markup to plain
text, lowercase, drop every character that is not a letter, digit, space,
hyphen or underscore, replace spaces with hyphens, then append `-N` for
the Nth repeat. A heading that already ends in `-1` and collides with a
generated suffix keeps counting, as GitHub does.

## Layout

```
internal/toc/heading.go      Parse
internal/toc/anchor.go       Slugger, New, githubSlugger
internal/toc/build.go        Build (relative depth)
internal/toc/render.go       Render
internal/toc/*_test.go       table tests, no I/O
internal/tree/block.go       Splice and new Strip take a Markers{Start, End} value
cmd/toc.go                   toc subcommand
cmd/toc_test.go              reuse newRepo; add a toc run helper
```

`Splice` currently hard-codes the tree markers. Threading a `Markers`
parameter through is the first refactor and changes no behaviour.

## Edge cases (each a table-test row)

- Heading inside a ``` or ~~~ fence, including fences with info strings and an unclosed fence: skipped.
- `#hashtag` with no space, `####### seven`: not headings.
- Trailing closing hashes `## Title ##`: stripped.
- Up to three leading spaces: still a heading. Four: indented code, not a heading.
- Inline markup (`` `code` ``, `**bold**`, `[link](url)`): kept as text for display, stripped for the slug.
- Duplicates: `-1`, `-2`; collisions with an existing `-1` suffix keep counting.
- Headings inside the old Contents Block: removed by `Strip` before parsing.
- No headings: empty Block, not an error.
- Both Blocks inserted fresh into one file: the second inserted lands above the first. Stable after that. Documented, not fixed.

## Testing

- `internal/toc`: pure table tests for Parse, Slug, Build, Render.
- `internal/tree/block_test.go`: parameterised markers, `Strip`, both Blocks coexisting.
- `cmd/toc_test.go`: single file, `--all` opt-in by marker, `--check` listing every stale file, run twice is byte-identical.

## Implementation order

1. CONTEXT.md terms and ADR 0002.
2. `Markers` parameter on `Splice`, plus `Strip`. Existing tree tests stay green.
3. Parse, Slugger, Build, Render, each red-green.
4. `toc` subcommand, single file.
5. `--all`, `--check`, `--dry-run`.
6. README Usage, pre-commit hook entry `md-toc-toc` alongside `md-toc-tree`.
