[![Circle CI](https://circleci.com/gh/segmentio/golines.svg?style=svg&circle-token=b1d01d8b035ef0aa71ccd183580586a80cd85271)](https://circleci.com/gh/segmentio/golines)
# golines
Golines is a golang formatter that shortens long lines, in addition to all
of the formatting fixes done by [`gofmt`](https://golang.org/cmd/gofmt/).

## Motivation

The standard golang formatting tools (`gofmt`, `goimports`, etc.) are great, but
[deliberately don't shorten long lines](https://github.com/golang/go/issues/11915); instead, this
is an activity left to developers.

While there are different tastes when it comes to line lengths in go, we've generally found
that very long lines are more difficult to read than their shortened alternatives. As an example:

```go
func MyFunction(myFirstArgument string, mySecondArgument string, myThirdArgument string, myFourthArgument string, myFifthArgument string) (string, error) {
  ...
}
```

vs.

```go
func MyFunction(
    myFirstArgument string,
    mySecondArgument string,
    myThirdArgument string,
    myFourthArgument string,
    myFifthArgument string,
) (string, error) {
  ...
}
```

We built `golines` to give go developers the option to automatically shorten long lines, like
the one above, according to their preferences.

## Examples

See this [before](_fixtures/end_to_end.go) and [after](_fixtures/end_to_end__exp.go)
view a file with very long lines. More example pairs can be found in the
[`_fixtures`](_fixtures) directory.

## Usage

First, install the tool:

```
go get -u github.com/segmentio/golines
```

Then, run:

```
golines [paths to format]
```

The paths can be either directories or individual files. If no paths are
provided, then input is taken from stdin (as with `gofmt`).

By default, the results are printed to stdout. To overwrite the existing
files in place, use the `-w` flag.

## Options

Some other options are described in the sections below. Run `golines --help` to see
all available flags and settings.

#### Line length settings

By default, the tool tries to shorten lines that are longer than 100 columns
and assumes that 1 tab = 4 columns. The latter can be changed via the
`-m` and `-t` flags respectively.

#### Dry-run mode

Running the tool with the `--dry-run` flag will show pretty, git-style diffs
via an embedded Python script.

#### Comment shortening

Shortening long comment lines is harder than shortening code because comments can
have arbitrary structure and format. `golines` includes some basic
logic for shortening single-line (i.e., `//`-prefixed) comments, but this is turned
off by default since the quality isn't great. To enable this feature anyway, run
with the `--shorten-comments` flag.

#### Custom formatters

By default, the tool will use [`goimports`](https://godoc.org/golang.org/x/tools/cmd/goimports) as
the base formatter (if found), otherwise it will revert to `gofmt`. An explicit formatter can be
set via the `--base-formatter` flag.

#### Generated files

By default, the tool will not format any files that look like they're generated. If you
want to reformat these too, run with the `--no-ignore-generated` flag.

#### Struct tag reformatting

In addition to shortening long lines, the tool also aligns struct tag keys; see the
associated [before](_fixtures/struct_tags.go) and [after](_fixtures/struct_tags__exp.go)
examples in the `_fixtures` directory. To turn this behavior off, run with `--no-reformat-tags`.

## Developer Tooling Integration

### vim-go

Add the following lines to your vimrc, substituting `128` with your preferred line length:

```vim
let g:go_fmt_command = "golines"
let g:go_fmt_options = {
    \ 'golines': '-m 128',
    \ }
```

### Others

Coming soon.

## How It Works

For each input source file, `golines` runs through the following process:

1. Read the file, break it into lines
2. Add a specially-formatted annotation (comment) to each line that's longer
  than the configured maximum (by default, 100 columns)
3. Use [Dave Brophy's](https://github.com/dave) excellent
  [decorated syntax tree](https://github.com/dave/dst) library to parse the code
  plus added annotations
4. Do a depth-first traversal of the resulting tree, looking for nodes
  that have an annotation on them
5. If a node is part of a line that's too long, shorten it by altering
  the newlines around the node and/or its children
6. Repeat steps 2-5 until no more shortening can be done
7. Run the base formatter (e.g., `gofmt`) over the results, write these to either
  stdout or the source file

## Limitations

The tool has been tested on a variety of inputs, but it's not perfect. Among
other examples, the handling of long lines in comments could be improved. If you see
anything particularly egregious, please report via an issue.
