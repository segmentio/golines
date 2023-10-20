package internal

import (
	"testing"

	"github.com/dave/dst"
	"github.com/stretchr/testify/assert"
)

func TestAnnotationStrings(t *testing.T) {
	assert.Equal(t, "// __golines:shorten:5", CreateAnnotation(5))
	assert.Equal(t, 5, ParseAnnotation("// __golines:shorten:5"))
	assert.Equal(t, -1, ParseAnnotation("// __golines:shorten:not_a_number"))
	assert.Equal(t, -1, ParseAnnotation("// not an annotation"))
	assert.True(t, IsAnnotation("// __golines:shorten:5"))
	assert.False(t, IsAnnotation("// not an annotation"))
}

func TestHasAnnotation(t *testing.T) {
	node1 := &dst.Ident{
		Name: "x",
		Decs: dst.IdentDecorations{
			NodeDecs: dst.NodeDecs{
				Start: []string{
					"// not an annotation",
					CreateAnnotation(55),
				},
			},
		},
	}
	assert.True(t, HasAnnotation(node1))

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
	assert.False(t, HasAnnotation(node2))
}
