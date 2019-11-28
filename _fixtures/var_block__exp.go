package fixtures

import "fmt"

var (
	abcd = []string{
		"a really long string",
		"another really long string",
		"a third really long string",
		"a fourth really long string",
		fmt.Sprintf(
			"%s %s %s %s >>>>> %s %s",
			"first argument",
			"second argument",
			"third argument",
			"fourth argument",
			"fifth argument",
			"sixth argument",
		),
	}
)
