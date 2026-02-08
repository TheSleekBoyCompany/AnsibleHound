package opengraph

import (
	"ansible-hound/core/ansible"

	"github.com/Ramoreik/gopengraph"
	"github.com/Ramoreik/gopengraph/edge"
	"github.com/Ramoreik/gopengraph/node"
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

func GenerateEdge(edgeKind string, startId string, endId string) (e *edge.Edge) {

	e, err := edge.NewEdge(startId, endId, edgeKind, MATCH_BY_ID, MATCH_BY_ID, ANSIBLE_BASE, ANSIBLE_BASE, nil)
	if err != nil {
		log.Error(err)
	}

	return e
}

func GenerateEdgeCustom(edgeKind string, startId string, endId string, startMatchBy string, endMatchBy string, startNodeKind string, endNodeKind string) (e *edge.Edge) {

	e, err := edge.NewEdge(startId, endId, edgeKind, startMatchBy, endMatchBy, startNodeKind, endNodeKind, nil)
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
