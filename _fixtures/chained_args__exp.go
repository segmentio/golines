package fixtures

import "strings"

type Logger struct{}

func (l Logger) Str(k, v string) Logger                  { return l }
func (l Logger) Msgf(format string, args ...interface{}) {}

func Map[T, U any](collection []T, fn func(T, int) U) []U {
	return nil
}

func testChainedArgs() {
	messages := []string{}
	l := Logger{}

	// Long argument inside a chained method call - should split the arguments
	l.Str("key", strings.Join(Map(messages, func(message string, _ int) string { return message + message + message }), ",")).
		Msgf("Message with format %s", "arg")
}
