package fixtures

import "fmt"

func CaseLists() {
	x := "hello"

	switch x {
	case "a really long string", "another really long string", "a third really long string", "a fourth really long string":
		fmt.Println("hello")
	}
}
