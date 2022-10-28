package shorten

import (
	"testing"

	"github.com/dave/dst"
	"github.com/stretchr/testify/assert"
)

func TestAnnotationStrings(t *testing.T) {
	assert.Equal(t, "// __golines:shorten:5", createAnnotation(5))
	assert.Equal(t, 5, parseAnnotation("// __golines:shorten:5"))
	assert.Equal(t, -1, parseAnnotation("// __golines:shorten:not_a_number"))
	assert.Equal(t, -1, parseAnnotation("// not an annotation"))
	assert.True(t, isAnnotation("// __golines:shorten:5"))
	assert.False(t, isAnnotation("// not an annotation"))
}

func TestHasAnnotation(t *testing.T) {
	node1 := &dst.Ident{
		Name: "x",
		Decs: dst.IdentDecorations{
			NodeDecs: dst.NodeDecs{
				Start: []string{
					"// not an annotation",
					createAnnotation(55),
				},
			},
		},
	}
	assert.True(t, hasAnnotation(node1))

	node2 := &dst.Ident{
		Name: "x",
		Decs: dst.IdentDecorations{
			NodeDecs: dst.NodeDecs{
				Start: []string{
					"// not an annotation",
				},
			},
		},
	}
	assert.False(t, hasAnnotation(node2))
}
