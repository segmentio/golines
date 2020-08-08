package fixtures

import "fmt"

type MyStruct struct {
	Field1 string `json:"field1" info:"something"`

	// Field6 example adapted from https://github.com/segmentio/golines/issues/15.
	ALongField2 string `json:"long_field2" info:"something else"                            tag:"a really long tag that extends us beyond 100 chars"`
	Field3      string `json:"field3"      info:"third thing"`
	Field4      string `json:"field3"                                                       tag:"something"`
	Field5      int    `                                                                    tag:"something else"`
	Field6      string `json:"somevalue"   info:"http://username:password@example.com:1234"`
}

type MyStruct2 struct {
	Field1 string

	ALongField2 string
	Field3      string            `json:"field3"         info:"here"`
	Field5      map[string]string `json:"something else"`
	MyStruct    `                  json:"mystruct tag"               info2:"here"`
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

type Struct4 struct {
	Field1   []int `json:"field"`
	MyStruct `      json:"field"`
}

type Struct5 struct {
	Field1   *int `json:"field"`
	MyStruct `     json:"field"`
}

type Struct6 struct {
	Field1   chan<- int `json:"field"`
	MyStruct `           json:"field"`
}

type Struct7 struct {
	Field1   <-chan int `json:"field"`
	Field2   string     `json:"field"`
	MyStruct `           json:"field"`
}

type Struct8 struct {
	Field1   chan int `json:"field"`
	MyStruct `         json:"field"`
}

type Struct9 struct {
	Field1   MyStruct `json:"field"`
	MyStruct `         json:"field"`
}
