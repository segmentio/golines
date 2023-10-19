package fixtures

import "fmt"

type Chain struct{}

func (c *Chain) ChainCall(arg1 string, arg2 string, arg3 string) *Chain {
	return c
}

func NewChain() *Chain {
	return &Chain{}
}

func ChainedCalls() {
	c := Chain{}
	c.ChainCall("a long argument", "another long argument", "a third long argument").
		ChainCall("a long argument2", "another long argument2", "a third long argument2").
		ChainCall("a long argument3", "another long argument3", "a third long argument3")
	NewChain().ChainCall(
		"a really really really really really long argument4",
		"another really really really really really long argument4",
		fmt.Sprintf("%v", "this is a long method"),
	).
		ChainCall("a really really really really really long argument5", "another really really really really really long argument5", "a third really really really really really long argument5").
		ChainCall("a", "b", fmt.Sprintf("%v", "this is a long method"))
	NewChain().ChainCall("a", "b", "c").ChainCall("d", "e", "f")
	NewChain().ChainCall("a", "b", "c").
		ChainCall("d", "e", "f").ChainCall("g", "h", "i")
}
