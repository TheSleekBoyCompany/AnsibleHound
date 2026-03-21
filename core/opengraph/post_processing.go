package opengraph

import (
	"slices"

	"github.com/Ramoreik/awxlimit/pkg/awxlimit"
	"github.com/Ramoreik/gopengraph"
	"github.com/Ramoreik/gopengraph/edge"
	"github.com/Ramoreik/gopengraph/node"
	"github.com/charmbracelet/log"
)

func isInEndNodeKinds(graph *gopengraph.OpenGraph, edge *edge.Edge, kind string) (ok bool) {
	endNode := graph.GetNode(edge.GetEndNodeID())
	if slices.Contains(endNode.GetKinds(), kind) {
		ok = true
	}
	return ok
}

func isInStartNodeKinds(graph *gopengraph.OpenGraph, edge *edge.Edge, kind string) (ok bool) {
	endNode := graph.GetNode(edge.GetStartNodeID())
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

func getInventoryMembers(graph *gopengraph.OpenGraph, inv *node.Node) (inventory awxlimit.Inventory) {
	edges := graph.GetEdgesFromNode(inv.GetID())
	hosts := []string{}
	groups := []awxlimit.Group{}
	for _, edge := range edges {
		if edge.GetKind() == "ATContains" {

			if isInEndNodeKinds(graph, edge, "ATGroup") {
				groupNode := graph.GetNode(edge.GetEndNodeID())
				groupEdges := graph.GetEdgesFromNode(groupNode.GetID())
				groupHosts := []string{}
				for _, groupEdge := range groupEdges {
					groupHost := graph.GetNode(groupEdge.GetEndNodeID())
					groupHosts = append(groupHosts, groupHost.GetProperty("name").(string))
				}
				group := awxlimit.Group{
					Name:     groupNode.GetProperty("name").(string),
					Hosts:    groupHosts,
					Children: []string{},
				}
				groups = append(groups, group)
			}

			if isInEndNodeKinds(graph, edge, "ATHost") {
				hostNode := graph.GetNode(edge.GetEndNodeID())
				hosts = append(hosts, hostNode.GetProperty("name").(string))
			}
		}
	}
	inventory = awxlimit.Inventory{
		Hosts:  hosts,
		Groups: groups,
	}
	return inventory
}

func PostProcessingCredentials(graph *gopengraph.OpenGraph) {

	log.Info("Handling post processing edges for Credentials.")

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

				// ATValidFor - Maps credentials directly to machines, based on JobTemplates using them
				// 1: Check for ATUses edge between Credential and ATJobTemplate
				// 2: Check for the `limit` value of the JobTemplate and the Inventory
				// 3: Find the hosts matching the `limit` and `inventory` values. (https://docs.ansible.com/projects/ansible/latest/inventory_guide/intro_patterns.html)
				// 4: Create ATValidFor edge between `ATCredential` and `ATHost`
				if edge.GetKind() == "ATUses" && isInStartNodeKinds(graph, edge, "ATJobTemplate") {
					var inventoryNode *node.Node
					jobTemplateNode := graph.GetNode(edge.GetStartNodeID())
					limit := jobTemplateNode.GetProperty("limit", "")
					jobTemplateEdges := graph.GetEdgesFromNode(jobTemplateNode.GetID())
					for _, jobTemplateEdge := range jobTemplateEdges {
						if jobTemplateEdge.GetKind() == "ATUses" && isInEndNodeKinds(graph, jobTemplateEdge, "ATInventory") {
							inventoryNode = graph.GetNode(jobTemplateEdge.GetEndNodeID())
						}
					}
					if inventoryNode != nil {
						matched := []string{}
						if limit != "" {
							inventory := getInventoryMembers(graph, inventoryNode)
							matched, _ = awxlimit.MatchHosts(limit.(string), inventory)
						}
						hostEdges := graph.GetEdgesFromNode(inventoryNode.GetID())
						for _, hostEdge := range hostEdges {
							if limit == "" || slices.Contains(matched, graph.GetNode(hostEdge.GetEndNodeID()).GetProperty("name").(string)) {
								e := GenerateEdge("ATValidFor", credentialNode.GetID(), hostEdge.GetEndNodeID())
								graph.AddEdge(e)
							}
						}
					}
				}

				if edge.GetKind() == "ATUse" || edge.GetKind() == "ATAdmin" {

					if machineCredentialType == "ssh" {
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

					if machineCredentialType == "password" {
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

					for _, ie := range identityEdges {
						if identityCanControlInventory(graph, ie) || identityCanControlPlaybook(graph, ie) {
							edge = GenerateEdge("ATCanUseInADHOCCommands",
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
