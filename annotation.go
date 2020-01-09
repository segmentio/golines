package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/dave/dst"
)

const annotationPrefix = "// __golines:shorten:"

// CreateAnnotation generates the text of a comment that will annotate long lines.
func CreateAnnotation(length int) string {
	return fmt.Sprintf(
		"%s%d",
		annotationPrefix,
		length,
	)
}

// IsAnnotation determines whether the given line is an annotation created with CreateAnnotation.
func IsAnnotation(line string) bool {
	return strings.HasPrefix(
		strings.Trim(line, " \t"),
		annotationPrefix,
	)
}

// HasAnnotation determines whether the given AST node has a line length annotation on it.
func HasAnnotation(node dst.Node) bool {
	startDecorations := node.Decorations().Start.All()
	return len(startDecorations) > 0 &&
		IsAnnotation(startDecorations[len(startDecorations)-1])
}

// HasAnnotationRecursive determines whether the given node or one of its children has a
// golines annotation on it. It's currently implemented for function declarations, fields,
// and call expressions only.
func HasAnnotationRecursive(node dst.Node) bool {
	if HasAnnotation(node) {
		return true
	}

	switch n := node.(type) {
	case *dst.FuncDecl:
		if n.Type != nil && n.Type.Params != nil {
			for _, item := range n.Type.Params.List {
				if HasAnnotationRecursive(item) {
					return true
				}
			}
		}
	case *dst.Field:
		return HasAnnotation(n)
	case *dst.CallExpr:
		for _, arg := range n.Args {
			if HasAnnotation(arg) {
				return true
			}
		}
	}

	return false
}

// ParseAnnotation returns the line length encoded in a golines annotation. If none is found,
// it returns -1.
func ParseAnnotation(line string) int {
	if IsAnnotation(line) {
		components := strings.SplitN(line, ":", 3)
		val, err := strconv.Atoi(components[2])
		if err != nil {
			return -1
		}
		return val
	}
	return -1
}
