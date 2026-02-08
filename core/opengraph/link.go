package opengraph

import (
	"ansible-hound/core/ansible"
	"ansible-hound/core/gather"
	"path"
	"strings"

	"github.com/Ramoreik/gopengraph"
	"github.com/Ramoreik/gopengraph/edge"
	"github.com/charmbracelet/log"
)

func LinkOrganization(graph *gopengraph.OpenGraph, instanceOID string,
	organizations map[int]*ansible.Organization,
	inventories map[int]*ansible.Inventory, jobTemplates map[int]*ansible.JobTemplate,
	credentials map[int]*ansible.Credential, projects map[int]*ansible.Project) {

	log.Info("Linking Instance and Organizations.")
	edgeKind := "ATContains"
	for _, organization := range organizations {
		edge := GenerateEdge(edgeKind, instanceOID, organization.OID)
		AddEdge(graph, edge)
	}

	log.Info("Linking Organizations and Inventories.")
	edgeKind = "ATContains"
	for _, inventory := range inventories {
		if gather.HasAccessTo(organizations, inventory.Organization) {
			edge := GenerateEdge(edgeKind, organizations[inventory.Organization].OID, inventory.OID)
			AddEdge(graph, edge)
		}
	}

	log.Info("Linking Organizations and Job Templates.")
	edgeKind = "ATContains"
	for _, jobTemplate := range jobTemplates {
		if gather.HasAccessTo(organizations, jobTemplate.Organization) {
			edge := GenerateEdge(edgeKind, organizations[jobTemplate.Organization].OID, jobTemplate.OID)
			AddEdge(graph, edge)
		}
	}

	log.Info("Linking Organizations and Credentials.")
	edgeKind = "ATContains"
	for _, credential := range credentials {
		if gather.HasAccessTo(organizations, credential.Organization) {
			edge := GenerateEdge(edgeKind, organizations[credential.Organization].OID, credential.OID)
			AddEdge(graph, edge)
		}
	}

	log.Info("Linking Organizations and Projects")
	edgeKind = "ATContains"
	for _, project := range projects {
		if gather.HasAccessTo(organizations, project.Organization) {
			edge := GenerateEdge(edgeKind, organizations[project.Organization].OID, project.OID)
			AddEdge(graph, edge)
		}
	}

}

func LinkInventory(graph *gopengraph.OpenGraph, inventories map[int]*ansible.Inventory,
	hosts map[int]*ansible.Host, groups map[int]*ansible.Group) {

	log.Info("Linking Inventories and Hosts.")
	edgeKind := "ATContains"
	for _, host := range hosts {
		if gather.HasAccessTo(inventories, host.Inventory) {
			edge := GenerateEdge(edgeKind, inventories[host.Inventory].OID, host.OID)
			AddEdge(graph, edge)
		}
	}

	log.Info("Linking Inventories and Groups.")
	edgeKind = "ATContains"
	for _, group := range groups {
		if gather.HasAccessTo(inventories, group.Inventory) {
			edge := GenerateEdge(edgeKind, inventories[group.Inventory].OID, group.OID)
			AddEdge(graph, edge)
		}
	}

	log.Info("Linking Groups and Hosts.")
	edgeKind = "ATContains"
	for _, group := range groups {
		for _, host := range group.Hosts {
			if gather.HasAccessTo(hosts, host.ID) {
				edge := GenerateEdge(edgeKind, groups[group.ID].OID, host.OID)
				AddEdge(graph, edge)
			}
		}
	}

}

func LinkJobTemplates(graph *gopengraph.OpenGraph, jobTemplates map[int]*ansible.JobTemplate,
	jobs map[int]*ansible.Job, projects map[int]*ansible.Project,
	inventories map[int]*ansible.Inventory, credentials map[int]*ansible.Credential,
	credentialTypes map[int]*ansible.CredentialType) {

	log.Info("Linking Job Templates and Jobs.")
	edgeKind := "ATContains"
	for _, job := range jobs {
		if gather.HasAccessTo(jobTemplates, job.UnifiedJobTemplate) {
			edge := GenerateEdge(edgeKind, jobTemplates[job.UnifiedJobTemplate].OID, job.OID)
			AddEdge(graph, edge)
		}
	}

	log.Info("Linking Job Templates and Projects.")
	edgeKind = "ATUses"
	for _, jobTemplate := range jobTemplates {
		if gather.HasAccessTo(projects, jobTemplate.Project) {
			edge := GenerateEdge(edgeKind, jobTemplate.OID, projects[jobTemplate.Project].OID)
			AddEdge(graph, edge)
		}
	}

	log.Info("Linking Job Template and Inventories.")
	edgeKind = "ATUses"
	for _, jobTemplate := range jobTemplates {
		if gather.HasAccessTo(inventories, jobTemplate.Inventory) {
			edge := GenerateEdge(edgeKind, jobTemplate.OID, inventories[jobTemplate.Inventory].OID)
			AddEdge(graph, edge)
		}
	}

	log.Info("Linking Job Template and Credentials.")
	edgeKind = "ATUses"
	for _, jobTemplate := range jobTemplates {
		for _, credential := range jobTemplate.Credentials {
			edge := GenerateEdge(edgeKind, jobTemplate.OID, credential.OID)
			AddEdge(graph, edge)
		}
	}

	log.Info("Linking Credentials to Credential Type.")
	edgeKind = "ATUsesType"
	for _, credential := range credentials {
		if gather.HasAccessTo(credentialTypes, credential.CredentialType) {
			edge := GenerateEdge(edgeKind, credential.OID, credentialTypes[credential.CredentialType].OID)
			AddEdge(graph, edge)
		}
	}

}

