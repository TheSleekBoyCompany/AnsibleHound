package opengraph

import (
	"github.com/TheManticoreProject/gopengraph"
	"github.com/TheManticoreProject/gopengraph/node"
)

const SOURCE_KIND = "AnsibleBase"

type OutputJson struct {
	Metadata Metadata `json:"metadata"`
	Graph    Graph    `json:"graph"`
}

type Metadata struct {
	SourceKind string `json:"source_kind,omitempty"`
}

type Graph struct {
	Nodes []*node.Node `json:"nodes"`
	Edges []Edge       `json:"edges"`
}

type Edge struct {
	Kind  string       `json:"kind"`
	Start StartEndNode `json:"start"`
	End   StartEndNode `json:"end"`
}

type StartEndNode struct {
	Value string `json:"value"`
	Kind  string `json:"kind,omitempty"`
}

func InitGraph() (graph gopengraph.OpenGraph) {
	graph = *gopengraph.NewOpenGraph(SOURCE_KIND)
	return graph
}

func AddNodes(graph *gopengraph.OpenGraph, nodes []*node.Node) {
	for _, n := range nodes {
		graph.AddNode(n)
	}
}
