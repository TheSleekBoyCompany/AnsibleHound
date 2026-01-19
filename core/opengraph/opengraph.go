package opengraph

import (
	"ansible-hound/core/ansible"

	"github.com/TheManticoreProject/gopengraph"
	"github.com/TheManticoreProject/gopengraph/edge"
	"github.com/TheManticoreProject/gopengraph/node"
	"github.com/charmbracelet/log"
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

func AddEdge(graph *gopengraph.OpenGraph, edge *edge.Edge) {
	if !graph.AddEdge(edge) {
		log.Debugf("Edge failed validation, it was either a duplicate or one of the nodes did not exist in the graph.")
		log.Debugf("(%s)-[%s]-(%s)", edge.GetStartNodeID(), edge.GetKind(), edge.GetEndNodeID())
	}
}

func GenerateEdge(edgeKind string, startId string, endId string, startKind ...string) (e *edge.Edge) {

	e, err := edge.NewEdge(startId, endId, edgeKind, nil)
	if err != nil {
		log.Error(err)
	}

	return e
}

func GenerateNodes[T ansible.AnsibleType](objects map[int]T) (nodes []*node.Node) {
	for _, object := range objects {
		nodes = append(nodes, object.ToBHNode())
	}
	return nodes
}
