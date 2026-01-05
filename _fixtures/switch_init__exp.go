package fixtures

func getType(a, b, c, d, e, f string) int {
	return 0
}

func testSwitchInit() {
	// Switch statement with long init clause
	switch v := getType(
		"argument1",
		"argument2",
		"argument3",
		"argument4",
		"argument5",
		"argument6",
	); v {
	case 1:
		println("one")
	case 2:
		println("two")
	default:
		println("other")
	}
}
