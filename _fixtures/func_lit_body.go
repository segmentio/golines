package fixtures

type LongTypeName struct {
	Value int
}

func (l LongTypeName) Compare(other LongTypeName) int {
	return l.Value - other.Value
}

func SortFunc[T any](collection []T, fn func(T, T) int) {}

func testFuncLitBody() {
	// Function literal with long body that should be expanded
	items := []LongTypeName{}
	SortFunc(
		items,
		func(a, b LongTypeName) int { return a.Compare(b) + a.Compare(b) + a.Compare(b) + a.Compare(b) },
	)

	// Single line function literal that fits - should stay as is
	SortFunc(items, func(a, b LongTypeName) int { return a.Compare(b) })
}

