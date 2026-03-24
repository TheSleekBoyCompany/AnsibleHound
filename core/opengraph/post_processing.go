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
	if isInEndNodeKinds(graph, identity, PROJECT_NODE) {
		if identity.GetKind() == ADMIN_ROLE_EDGE {
			ok = true
		}
	}
	if isInEndNodeKinds(graph, identity, ORGANIZATION_NODE) {
		if identity.GetKind() == ADMIN_ROLE_EDGE || identity.GetKind() == "ATProjectAdmin" {
			ok = true
		}
	}
	return ok
}

func identityCanControlJobTemplate(graph *gopengraph.OpenGraph, identity *edge.Edge) (ok bool) {
	if isInEndNodeKinds(graph, identity, JOB_TEMPLATE_NODE) {
		if identity.GetKind() == ADMIN_ROLE_EDGE {
			ok = true
		}
	}

	if isInEndNodeKinds(graph, identity, ORGANIZATION_NODE) {
		if identity.GetKind() == ADMIN_ROLE_EDGE || identity.GetKind() == "ATProjectAdmin" {
			ok = true
		}
	}
	return ok
}

func identityCanControlInventory(graph *gopengraph.OpenGraph, identity *edge.Edge) (ok bool) {
	if isInEndNodeKinds(graph, identity, INVENTORY_NODE) {
		if identity.GetKind() == ADMIN_ROLE_EDGE {
			ok = true
		}
	}
	if isInEndNodeKinds(graph, identity, ORGANIZATION_NODE) {
		if identity.GetKind() == INVENTORY_ADMIN_ROLE_EDGE || identity.GetKind() == "ATAdmin" {
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
		if edge.GetKind() == CONTAINS_EDGE {

			if isInEndNodeKinds(graph, edge, GROUP_NODE) {
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

			if isInEndNodeKinds(graph, edge, HOST_NODE) {
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

	credentialNodes := graph.GetNodesByKind(CREDENTIAL_NODE)
	for _, credentialNode := range credentialNodes {

		var credentialTypeName string
		var credentialTypeNode *node.Node
		for _, edge := range graph.GetEdgesFromNode(credentialNode.GetID()) {
			if edge.GetKind() == USES_TYPE_EDGE {
				credentialTypeNode = graph.GetNode(edge.GetEndNodeID())
				credentialTypeName = credentialTypeNode.GetProperty("name").(string)
			}
		}

		edges := graph.GetEdgesToNode(credentialNode.GetID())
		for _, edge := range edges {
			identityNode := graph.GetNode(edge.GetStartNodeID())
			identityEdges := graph.GetEdgesFromNode(identityNode.GetID())
			switch credentialTypeName {

			case SCM_CREDENTIAL_TYPE:
				if edge.GetKind() == USE_ROLE_EDGE || edge.GetKind() == ADMIN_ROLE_EDGE {
					for _, ie := range identityEdges {
						if identityCanControlPlaybook(graph, ie) {
							edge = GenerateEdge(COMPROMISE_WITH_FAKE_SCM_SERVER_POST_PROCESSING_EDGE,
								edge.GetStartNodeID(), credentialNode.GetID())
							graph.AddEdge(edge)
							break
						}
					}
				}

			case SECRET_SERVER_CREDENTIAL_TYPE:
				if edge.GetKind() == ADMIN_ROLE_EDGE {
					edge = GenerateEdge(COMPROMISE_WITH_REQUESTBIN_POST_PROCESSING_EDGE,
						edge.GetStartNodeID(), credentialNode.GetID())
					graph.AddEdge(edge)
				}

			case MACHINE_CREDENTIAL_TYPE:
				machineCredentialType := credentialNode.GetProperty("machine_credential_type").(string)
				if edge.GetKind() == USES_EDGE && isInStartNodeKinds(graph, edge, JOB_TEMPLATE_NODE) {
					var inventoryNode *node.Node
					jobTemplateNode := graph.GetNode(edge.GetStartNodeID())
					limit := jobTemplateNode.GetProperty("limit", "")
					jobTemplateEdges := graph.GetEdgesFromNode(jobTemplateNode.GetID())
					for _, jobTemplateEdge := range jobTemplateEdges {
						if jobTemplateEdge.GetKind() == USES_EDGE && isInEndNodeKinds(graph, jobTemplateEdge, INVENTORY_NODE) {
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
								e := GenerateEdge(VALID_FOR_POST_PROCESSING_EDGE, credentialNode.GetID(), hostEdge.GetEndNodeID())
								graph.AddEdge(e)
							}
						}
					}
				}

				if edge.GetKind() == USE_ROLE_EDGE || edge.GetKind() == ADMIN_ROLE_EDGE {

					if machineCredentialType == MACHINE_CREDENTIAL_SSH_SUBTYPE {
						for _, ie := range identityEdges {
							if identityCanControlJobTemplate(graph, ie) || identityCanControlPlaybook(graph, ie) {
								edge = GenerateEdge(SSH_HIJACK_AGENT_POST_PROCESSING_EDGE,
									edge.GetStartNodeID(), credentialNode.GetID())
								graph.AddEdge(edge)
								break
							}
						}
					}

					if machineCredentialType == MACHINE_CREDENTIAL_PASSWORD_SUBTYPE {
						for _, ie := range identityEdges {
							if identityCanControlInventory(graph, ie) {
								edge = GenerateEdge(COMPROMISE_WITH_HONEYPOT_POST_PROCESSING_EDGE,
									edge.GetStartNodeID(), credentialNode.GetID())
								graph.AddEdge(edge)
								break
							}
						}
					}

					for _, ie := range identityEdges {
						if identityCanControlInventory(graph, ie) || identityCanControlPlaybook(graph, ie) {
							edge = GenerateEdge(CAN_USE_IN_ADHOC_COMMANDS_POST_PROCESSING_EDGE,
								edge.GetStartNodeID(), credentialNode.GetID())
							graph.AddEdge(edge)
							break
						}
					}

				}

			default:
				if edge.GetKind() == USE_ROLE_EDGE || edge.GetKind() == ADMIN_ROLE_EDGE {
					hasEnvInjectors := credentialTypeNode.GetProperty(
						"injector_env", false).(bool)
					hasEVInjectors := credentialTypeNode.GetProperty(
						"injector_extra_vars", false).(bool)
					if hasEVInjectors || hasEnvInjectors {
						for _, ie := range identityEdges {
							if identityCanControlPlaybook(graph, ie) {
								edge = GenerateEdge(COMPROMISE_WITH_PLAYBOOK_POST_PROCESSING_EDGE,
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
