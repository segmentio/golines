package fixtures

import "fmt"

func myFunc(a string, b string, c string, d string, e string)

func TestGoroutines() {
	go myFunc("a really long first argument", "a really long second argument", "a really long third arument", "fourth argument", "fifth argument")

	go func() {
		z1 := fmt.Sprintf("This is a really long statement that should be broken up %s %s %s", argument1, argument2, argument3)
	}()
}
