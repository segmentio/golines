package fixtures

import "fmt"

type MyStruct struct {
	Field1 string `json:"field1" info:"something"`

	// Field6 example adapted from https://github.com/segmentio/golines/issues/15.
	ALongField2 string `json:"long_field2" info:"something else"                            tag:"a really long tag that extends us beyond 100 chars"`
	Field3      string `json:"field3"      info:"third thing"`
	Field4      string `json:"field3"                                                       tag:"something"`
	Field5      string `                                                                    tag:"something else"`
	Field6      string `json:"somevalue"   info:"http://username:password@example.com:1234"`
}

type MyStruct2 struct {
	Field1 string

	ALongField2 string
	Field3      string `json:"field3"`
	Field4      string
	Field5      string `json:"something else"`
}

func myfunc() {
	s := 4

	type Struct3 struct {
		Field1 string `json:"field1"            info:"something"`
		Field2 string `json:"field2 long value" info:"third thing"`
	}

	s2 := Struct3{}
	fmt.Println(s, s2)
}
