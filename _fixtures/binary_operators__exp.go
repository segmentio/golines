package fixtures

import (
	"fmt"
	"log"
)

func testBinaryOperators() {
	z := argument1 + argument2 + fmt.Sprintf(
		"This is a really long statement that should be broken up %s %s %s",
		argument1,
		argument2,
		argument3,
	)
}
