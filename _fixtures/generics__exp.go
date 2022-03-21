package main

func SumIntsOrFloats(
	m map[K]V,
	longArgument1 int,
	longArgument2 int,
	longArgument3 int,
) V {
	var s V
	for _, v := range m {
		s += v
	}
	return s
}

func SumIntsOrFloats2(m map[K]V) V {
	var s V
	for _, v := range m {
		s += v
	}
	return s
}
