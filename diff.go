package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/pmezard/go-difflib/difflib"
	"golang.org/x/term"
)

// PrettyDiff prints colored, git-style diffs to the console.
func PrettyDiff(path string, contents []byte, results []byte) error {
	if bytes.Equal(contents, results) {
		return nil
	}

	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(string(contents)),
		B:        difflib.SplitLines(string(results)),
		FromFile: path,
		ToFile:   path + ".shortened",
		Context:  3,
	}

	text, err := difflib.GetUnifiedDiffString(diff)
	if err != nil {
		return err
	}

	ansiGreen := "\033[92m"
	ansiRed := "\033[91m"
	ansiBlue := "\033[94m"
	ansiEnd := "\033[0m"
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimRight(line, " ")
		switch {
		case !term.IsTerminal(int(os.Stdout.Fd())) && len(line) > 0:
			fmt.Printf("%s\n", line)
		case strings.HasPrefix(line, "+"):
			fmt.Printf("%s%s%s\n", ansiGreen, line, ansiEnd)
		case strings.HasPrefix(line, "-"):
			fmt.Printf("%s%s%s\n", ansiRed, line, ansiEnd)
		case strings.HasPrefix(line, "^"):
			fmt.Printf("%s%s%s\n", ansiBlue, line, ansiEnd)
		case len(line) > 0:
			fmt.Printf("%s\n", line)
		}
	}
	fmt.Println("")

	return nil
}
