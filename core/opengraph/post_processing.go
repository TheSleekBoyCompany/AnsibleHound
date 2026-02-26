package opengraph

import (
	"github.com/Ramoreik/gopengraph"
)

func PostProcessingCredentials(graph *gopengraph.OpenGraph) {
	credentialNodes := graph.GetNodesByKind("ATCredential")
	for _, node := range credentialNodes {

		var credentialType string
		for _, edge := range graph.GetEdgesFromNode(node.GetID()) {
			if edge.GetKind() == "ATUsesType" {
				credentialType = graph.GetNode(edge.GetEndNodeID()).GetProperty("name").(string)
			}
		}

		edges := graph.GetEdgesToNode(node.GetID())
		for _, edge := range edges {
			if credentialType == "Machine" {
				machineCredentialType := node.GetProperty("machine_credential_type").(string)

				if machineCredentialType == "ssh" && (edge.GetKind() == "ATUse" || edge.GetKind() == "ATAdmin") {
					edge = GenerateEdge("ATSSHHijackAgent", edge.GetStartNodeID(), node.GetID())
					graph.AddEdge(edge)
				}

				if machineCredentialType == "password" && (edge.GetKind() == "ATUse" || edge.GetKind() == "ATAdmin") {
					edge = GenerateEdge("ATCompromiseWithHoneypot", edge.GetStartNodeID(), node.GetID())
					graph.AddEdge(edge)
				}
			}
		}
	}
}
