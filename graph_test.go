package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/dave/dst/decorator"
	"github.com/stretchr/testify/assert"
)

var (
	testCode = `
package mypackage

func myfunc() {
	return 5
}
`

	expDot = `
digraph {
	File_0_0[label=<File>,shape="box"]
	File_0_0->Ident_1_0[label="Name",fontsize=12.0]
	File_0_0->FuncDecl_1_1[label="Decls",fontsize=12.0]
	Ident_1_0[label=<Ident<br/><font point-size="11.0" face="courier" color="#777777">mypackage</font>>,shape="box"]
	FuncDecl_1_1[label=<FuncDecl>,shape="box"]
	FuncDecl_1_1->Ident_2_0[label="Name",fontsize=12.0]
	FuncDecl_1_1->FieldList_2_1[label="Params",fontsize=12.0]
	FuncDecl_1_1->BlockStmt_2_2[label="Body",fontsize=12.0]
	Ident_2_0[label=<Ident<br/><font point-size="11.0" face="courier" color="#777777">myfunc</font>>,shape="box"]
	FieldList_2_1[label=<FieldList>,shape="box"]
	BlockStmt_2_2[label=<BlockStmt>,shape="box"]
	BlockStmt_2_2->ReturnStmt_3_0[label="List",fontsize=12.0]
	ReturnStmt_3_0[label=<ReturnStmt>,shape="box"]
	ReturnStmt_3_0->BasicLit_4_0[label="Results",fontsize=12.0]
	BasicLit_4_0[label=<BasicLit<br/><font point-size="11.0" face="courier" color="#777777">5</font>>,shape="box"]
}
`
)

func TestCreateDot(t *testing.T) {
	node, err := decorator.Parse(testCode)
	assert.Nil(t, err)

	out := &bytes.Buffer{}
	err = CreateDot(node, out)
	assert.Nil(t, err)

	assert.Equal(t, strings.TrimSpace(expDot), out.String())
}
