package main

import (
	"fmt"
	"html"
	"io"
	"strings"

	"github.com/dave/dst"
)

// GraphNode is a representation of a node in the AST graph.
type GraphNode struct {
	Type  string
	Value string
	Node  dst.Node
	Edges []*GraphEdge

	// Used for keeping track of node position during rendering
	level int
	seq   int
}

func (n *GraphNode) id() string {
	return fmt.Sprintf("%s_%d_%d", n.Type, n.level, n.seq)
}

// GraphEdge is a representation of an edge in the AST graph.
type GraphEdge struct {
	Dest         *GraphNode
	Relationship string
}

// CreateDot creates a dot representation of the graph associated with a dst node.
func CreateDot(node dst.Node, out io.Writer) error {
	root := NodeToGraphNode(node)

	dotGraph, err := WalkGraph(root)
	if err != nil {
		return err
	}

	_, err = out.Write([]byte(dotGraph))
	return err
}

// WalkGraph walks the graph starting at the argument root and returns
// a graphviz (dot) representation.
func WalkGraph(root *GraphNode) (string, error) {
	toProcess := []*GraphNode{root}
	processed := []*GraphNode{}
	outLines := []string{"digraph {"}

	var currLevel int
	var currSeq int

	// First, loop through the graph nodes to assign proper ids
	for len(toProcess) != 0 {
		currNode := toProcess[0]

		if currNode.level > currLevel {
			currLevel = currNode.level
			currSeq = 0
		}

		currNode.seq = currSeq
		currSeq++

		processed = append(processed, currNode)
		toProcess = toProcess[1:]

		for _, edge := range currNode.Edges {
			edge.Dest.level = currLevel + 1
			toProcess = append(toProcess, edge.Dest)
		}
	}

	// Then, fill out the graph in dot format
	for _, node := range processed {
		var nodeLabel string
		var nodeFormat string

		if HasAnnotation(node.Node) {
			nodeFormat = ",penwidth=3.0"
		}

		if node.Value != "" {
			nodeLabel = fmt.Sprintf(
				"%s<br/><font point-size=\"11.0\" face=\"courier\" color=\"#777777\">%s</font>",
				node.Type,
				html.EscapeString(node.Value),
			)
		} else {
			nodeLabel = node.Type
		}

		outLines = append(
			outLines,
			fmt.Sprintf(
				"\t%s[label=<%s>,shape=\"box\"%s]",
				node.id(),
				nodeLabel,
				nodeFormat,
			),
		)

		for _, edge := range node.Edges {
			outLines = append(
				outLines,
				fmt.Sprintf(
					"\t%s->%s[label=\"%s\",fontsize=12.0]",
					node.id(),
					edge.Dest.id(),
					edge.Relationship,
				),
			)
		}
	}

	outLines = append(outLines, "}")
	return strings.Join(outLines, "\n"), nil
}
