[![golines test](https://github.com/segmentio/golines/actions/workflows/test.yml/badge.svg)](https://github.com/segmentio/golines/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/segmentio/golines)](https://goreportcard.com/report/github.com/segmentio/golines)
[![GoDoc](https://godoc.org/github.com/segmentio/golines?status.svg)](https://godoc.org/github.com/segmentio/golines)

# golines

Golines is a Go code formatter that shortens long lines, in addition to all
of the formatting fixes done by [`gofmt`](https://golang.org/cmd/gofmt/).

## Maintenance

As of late 2024, [segmentio/golines](https://github.com/segmentio/golines/) has
functionally been in maintenance mode and several dependencies appear to be
similarly unmaintained. At some point in Q4 2025, this repository
[will be archived](https://docs.github.com/en/repositories/archiving-a-github-repository/archiving-repositories)
unless active maintainership can be found within Twilio Segment.
The code will remain available and the terms of the license will not be changed.

## Motivation

The standard Go formatting tools (`gofmt`, `goimports`, etc.) are great, but
[deliberately don't shorten long lines](https://github.com/golang/go/issues/11915);
instead, this is an activity left to developers.

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

We built `golines` to give Go developers the option to automatically shorten long lines, like
the one above, according to their preferences.

More background and technical details are available in
[this blog post](https://yolken.net/blog/cleaner-go-code-golines).

## Examples

See this [before](_fixtures/end_to_end.go) and [after](_fixtures/end_to_end__exp.go)
view of a file with very long lines. More example pairs can be found in the
[`_fixtures`](_fixtures) directory.

## Version support

Since v0.10.0, releases of `golines` have required at least Go 1.18 due to
generics-related dependencies. As of v0.13.0, `golines` requires a minimum of
Go 1.23 due to transitive requirements introduced by dependencies.

Generally, the [minimum version](https://go.dev/ref/mod#go-mod-file-go) in [`go.mod`](./go.mod)
is the absolute minimum required version of Go for any given version of `golines.`

If you need to use `golines` with an older version of go, install the tool from
the `v0.9.x` or `v0.12.x` releases.

## Usage

First, install the tool. If you're using Go 1.21 or newer, run:

```text
go install github.com/segmentio/golines@latest
```

Otherwise, for older Go versions, run:

```text
go install github.com/segmentio/golines@v0.9.0
```

Then, run:

```text
golines [paths to format]
```

The paths can be either directories or individual files. If no paths are
provided, then input is taken from `stdin` (as with `gofmt`).

By default, the results are printed to `stdout`. To overwrite the existing
files in place, use the `-w` flag.

## Options

Some other options are described in the sections below. Run `golines --help` to
see all available flags and settings.

### Line length settings

By default, the tool tries to shorten lines that are longer than 100 columns
and assumes that 1 tab = 4 columns. The latter can be changed via the
`-m` and `-t` flags respectively.

#### Dry-run mode

Running the tool with the `--dry-run` flag will show pretty, git-style diffs.

#### Comment shortening

Shortening long comment lines is harder than shortening code because comments can
have arbitrary structure and format. `golines` includes some basic
logic for shortening single-line (i.e., `//`-prefixed) comments, but this is turned
off by default since the quality isn't great. To enable this feature anyway, run
with the `--shorten-comments` flag.

#### Custom formatters

By default, the tool will use [`goimports`](https://godoc.org/golang.org/x/tools/cmd/goimports)
as the base formatter (if found), otherwise it will revert to `gofmt`. An explicit
formatter can be set via the `--base-formatter` flag; the command provided here
should accept its input via `stdin` and write its output to `stdout`.

#### Generated files

By default, the tool will not format any files that look like they're generated.
If you want to reformat these too, run with the flag `--ignore-generated=false`.

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

The original behavior can be used by running the tool with the
`--no-chain-split-dots` flag.

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
3. Set the `emeraldwalk.runonsave` key as follows
   (adding other flags to the `golines` command as desired):
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

1. Save the settings and restart VSCode

### Goland

1. Go into the Goland settings and click "Tools" -> "File Watchers" then click the plus to create a new file watcher
2. Set the following properties:
   - __Name:__ `golines`
   - __File type:__ `Go files`
   - __Scope:__ `Project Files`
   - __Program:__ `golines`
   - __Arguments:__ `$FilePath$ -w`
   - __Output paths to refresh:__ `$FilePath$`
3. In the "Advanced Options" section uncheck the __Auto-save edited files to trigger the watcher__ setting
4. Confirm by clicking OK
5. Activate your newly created file watcher in the Goland settings under "Tools" -> "Actions on save"

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
  `stdout` or the source file

See [this blog post](https://yolken.net/blog/cleaner-go-code-golines) for more technical details.

## Limitations

The tool has been tested on a variety of inputs, but it's not perfect. Among
other examples, the handling of long lines in comments could be improved. If you see
anything particularly egregious, please report via an issue.
