package fixtures

import "fmt"

func TestLoops() {
	for i := 0; i < 5; i++ {
		z1 := fmt.Sprintf("This is a really long statement that should be broken up %s %s %s", argument1, argument2, argument3)
	}

	for _, _ := range []string{"a", "b", "c"} {
		z2 := fmt.Sprintf("This is a really long statement that should be broken up %s %s %s", argument1, argument2, argument3)
	}

	fmt.Println(z1, z2)
}
