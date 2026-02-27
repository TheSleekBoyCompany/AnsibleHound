package opengraph

import (
	"slices"

	"github.com/Ramoreik/gopengraph"
	"github.com/Ramoreik/gopengraph/edge"
	"github.com/Ramoreik/gopengraph/node"
)

func isInEndNodeKinds(graph *gopengraph.OpenGraph, edge *edge.Edge, kind string) (ok bool) {
	endNode := graph.GetNode(edge.GetEndNodeID())
	if slices.Contains(endNode.GetKinds(), kind) {
		ok = true
	}
	return ok
}

func identityCanControlPlaybook(graph *gopengraph.OpenGraph, identity *edge.Edge) (ok bool) {
	if isInEndNodeKinds(graph, identity, "ATProject") {
		if identity.GetKind() == "ATAdmin" {
			ok = true
		}
	}
	if isInEndNodeKinds(graph, identity, "ATOrganization") {
		if identity.GetKind() == "ATAdmin" || identity.GetKind() == "ATProjectAdmin" {
			ok = true
		}
	}
	return ok
}

func identityCanControlJobTemplate(graph *gopengraph.OpenGraph, identity *edge.Edge) (ok bool) {
	if isInEndNodeKinds(graph, identity, "ATJobTemplate") {
		if identity.GetKind() == "ATAdmin" {
			ok = true
		}
	}

	if isInEndNodeKinds(graph, identity, "ATOrganization") {
		if identity.GetKind() == "ATAdmin" || identity.GetKind() == "ATProjectAdmin" {
			ok = true
		}
	}
	return ok
}

func identityCanControlInventory(graph *gopengraph.OpenGraph, identity *edge.Edge) (ok bool) {
	if isInEndNodeKinds(graph, identity, "ATInventory") {
		if identity.GetKind() == "ATAdmin" {
			ok = true
		}
	}
	if isInEndNodeKinds(graph, identity, "ATOrganization") {
		if identity.GetKind() == "ATInventoryAdmin" || identity.GetKind() == "ATAdmin" {
			ok = true
		}
	}
	return ok
}

func PostProcessingCredentials(graph *gopengraph.OpenGraph) {
	credentialNodes := graph.GetNodesByKind("ATCredential")
	for _, credentialNode := range credentialNodes {

		var credentialTypeName string
		var credentialTypeNode *node.Node
		for _, edge := range graph.GetEdgesFromNode(credentialNode.GetID()) {
			if edge.GetKind() == "ATUsesType" {
				credentialTypeNode = graph.GetNode(edge.GetEndNodeID())
				credentialTypeName = credentialTypeNode.GetProperty("name").(string)
			}
		}

		edges := graph.GetEdgesToNode(credentialNode.GetID())
		for _, edge := range edges {
			identityNode := graph.GetNode(edge.GetStartNodeID())
			identityEdges := graph.GetEdgesFromNode(identityNode.GetID())
			switch credentialTypeName {

			case "Source Control":
				// NEEDED: User has to be able to modify a Project to configure a malicious SCM Server
				// Edge case: If a project uses an SCM credential that the user does not have `ATUse` on, they can compromise it anyway.
				if edge.GetKind() == "ATUse" || edge.GetKind() == "ATAdmin" {
					for _, ie := range identityEdges {
						if identityCanControlPlaybook(graph, ie) {
							edge = GenerateEdge("ATCompromiseWithFakeSCM",
								edge.GetStartNodeID(), credentialNode.GetID())
							graph.AddEdge(edge)
							break
						}
					}
				}

			case "Thycotic Secret Server":
				if edge.GetKind() == "ATAdmin" {
					edge = GenerateEdge("ATCompromiseWithRequestbin",
						edge.GetStartNodeID(), credentialNode.GetID())
					graph.AddEdge(edge)
				}

			case "Machine":
				machineCredentialType := credentialNode.GetProperty("machine_credential_type").(string)
				if machineCredentialType == "ssh" && (edge.GetKind() == "ATUse" || edge.GetKind() == "ATAdmin") {
					// NEEDED: Control over the ansible playbook, to execute arbitrary commands on the execution environment with `delegate_to`
					for _, ie := range identityEdges {
						if identityCanControlJobTemplate(graph, ie) || identityCanControlPlaybook(graph, ie) {
							edge = GenerateEdge("ATSSHHijackAgent",
								edge.GetStartNodeID(), credentialNode.GetID())
							graph.AddEdge(edge)
							break
						}
					}
				}

				if machineCredentialType == "password" && (edge.GetKind() == "ATUse" || edge.GetKind() == "ATAdmin") {
					// NEEDED: user able to modify any inventory, in order to target an attacker controlled server.
					for _, ie := range identityEdges {
						if identityCanControlInventory(graph, ie) {
							edge = GenerateEdge("ATCompromiseWithHoneypot",
								edge.GetStartNodeID(), credentialNode.GetID())
							graph.AddEdge(edge)
							break
						}
					}
				}

			default:
				// NEEDED: User can control the playbook being executed.
				// NEEDED: Credential Type has injectors
				if edge.GetKind() == "ATUse" || edge.GetKind() == "ATAdmin" {
					hasEnvInjectors := credentialTypeNode.GetProperty(
						"injector_env", false).(bool)
					hasEVInjectors := credentialTypeNode.GetProperty(
						"injector_extra_vars", false).(bool)
					if hasEVInjectors || hasEnvInjectors {
						for _, ie := range identityEdges {
							if identityCanControlPlaybook(graph, ie) {
								edge = GenerateEdge("ATCompromiseWithMaliciousPlaybook",
									edge.GetStartNodeID(), credentialNode.GetID())
								graph.AddEdge(edge)
								break
							}
						}
					}
				}

			}
		}
	}
}
