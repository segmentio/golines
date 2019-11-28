package fixtures

type Chain struct{}

func (c *Chain) ChainCall(a string, b string, c string) *Chain {
	return c
}

func ChainedCalls() {
	c := Chain{}
	c.ChainCall("a long argument", "another long argument", "a third long argument").ChainCall("a long argument2", "another long argument2", "a third long argument2").ChainCall("a long argument3", "another long argument3", "a third long argument3")
}
