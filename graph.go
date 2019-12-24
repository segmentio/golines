package main

import (
	"fmt"
	"os"

	"github.com/dave/dst"
	"github.com/emicklei/dot"
)

type GraphNode struct {
	Level    int
	Sequence int
	Type     string
	Value    string
	Node     dst.Node
	Edges    []*GraphEdge
}

func (n *GraphNode) ID() string {
	return fmt.Sprintf("%s_%d_%d", n.Type, n.Level, n.Sequence)
}

type GraphEdge struct {
	Dest         *GraphNode
	Relationship string
}

func CreateDot(node dst.Node) {
	g := dot.NewGraph(dot.Directed)
	root := NodeToGraphNode(node, 0, 0)

	g.Node(root.ID())
	g.Write(os.Stdout)
}
