package shorten

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/dave/dst"
)

const annotationPrefix = "// __golines:shorten:"

// createAnnotation generates the text of a comment that will annotate long lines.
func createAnnotation(length int) string {
	return fmt.Sprintf(
		"%s%d",
		annotationPrefix,
		length,
	)
}

// isAnnotation determines whether the given line is an annotation created with createAnnotation.
func isAnnotation(line string) bool {
	return strings.HasPrefix(
		strings.Trim(line, " \t"),
		annotationPrefix,
	)
}

// hasAnnotation determines whether the given AST node has a line length annotation on it.
func hasAnnotation(node dst.Node) bool {
	startDecorations := node.Decorations().Start.All()
	return len(startDecorations) > 0 &&
		isAnnotation(startDecorations[len(startDecorations)-1])
}

// hasTailAnnotation determines whether the given AST node has a line length annotation at its
// end. This is needed to catch long function declarations with inline interface definitions.
func hasTailAnnotation(node dst.Node) bool {
	endDecorations := node.Decorations().End.All()
	return len(endDecorations) > 0 &&
		isAnnotation(endDecorations[len(endDecorations)-1])
}

// hasAnnotationRecursive determines whether the given node or one of its children has a
// golines annotation on it. It's currently implemented for function declarations, fields,
// call expressions, and selector expressions only.
func hasAnnotationRecursive(node dst.Node) bool {
	if hasAnnotation(node) {
		return true
	}

	switch n := node.(type) {
	case *dst.FuncDecl:
		if n.Type != nil && n.Type.Params != nil {
			for _, item := range n.Type.Params.List {
				if hasAnnotationRecursive(item) {
					return true
				}
			}
		}
	case *dst.Field:
		return hasTailAnnotation(n) || hasAnnotationRecursive(n.Type)
	case *dst.SelectorExpr:
		return hasAnnotation(n.Sel) || hasAnnotation(n.X)
	case *dst.CallExpr:
		if hasAnnotationRecursive(n.Fun) {
			return true
		}

		for _, arg := range n.Args {
			if hasAnnotation(arg) {
				return true
			}
		}
	case *dst.InterfaceType:
		for _, field := range n.Methods.List {
			if hasAnnotationRecursive(field) {
				return true
			}
		}
	}

	return false
}

// parseAnnotation returns the line length encoded in a golines annotation. If none is found,
// it returns -1.
func parseAnnotation(line string) int {
	if isAnnotation(line) {
		components := strings.SplitN(line, ":", 3)
		val, err := strconv.Atoi(components[2])
		if err != nil {
			return -1
		}
		return val
	}
	return -1
}
