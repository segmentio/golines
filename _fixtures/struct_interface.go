package fixtures

type MyStruct struct {
	i             Int
	LongInterface interface {
		SuperLongFunctionName(anArgument string, anotherReallyLongArgument string, superDuperLongArgument string, definitelyTheLongestArgument string) error
	}
}

type MyShortStruct struct {
	s              string
	ShortInterface interface {
		SuperShortFunctionName(shortArg string) error
	}
}
