package fixtures

func doSomething(a, b, c, d, e string) error {
	return nil
}

func testIfInit() {
	// If statement with long init clause
	if err := doSomething(
		"argument1",
		"argument2",
		"argument3",
		"argument4",
		"argument5",
	); err != nil {
		println(err)
	}

	// Nested if with long init
	if x := doSomething("value1", "value2", "value3", "value4", "value5"); x == nil {
		if y := doSomething("nested1", "nested2", "nested3", "nested4", "nested5"); y == nil {
			println("both nil")
		}
	}
}