func LinkUserRoles(graph *gopengraph.OpenGraph, users map[int]*ansible.User,
	organizations map[int]*ansible.Organization, inventories map[int]*ansible.Inventory,
	teams map[int]*ansible.Team, credentials map[int]*ansible.Credential,
	jobTemplates map[int]*ansible.JobTemplate) {

	log.Info("Linking User Roles.")
	for _, user := range users {

		for _, role := range user.Roles {
			edgeKind := "AT" + strings.ReplaceAll(role.Name, " ", "")

			var edge *edge.Edge

			switch role.SummaryFields.ResourceType {

			case "organization":
				if gather.HasAccessTo(organizations, role.SummaryFields.ResourceId) {
					edge = GenerateEdge(edgeKind, user.OID, organizations[role.SummaryFields.ResourceId].OID)
				}
			case "inventory":
				if gather.HasAccessTo(inventories, role.SummaryFields.ResourceId) {
					edge = GenerateEdge(edgeKind, user.OID, inventories[role.SummaryFields.ResourceId].OID)
				}
			case "team":
				if gather.HasAccessTo(teams, role.SummaryFields.ResourceId) {
					edge = GenerateEdge(edgeKind, user.OID, teams[role.SummaryFields.ResourceId].OID)
				}
			case "credential":
				if gather.HasAccessTo(credentials, role.SummaryFields.ResourceId) {
					edge = GenerateEdge(edgeKind, user.OID, credentials[role.SummaryFields.ResourceId].OID)
				}
			case "job_template":
				if gather.HasAccessTo(jobTemplates, role.SummaryFields.ResourceId) {
					edge = GenerateEdge(edgeKind, user.OID, jobTemplates[role.SummaryFields.ResourceId].OID)
				}
			}
			if edge != nil {
				AddEdge(graph, edge)
			}
		}
	}
}

func LinkTeamRoles(graph *gopengraph.OpenGraph, users map[int]*ansible.User,
	organizations map[int]*ansible.Organization, inventories map[int]*ansible.Inventory,
	teams map[int]*ansible.Team, credentials map[int]*ansible.Credential,
	jobTemplates map[int]*ansible.JobTemplate) {

	log.Info("Linking Team Roles.")
	for _, team := range teams {
		for _, role := range team.Roles {
			edgeKind := "AT" + strings.ReplaceAll(role.Name, " ", "")

			var edge *edge.Edge

			switch role.SummaryFields.ResourceType {

			case "organization":
				if gather.HasAccessTo(organizations, role.SummaryFields.ResourceId) {
					edge = GenerateEdge(edgeKind, team.OID, organizations[role.SummaryFields.ResourceId].OID)
				}
			case "inventory":
				if gather.HasAccessTo(inventories, role.SummaryFields.ResourceId) {
					edge = GenerateEdge(edgeKind, team.OID, inventories[role.SummaryFields.ResourceId].OID)
				}
			case "team":
				if gather.HasAccessTo(teams, role.SummaryFields.ResourceId) {
					edge = GenerateEdge(edgeKind, team.OID, teams[role.SummaryFields.ResourceId].OID)
				}
			case "credential":
				if gather.HasAccessTo(credentials, role.SummaryFields.ResourceId) {
					edge = GenerateEdge(edgeKind, team.OID, credentials[role.SummaryFields.ResourceId].OID)
				}
			case "job_template":
				if gather.HasAccessTo(jobTemplates, role.SummaryFields.ResourceId) {
					edge = GenerateEdge(edgeKind, team.OID, jobTemplates[role.SummaryFields.ResourceId].OID)
				}
			}
			if edge != nil {
				AddEdge(graph, edge)
			}
		}
	}
}

func LinkAD(graph *gopengraph.OpenGraph, ldap gather.AHLdap, users map[int]*ansible.User) {

	if (ldap != gather.AHLdap{}) {

		log.Info("Linking LDAP Users.")
		edgeKind := "SyncedToAHUser"

		conn, err := gather.Connect(ldap)

		if err != nil {
			log.Fatal("Connection failed:", err)
		}
		defer conn.Close()

		for _, user := range users {
			if user.ExternalAccount == LDAP_VALUE {
				if user.LdapDn != "" {
					ldap_dn := user.LdapDn
					objectSid, _ := gather.Search(conn, ldap_dn)
					edge := GenerateEdgeCustom(edgeKind, objectSid, user.OID, MATCH_BY_ID, MATCH_BY_ID, ACTIVE_DIRECTORY_BASE, ANSIBLE_BASE)
					graph.AddEdgeWithoutValidation(edge)
				}
			}
		}
	} else {
		log.Warn("Skipping linking LDAP Users since the user used for authentication is a local user")
	}
}

func LinkGitHub(graph *gopengraph.OpenGraph, projects map[int]*ansible.Project) {

	if projects != nil {

		log.Info("Linking GitHub Repositories.")
		edgeKind := "ATIsLinkedTo"

		for _, project := range projects {
			if project.ScmType == GIT_SCM_TYPE {
				if project.ScmUrl != "" {
					scmUrl := project.ScmUrl
					repositoryName := strings.TrimSuffix(path.Base(scmUrl), DOT_GIT_SCM_TYPE)
					edge := GenerateEdgeCustom(edgeKind, project.OID, repositoryName, MATCH_BY_ID, MATCH_BY_NAME, ANSIBLE_BASE, GITHUB_BASE)
					graph.AddEdgeWithoutValidation(edge)
				}
			}
		}
	}
}
