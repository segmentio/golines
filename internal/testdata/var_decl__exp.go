package fixtures

import "fmt"

var abc, cde = []string{
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
}, []string{}

var myMap = map[string]string{
	"first key":  "first value",
	"second key": "second value",
	"third key":  "third value",
	"fourth key": "fourth value",
	"fifth key":  "fifth value",
}
