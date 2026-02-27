package opengraph

import (
	"slices"

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

			if credentialType == "Thycotic Secret Server" {
				if edge.GetKind() == "ATAdmin" {
					edge = GenerateEdge("ATCompromiseWithRequestbin", edge.GetStartNodeID(), node.GetID())
					graph.AddEdge(edge)
				}
			}

			if credentialType == "Machine" {
				machineCredentialType := node.GetProperty("machine_credential_type").(string)
				identityNode := graph.GetNode(edge.GetStartNodeID())
				identityEdges := graph.GetEdgesFromNode(identityNode.GetID())

				if machineCredentialType == "ssh" && (edge.GetKind() == "ATUse" || edge.GetKind() == "ATAdmin") {
					// NEEDED: Control over the ansible playbook, to execute arbitrary commands on the execution environment with `delegate_to`
					var ok bool
					for _, edge := range identityEdges {
						if slices.Contains(graph.GetNode(edge.GetEndNodeID()).GetKinds(), "ATJobTemplate") {
							// If Admin of a Job Template you can potentially exploit an injection in the playbook.
							if edge.GetKind() == "ATAdmin" {
								ok = true
							}
						}
						if slices.Contains(graph.GetNode(edge.GetEndNodeID()).GetKinds(), "ATOrganization") {
							if edge.GetKind() == "ATAdmin" || edge.GetKind() == "ATJobTemplateAdmin" || edge.GetKind() == "ATProjectAdmin" {
								ok = true
							}
						}
						if slices.Contains(graph.GetNode(edge.GetEndNodeID()).GetKinds(), "ATProject") {
							// If ATAdmin of a project you can change the repo and configure a public one.
							if edge.GetKind() == "ATAdmin" {
								ok = true
							}
						}
					}
					if ok {
						edge = GenerateEdge("ATSSHHijackAgent", edge.GetStartNodeID(), node.GetID())
						graph.AddEdge(edge)
					}
				}

				if machineCredentialType == "password" && (edge.GetKind() == "ATUse" || edge.GetKind() == "ATAdmin") {
					// NEEDED: user able to modify any inventory, in order to target an attacker controlled server.
					var ok bool
					for _, edge := range identityEdges {
						if slices.Contains(graph.GetNode(edge.GetEndNodeID()).GetKinds(), "ATInventory") {
							if edge.GetKind() == "ATAdmin" {
								ok = true
							}
						}
						if slices.Contains(graph.GetNode(edge.GetEndNodeID()).GetKinds(), "ATOrganization") {
							if edge.GetKind() == "ATInventoryAdmin" || edge.GetKind() == "ATAdmin" {
								ok = true
							}
						}
					}
					if ok {
						edge = GenerateEdge("ATCompromiseWithHoneypot", edge.GetStartNodeID(), node.GetID())
						graph.AddEdge(edge)
					}
				}

			}
		}
	}
}
