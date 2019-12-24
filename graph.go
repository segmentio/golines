package main

import (
	"fmt"
	"os"

	"github.com/dave/dst"
	"github.com/emicklei/dot"
)

type GraphNode struct {
	Type  string
	Value string
	Node  dst.Node
	Edges []*GraphEdge
}

type GraphEdge struct {
	Dest         *GraphNode
	Relationship string
}

func CreateDot(node dst.Node) {
	graph := dot.NewGraph(dot.Directed)
	root := NodeToGraphNode(node)
	WalkGraph(graph, root)

	graph.Write(os.Stdout)
}

func WalkGraph(graph *dot.Graph, root *GraphNode) dot.Node {
	var nodeLabel string

	if root.Value != "" {
		nodeLabel = fmt.Sprintf("%s\n%s", root.Type, root.Value)
	} else {
		nodeLabel = fmt.Sprintf("%s", root.Type)
	}

	dotRoot := graph.Node(nodeLabel).Box()

	for _, edge := range root.Edges {
		dotDest := WalkGraph(graph, edge.Dest)
		graph.Edge(dotRoot, dotDest, edge.Relationship)
	}

	return dotRoot
}
