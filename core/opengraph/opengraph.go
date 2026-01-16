package opengraph

import (
	"github.com/TheManticoreProject/gopengraph"
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
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
}

type Node struct {
	Id         string            `json:"id"`
	Kinds      []string          `json:"kinds,omitempty"`
	Properties map[string]string `json:"properties"`
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
