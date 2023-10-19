package fixtures

import "fmt"

func testFunc3() {
	if "hello this is a big string" == "this is a small string" &&
		"this is another big string" == "this is an even bigger string >>>" {
		fmt.Print("inside if statement")
	}
}
