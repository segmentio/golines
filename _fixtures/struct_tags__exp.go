package fixtures

type MyStruct struct {
	Field1 string `json:"field1" info:"something"`

	ALongField2 string `json:"long_field2" info:"something else" tag:"a really long tag that extends us beyond 100 chars"`
	Field3      string `json:"field3"      info:"third thing"`
	Field4      string `json:"field3"                            tag:"something"`
	Field5      string `                                         tag:"something else"`
}

type MyStruct2 struct {
	Field1 string

	ALongField2 string
	Field3      string `json:"field3"`
	Field4      string
	Field5      string `json:"something else"`
}
