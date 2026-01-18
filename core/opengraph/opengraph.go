package opengraph

import (
	"github.com/TheManticoreProject/gopengraph"
	"github.com/TheManticoreProject/gopengraph/node"
)

const SOURCE_KIND = "AnsibleBase"

func InitGraph() (graph gopengraph.OpenGraph) {
	graph = *gopengraph.NewOpenGraph(SOURCE_KIND)
	return graph
}

func AddNodes(graph *gopengraph.OpenGraph, nodes []*node.Node) {
	for _, n := range nodes {
		graph.AddNode(n)
	}
}
