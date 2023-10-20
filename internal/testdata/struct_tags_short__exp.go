package fixtures

type MyStruct struct {
	Field1 string `json:"field1" info:"something"`

	ALongField2 string `json:"long_field2" info:"something else"`
	Field3      string `json:"field3"      info:"third thing"`
}
