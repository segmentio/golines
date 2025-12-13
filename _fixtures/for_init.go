package fixtures

func getValue(a, b, c, d, e string) int {
	return 0
}

func testForInit() {
	// For statement with long init clause
	for i := getValue("argument1", "argument2", "argument3", "argument4", "argument5"); i < 10; i++ {
		println(i)
	}

	// For with long condition
	for j := 0; j < getValue("value1", "value2", "value3", "value4", "value5"); j++ {
		println(j)
	}
}

