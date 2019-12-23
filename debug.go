package main

import (
	"github.com/dave/dst"
	"github.com/davecgh/go-spew/spew"
)

func SpewFile(file *dst.File) {
	spew.Sdump(file)
}

func DotFile(file *dst.File) {

}
