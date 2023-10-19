package fixtures

import "fmt"

func DeferStmt() {
	defer func() {
		fmt.Printf(
			"This is a really long statement that should be broken up %s %s %s",
			argument1,
			argument2,
			argument3,
		)
	}()

	defer fmt.Printf(
		"This is a really long statement that should be broken up %s %s %s",
		argument1,
		argument2,
		argument3,
	)
}
