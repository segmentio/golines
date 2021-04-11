[![Circle CI](https://circleci.com/gh/segmentio/golines.svg?style=svg&circle-token=b1d01d8b035ef0aa71ccd183580586a80cd85271)](https://circleci.com/gh/segmentio/golines)
[![Go Report Card](https://goreportcard.com/badge/github.com/segmentio/golines)](https://goreportcard.com/report/github.com/segmentio/golines)
[![GoDoc](https://godoc.org/github.com/segmentio/golines?status.svg)](https://godoc.org/github.com/segmentio/golines)
[![Coverage](https://img.shields.io/badge/Go%20Coverage-84%25-brightgreen.svg?longCache=true&style=flat)](https://gocover.io/github.com/segmentio/golines?version=1.13.x)

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
myMap := map[string]string{"first key": "first value", "second key": "second value", "third key": "third value", "fourth key": "fourth value", "fifth key": "fifth value"}
```

vs.

```go
myMap := map[string]string{
	"first key": "first value",
	"second key": "second value",
	"third key": "third value",
	"fourth key": "fourth value",
	"fifth key": "fifth value",
}
```

We built `golines` to give go developers the option to automatically shorten long lines, like
the one above, according to their preferences.

More background and technical details are available in
[this blog post](https://yolken.net/blog/cleaner-go-code-golines).

## Examples

See this [before](_fixtures/end_to_end.go) and [after](_fixtures/end_to_end__exp.go)
view of a file with very long lines. More example pairs can be found in the
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

#### Chained method splitting

There are several possible ways to split lines that are part of
[method chains](https://en.wikipedia.org/wiki/Method_chaining). The original
approach taken by `golines` was to split on the args, e.g.:

```go
myObj.Method(
	arg1,
	arg2,
	arg3,
).AnotherMethod(
	arg1,
	arg2,
).AThirdMethod(
	arg1,
	arg2,
)
```

Starting in version 0.3.0, the tool now splits on the dots by default, e.g.:

```go
myObj.Method(arg1, arg2, arg3).
	AnotherMethod(arg1, arg2).
	AThirdMethod(arg1, arg2)
```

The original behavior can be used by running the tool with the `--no-chain-split-dots`
flag.

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

### Visual Studio Code

1. Install the [Run on Save](https://marketplace.visualstudio.com/items?itemName=emeraldwalk.RunOnSave) extension
2. Go into the VSCode settings menu, scroll down to the section for the "Run on Save"
  extension, click the "Edit in settings.json" link
3. Set the `emeraldwalk.runonsave` key as follows (adding other flags to the `golines`
  command as desired):

```
    "emeraldwalk.runonsave": {
        "commands": [
            {
                "match": "\\.go$",
                "cmd": "golines ${file} -w"
            }
        ]
    }
```

4. Save the settings and restart VSCode

### Others

Coming soon.

## How It Works

For each input source file, `golines` runs through the following process:

1. Read the file, break it into lines
2. Add a specially-formatted annotation (comment) to each line that's longer
  than the configured maximum
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

See [this blog post](https://yolken.net/blog/cleaner-go-code-golines) for more technical details.

## Limitations

The tool has been tested on a variety of inputs, but it's not perfect. Among
other examples, the handling of long lines in comments could be improved. If you see
anything particularly egregious, please report via an issue.
