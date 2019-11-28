package fixtures

import "fmt"

// Short prefix
// This is a really, really long comment on a single line. We should try to break it up if possible
// because it's longer than 100 chars.
// Short suffix
//

// Another comment

/*
	A block comment. Really long lines in here aren't currently processed because they're a bit harder to handle.

	func(aReallyLongArgument string, anotherReallyLongArgument string, aThirdReallyLongArgument string) (string, error) {
		return "", nil
	}
*/
// go:generate this is a really long go generate line. We don't want to shorten this because that could cause problems running go generate.
func testFunc() {
	for i := 0; i < 10; i++ {
		if i > 5 {
			// This is a another really, really long comment on a single line. We should try to break
			// it up if possible because it's longer than 100 chars.
			fmt.Print("hello")
		}
	}
}
