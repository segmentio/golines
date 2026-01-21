package fixtures

var _ = &mybigtype{
	arg1: "biglongstring1",
	arg2: longfunctioncall(param1),
	arg3: []string{"list1"},
	arg4: bigvariable,
}
