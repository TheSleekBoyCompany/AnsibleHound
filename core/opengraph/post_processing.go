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
			identityNode := graph.GetNode(edge.GetStartNodeID())
			identityEdges := graph.GetEdgesFromNode(identityNode.GetID())
			switch credentialType {

			case "Source Control":
				// NEEDED: User has to be able to modify a Project to configure a malicious SCM Server
				// Edge case: If a project uses an SCM credential that the user does not have `ATUse` on, they can compromise it anyway.
				var ok bool
				if edge.GetKind() == "ATUse" || edge.GetKind() == "ATAdmin" {
					for _, ie := range identityEdges {
						if slices.Contains(graph.GetNode(ie.GetEndNodeID()).GetKinds(), "ATProject") {
							if ie.GetKind() == "ATAdmin" {
								ok = true
								break
							}
						}
						if slices.Contains(graph.GetNode(ie.GetEndNodeID()).GetKinds(), "ATOrganization") {
							if ie.GetKind() == "ATAdmin" || ie.GetKind() == "ATProjectAdmin" {
								ok = true
								break
							}
						}
					}
				}
				if ok {
					edge = GenerateEdge("ATCompromiseWithFakeSCM", edge.GetStartNodeID(), node.GetID())
					graph.AddEdge(edge)
				}

			case "Thycotic Secret Server":
				if edge.GetKind() == "ATAdmin" {
					edge = GenerateEdge("ATCompromiseWithRequestbin", edge.GetStartNodeID(), node.GetID())
					graph.AddEdge(edge)
				}

			case "Machine":
				machineCredentialType := node.GetProperty("machine_credential_type").(string)

				if machineCredentialType == "ssh" && (edge.GetKind() == "ATUse" || edge.GetKind() == "ATAdmin") {
					// NEEDED: Control over the ansible playbook, to execute arbitrary commands on the execution environment with `delegate_to`
					var ok bool
					for _, ie := range identityEdges {
						if slices.Contains(graph.GetNode(ie.GetEndNodeID()).GetKinds(), "ATJobTemplate") {
							// If Admin of a Job Template you can potentially exploit an injection in the playbook.
							if ie.GetKind() == "ATAdmin" {
								ok = true
								break
							}
						}
						if slices.Contains(graph.GetNode(ie.GetEndNodeID()).GetKinds(), "ATOrganization") {
							if ie.GetKind() == "ATAdmin" || ie.GetKind() == "ATJobTemplateAdmin" || ie.GetKind() == "ATProjectAdmin" {
								ok = true
								break
							}
						}
						if slices.Contains(graph.GetNode(ie.GetEndNodeID()).GetKinds(), "ATProject") {
							// If ATAdmin of a project you can change the repo and configure a public one.
							if ie.GetKind() == "ATAdmin" {
								ok = true
								break
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
					for _, ie := range identityEdges {
						if slices.Contains(graph.GetNode(ie.GetEndNodeID()).GetKinds(), "ATInventory") {
							if ie.GetKind() == "ATAdmin" {
								ok = true
								break
							}
						}
						if slices.Contains(graph.GetNode(ie.GetEndNodeID()).GetKinds(), "ATOrganization") {
							if ie.GetKind() == "ATInventoryAdmin" || ie.GetKind() == "ATAdmin" {
								ok = true
								break
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
