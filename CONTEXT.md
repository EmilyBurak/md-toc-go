# md-toc-go domain glossary

A command-line tool that writes navigational Markdown (a notes index today, a table of contents later) into Markdown files inside a git repository.

## Language

**Note**:
A Markdown file tracked by git. Only Notes appear in a Tree.
_Avoid_: Document, page, file

**Tree**:
The nested index of every Note beneath a directory, rendered as a Markdown bullet list of links.
_Avoid_: File tree, listing, sitemap

**Target**:
The Markdown file a Tree is written into. A Target never appears in its own Tree.
_Avoid_: Output file, destination, README

**Index Note**:
A `README.md` inside a directory. The directory links to it and it is not listed as a separate child.
_Avoid_: Directory README, index file

**Title**:
The link text shown for a Note: its first H1 heading, or its filename without extension when it has none.
_Avoid_: Name, label, heading

**Block**:
The region of a Target, delimited by start and end marker comments, that md-toc-go owns and rewrites on every run.
_Avoid_: Section, region, generated area
