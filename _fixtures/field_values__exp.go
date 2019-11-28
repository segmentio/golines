package fixtures

import "fmt"

type MyStruct2 struct {
	name  string
	value string
}

func TestLongStructFields() {
	s := &MyStruct2{
		name: fmt.Sprintf(
			">>>>>>>>>>>>>>>>>>>>>> %s %s %s",
			"a really long first argument",
			"a really long second argument",
			"a third argument",
		),
		value: "short value",
	}
	fmt.Println(s)

	s2 := &MyStruct2{
		name:  "this is a really long name, I don't think we can split it up",
		value: "this is a really long value",
	}
}
