package fixtures

type Request struct {
	Entity       string
	Fields       []string
	VersionCheck bool
}

func update(ctx string, req Request) error {
	return nil
}

func testCompositeLit() {
	// Struct literal with long field line that should be split
	update(
		"ctx",
		Request{
			Entity:       "enrollment",
			Fields:       []string{"field1", "field2", "field3", "field4"},
			VersionCheck: true,
		},
	)

	// Struct literal with nested long slice
	update(
		"ctx",
		Request{
			Entity:       "test",
			Fields:       []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"},
			VersionCheck: false,
		},
	)
}
