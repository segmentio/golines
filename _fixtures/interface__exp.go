package fixtures

type MyEmptyInterface interface{}

type MyInterface2 interface {
	aReallyLongFunctionName(
		argument1 string,
		argument2 string,
		argument3 string,
		argument4 string,
		argument5 string,
		argument6 string,
	) (string, error)
	shortFunc(x int, y int) error
}
