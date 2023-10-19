package fixtures

import (
	"fmt"
)

func testBinaryOperators() {
	z := argument1 + argument2 + fmt.Sprintf("This is a really long statement that should be broken up %s %s %s", argument1, argument2, argument3)
	y := "hello this is a big string" || "this is a small string" || "the smallest string" || "this is another big string" || "this is an even bigger string >>>"
}
