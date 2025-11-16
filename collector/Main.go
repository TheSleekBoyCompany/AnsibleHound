package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	core "ansible-hound/core"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

func launchGathering(client http.Client, targetUrl *url.URL,
	username string, password string, outdir string) {

	output := core.OutputJson{
		Metadata: core.Metadata{
			SourceKind: "AnsibleBase",
		},
	}
	nodes := []core.Node{}
	edges := []core.Edge{}

	// -- Check if credentials are valid

	log.Info("Authenticating on Ansible Worx/Tower instance.")
	_, err := core.AuthenticateOnAnsibleInstance(client, *targetUrl, username, password, core.ME_ENDPOINT)
	if err != nil {
		log.Error("Unable to authenticate on Ansible WorX/Tower instance due to invalid credentials.")
		log.Error(err)
		return
	}

	// -- Gathering Ansible Instance information --

	log.Info("Gathering Ansible Worx/Tower instance information.")
	instance, err := core.GatherAnsibleInstance(client, *targetUrl)
	if err != nil {
		log.Fatalf("Unable to gather Ansible WorX/Tower information (%s).", targetUrl)
	}

	instance.Name = targetUrl.Host

	instanceNode := instance.ToBHNode()
	nodes = append(nodes, instanceNode)

	// -- Gathering all nodes --

	log.Info("Gathering Users.")
	users, err := core.GatherObject[*core.User](
		instance.InstallUUID, client, *targetUrl, username, password, core.USERS_ENDPOINT,
	)
	if err != nil {
		log.Error("An error occured while gathering Users, skipping.")
		log.Error(err)
	}

	log.Info("Gathering User Roles.")
	for i, user := range users {
		userRolesEndpoint := fmt.Sprintf(core.USER_ROLES_ENDPOINT, user.ID)
		roles, err := core.GatherObject[*core.Role](
			instance.InstallUUID, client, *targetUrl, username, password, userRolesEndpoint)
		if err != nil {
			log.Error("An error occured while gathering User Roles.")
			log.Error(err)
			continue
		}
		user.Roles = roles
		users[i] = user
	}
	userNodes := core.GenerateNodes(users)
	nodes = append(nodes, userNodes...)

	log.Info("Gathering Hosts.")
	hosts, err := core.GatherObject[*core.Host](
		instance.InstallUUID, client, *targetUrl, username, password, core.HOSTS_ENDPOINT,
	)
	if err != nil {
		log.Error("An error occured while gathering Hosts, skipping.")
		log.Error(err)
	}
	hostNodes := core.GenerateNodes(hosts)
	nodes = append(nodes, hostNodes...)

	log.Info("Gathering Groups.")
	groups, err := core.GatherObject[*core.Group](
		instance.InstallUUID, client, *targetUrl, username, password, core.GROUPS_ENDPOINT,
	)
	if err != nil {
		log.Error("An error occured while gathering Users, skipping.")
		log.Error(err)
	}
	groupsNodes := core.GenerateNodes(groups)
	nodes = append(nodes, groupsNodes...)

	log.Info("Gathering Group Hosts.")
	for i, group := range groups {

		groupHostsEndpoint := fmt.Sprintf(core.GROUP_HOSTS_ENDPOINT, group.ID)
		hosts, err := core.GatherObject[*core.Host](
			instance.InstallUUID, client, *targetUrl, username, password, groupHostsEndpoint,
		)
		if err != nil {
			log.Error("An error occured while gathering Team Roles.")
			log.Error(err)
			continue
		}
		group.Hosts = hosts
		groups[i] = group
	}

	log.Info("Gathering Jobs.")
	jobs, err := core.GatherObject[*core.Job](
		instance.InstallUUID, client, *targetUrl, username, password, core.JOBS_ENDPOINT,
	)
	if err != nil {
		log.Error("An error occured while gathering Jobs, skipping.")
		log.Error(err)
	}
	jobsNodes := core.GenerateNodes(jobs)
	nodes = append(nodes, jobsNodes...)

	log.Info("Gathering Job Templates.")
	jobTemplates, err := core.GatherObject[*core.JobTemplate](
		instance.InstallUUID, client, *targetUrl, username, password, core.JOB_TEMPLATE_ENDPOINT,
	)
	if err != nil {
		log.Error("An error occured while gathering Job Templates, skipping.")
		log.Error(err)
	}
	log.Info("Gathering Job Templates Credentials.")
	for i, jobTemplate := range jobTemplates {

		jobTemplatesCredentialsEndpoint := fmt.Sprintf(core.
			JOB_TEMPLATE_CREDENTIALS_ENDPOINT, jobTemplate.ID)
		credentials, err := core.GatherObject[*core.Credential](
			instance.InstallUUID, client, *targetUrl, username, password, jobTemplatesCredentialsEndpoint,
		)
		if err != nil {
			log.Error("An error occured while gathering Job Template Credentials.")
			log.Error(err)
			continue
		}

		jobTemplate.Credentials = credentials
		jobTemplates[i] = jobTemplate
	}
	jobTemplatesNodes := core.GenerateNodes(jobTemplates)
	nodes = append(nodes, jobTemplatesNodes...)

	log.Info("Gathering Inventories.")
	inventories, err := core.GatherObject[*core.Inventory](
		instance.InstallUUID, client, *targetUrl, username, password, core.INVENTORIES_ENDPOINT,
	)
	if err != nil {
		log.Error("An error occured while gathering Inventories, skipping.")
		log.Error(err)
	}
	inventoriesNodes := core.GenerateNodes(inventories)
	nodes = append(nodes, inventoriesNodes...)

	log.Info("Gathering Organizations.")
	organizations, err := core.GatherObject[*core.Organization](
		instance.InstallUUID, client, *targetUrl, username, password, core.ORGANIZATIONS_ENDPOINT,
	)
	if err != nil {
		log.Error("An error occured while gathering Organizations, skipping.")
		log.Error(err)
	}
	organizationNodes := core.GenerateNodes(organizations)
	nodes = append(nodes, organizationNodes...)

	log.Info("Gathering Credentials.")
	credentials, err := core.GatherObject[*core.Credential](
		instance.InstallUUID, client, *targetUrl, username, password, core.CREDENTIALS_ENDPOINT,
	)
	if err != nil {
		log.Error("An error occured while gathering Credentials, skipping.")
		log.Error(err)
	}
	credentialNodes := core.GenerateNodes(credentials)
	nodes = append(nodes, credentialNodes...)

	log.Info("Gathering Credential Types.")
	credentialTypes, err := core.GatherObject[*core.CredentialType](
		instance.InstallUUID, client, *targetUrl, username, password, core.CREDENTIAL_TYPES_ENDPOINT,
	)
	if err != nil {
		log.Error("An error occured while gathering Credential Types, skipping.")
		log.Error(err)
	}
	credentialTypesNodes := core.GenerateNodes(credentialTypes)
	nodes = append(nodes, credentialTypesNodes...)

	log.Info("Gathering Projects.")
	projects, err := core.GatherObject[*core.Project](
		instance.InstallUUID, client, *targetUrl, username, password, core.PROJECTS_ENDPOINT,
	)
	if err != nil {
		log.Error("An error occured while gathering Projects, skipping.")
		log.Error(err)
	}
	projectNodes := core.GenerateNodes(projects)
	nodes = append(nodes, projectNodes...)

	log.Info("Gathering Teams.")
	teams, err := core.GatherObject[*core.Team](
		instance.InstallUUID, client, *targetUrl, username, password, core.TEAMS_ENDPOINT,
	)
	if err != nil {
		log.Error("An error occured while gathering Teams, skipping.")
		log.Error(err)
	}

	log.Info("Gathering Team Roles.")
	for i, team := range teams {

		teamRolesEndpoint := fmt.Sprintf(core.TEAM_ROLES_ENDPOINT, team.ID)
		roles, err := core.GatherObject[*core.Role](
			instance.InstallUUID, client, *targetUrl, username, password, teamRolesEndpoint,
		)
		if err != nil {
			log.Error("An error occured while gathering Team Roles.")
			log.Error(err)
			continue
		}

		teamMembersEndpoint := fmt.Sprintf(core.TEAM_USERS_ENDPOINT, team.ID)
		members, err := core.GatherObject[*core.User](
			instance.InstallUUID, client, *targetUrl, username, password, teamMembersEndpoint,
		)
		if err != nil {
			log.Error("An error occured while gathering Team Members.")
			log.Error(err)
			continue
		}
		team.Roles = roles
		team.Members = members
		teams[i] = team
	}
	teamNodes := core.GenerateNodes(teams)
	nodes = append(nodes, teamNodes...)

	// -- Creating all edges --

	log.Info("Linking Instance and Organizations.")
	kind := "ATContains"
	for _, organization := range organizations {
		edge := core.GenerateEdge(kind, instance.OID, organization.OID)
		edges = append(edges, edge)
	}

	log.Info("Linking Organizations and Inventories.")
	kind = "ATContains"
	for _, inventory := range inventories {
		if core.HasAccessTo(organizations, inventory.Organization) {
			edge := core.GenerateEdge(kind, organizations[inventory.Organization].OID, inventory.OID)
			edges = append(edges, edge)
		}
	}

	log.Info("Linking Inventories and Hosts.")
	kind = "ATContains"
	for _, host := range hosts {
		if core.HasAccessTo(inventories, host.Inventory) {
			edge := core.GenerateEdge(kind, inventories[host.Inventory].OID, host.OID)
			edges = append(edges, edge)
		}
	}

	log.Info("Linking Inventories and Groups.")
	kind = "ATContains"
	for _, group := range groups {
		if core.HasAccessTo(inventories, group.Inventory) {
			edge := core.GenerateEdge(kind, inventories[group.Inventory].OID, group.OID)
			edges = append(edges, edge)
		}
	}

	log.Info("Linking Groups and Hosts.")
	kind = "ATContains"
	for _, group := range groups {
		for _, host := range group.Hosts {
			if core.HasAccessTo(hosts, host.ID) {
				edge := core.GenerateEdge(kind, groups[group.ID].OID, host.OID)
				edges = append(edges, edge)
			}
		}
	}

	log.Info("Linking Job Templates and Jobs.")
	kind = "ATContains"
	for _, job := range jobs {
		if core.HasAccessTo(jobTemplates, job.UnifiedJobTemplate) {
			edge := core.GenerateEdge(kind, jobTemplates[job.UnifiedJobTemplate].OID, job.OID)
			edges = append(edges, edge)
		}
	}

	log.Info("Linking Organizations and Job Templates.")
	kind = "ATContains"
	for _, jobTemplate := range jobTemplates {
		if core.HasAccessTo(organizations, jobTemplate.Organization) {
			edge := core.GenerateEdge(kind, organizations[jobTemplate.Organization].OID, jobTemplate.OID)
			edges = append(edges, edge)
		}
	}

	log.Info("Linking Organizations and Credentials.")
	kind = "ATContains"
	for _, credential := range credentials {
		if core.HasAccessTo(organizations, credential.Organization) {
			edge := core.GenerateEdge(kind, organizations[credential.Organization].OID, credential.OID)
			edges = append(edges, edge)
		}
	}

	log.Info("Linking Organizations and Projects")
	kind = "ATContains"
	for _, project := range projects {
		if core.HasAccessTo(organizations, project.Organization) {
			edge := core.GenerateEdge(kind, organizations[project.Organization].OID, project.OID)
			edges = append(edges, edge)
		}
	}

	log.Info("Linking Job Templates and Projects.")
	kind = "ATUses"
	for _, jobTemplate := range jobTemplates {
		if core.HasAccessTo(projects, jobTemplate.Project) {
			edge := core.GenerateEdge(kind, jobTemplate.OID, projects[jobTemplate.Project].OID)
			edges = append(edges, edge)
		}
	}

	log.Info("Linking Job Template and Inventories.")
	kind = "ATUses"
	for _, jobTemplate := range jobTemplates {
		if core.HasAccessTo(inventories, jobTemplate.Inventory) {
			edge := core.GenerateEdge(kind, jobTemplate.OID, inventories[jobTemplate.Inventory].OID)
			edges = append(edges, edge)
		}
	}

	log.Info("Linking Credentials to Credential Type.")
	kind = "ATUsesType"
	for _, credential := range credentials {
		if core.HasAccessTo(credentialTypes, credential.CredentialType) {
			edge := core.GenerateEdge(kind, credential.OID, credentialTypes[credential.CredentialType].OID)
			edges = append(edges, edge)
		}
	}

	log.Info("Linking Job Template and Credentials.")
	kind = "ATUses"
	for _, jobTemplate := range jobTemplates {
		for _, credential := range jobTemplate.Credentials {
			edge := core.GenerateEdge(kind, jobTemplate.OID, credential.OID)
			edges = append(edges, edge)
		}
	}

	log.Info("Linking User Roles.")
	for _, user := range users {

		for _, role := range user.Roles {
			kind := "AT" + strings.ReplaceAll(role.Name, " ", "")

			var edge core.Edge

			switch role.SummaryFields.ResourceType {

			case "organization":
				if core.HasAccessTo(organizations, role.SummaryFields.ResourceId) {
					edge = core.GenerateEdge(kind, user.OID, organizations[role.SummaryFields.ResourceId].OID)
					edges = append(edges, edge)
				}
			case "inventory":
				if core.HasAccessTo(inventories, role.SummaryFields.ResourceId) {
					edge = core.GenerateEdge(kind, user.OID, inventories[role.SummaryFields.ResourceId].OID)
					edges = append(edges, edge)
				}
			case "team":
				if core.HasAccessTo(teams, role.SummaryFields.ResourceId) {
					edge = core.GenerateEdge(kind, user.OID, teams[role.SummaryFields.ResourceId].OID)
					edges = append(edges, edge)
				}
			case "credential":
				if core.HasAccessTo(credentials, role.SummaryFields.ResourceId) {
					edge = core.GenerateEdge(kind, user.OID, credentials[role.SummaryFields.ResourceId].OID)
					edges = append(edges, edge)
				}
			case "job_template":
				if core.HasAccessTo(jobTemplates, role.SummaryFields.ResourceId) {
					edge = core.GenerateEdge(kind, user.OID, jobTemplates[role.SummaryFields.ResourceId].OID)
					edges = append(edges, edge)
				}
			}

		}
	}

	log.Info("Linking Team Roles.")
	for _, team := range teams {
		for _, role := range team.Roles {
			kind := "AT" + strings.ReplaceAll(role.Name, " ", "")

			var edge core.Edge

			switch role.SummaryFields.ResourceType {

			case "organization":
				if core.HasAccessTo(organizations, role.SummaryFields.ResourceId) {
					edge = core.GenerateEdge(kind, team.OID, organizations[role.SummaryFields.ResourceId].OID)
					edges = append(edges, edge)
				}
			case "inventory":
				if core.HasAccessTo(inventories, role.SummaryFields.ResourceId) {
					edge = core.GenerateEdge(kind, team.OID, inventories[role.SummaryFields.ResourceId].OID)
					edges = append(edges, edge)
				}
			case "team":
				if core.HasAccessTo(teams, role.SummaryFields.ResourceId) {
					edge = core.GenerateEdge(kind, team.OID, teams[role.SummaryFields.ResourceId].OID)
					edges = append(edges, edge)
				}
			case "credential":
				if core.HasAccessTo(credentials, role.SummaryFields.ResourceId) {
					edge = core.GenerateEdge(kind, team.OID, credentials[role.SummaryFields.ResourceId].OID)
					edges = append(edges, edge)
				}
			case "job_template":
				if core.HasAccessTo(jobTemplates, role.SummaryFields.ResourceId) {
					edge = core.GenerateEdge(kind, team.OID, jobTemplates[role.SummaryFields.ResourceId].OID)
					edges = append(edges, edge)
				}
			}
		}
	}

	output.Graph = core.Graph{
		Nodes: nodes,
		Edges: edges,
	}

	outputJson, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	core.WriteToFile(
		outputJson,
		path.Join(outdir, core.CalculateName("output")),
	)

}

var ingestCmd = &cobra.Command{
	Use:   "collect",
	Short: "Collector for Ansible AWX/Tower, for use with the Ansible project.",
	Run: func(cmd *cobra.Command, args []string) {

		target, _ := cmd.Flags().GetString("target")
		if target == "" {
			log.Fatal("Empty target provided.")
		}

		targetUrl, err := url.Parse(target)
		if err != nil {
			log.Fatal("Invalid target URL specified.\n%s", err)
		}

		username, _ := cmd.Flags().GetString("username")
		if username == "" {
			log.Fatal("Empty username provided.")
		}

		password, _ := cmd.Flags().GetString("password")
		if password == "" {
			log.Fatal("Empty password provided.")
		}

		verbose, _ := cmd.Flags().GetBool("verbose")
		if verbose {
			log.SetLevel(log.DebugLevel)
		}

		var proxyURL *url.URL
		proxy, _ := cmd.Flags().GetString("proxy")
		if proxy != "" {
			proxyURL, err = url.Parse(proxy)
			if err != nil {
				log.Fatalf("Invalid proxy URL specified.\n%s", err)
			}
		}

		outdir, _ := cmd.Flags().GetString("outdir")

		skipVerifySSL, _ := cmd.Flags().GetBool("skip-verify-ssl")

		client := core.InitClient(proxyURL, skipVerifySSL)

		launchGathering(client, targetUrl, username, password, outdir)

	},
}

func main() {

	ingestCmd.Flags().StringP("target", "t", "", "Target URL of AWX/Tower instance.")
	_ = ingestCmd.MarkFlagRequired("target")

	ingestCmd.Flags().StringP("username", "u", "", "Username to use for authenticatio.")
	_ = ingestCmd.MarkFlagRequired("username")

	ingestCmd.Flags().StringP("password", "p", "", "Password to use for authentication.")
	_ = ingestCmd.MarkFlagRequired("password")

	ingestCmd.Flags().StringP("proxy", "", "", "(optional) Configure HTTP/HTTPS proxy.")
	ingestCmd.Flags().StringP("outdir", "", "", "(optional) Output directory for the json files.")
	ingestCmd.Flags().BoolP("verbose", "v", false, "(optional) Enable debug logs.")
	ingestCmd.Flags().BoolP("skip-verify-ssl", "k", false, "(optional) Skips SSL/TLS verification.")

	if err := ingestCmd.Execute(); err != nil {
		msg := fmt.Sprintf("CLI error: %v\n", err)
		log.Error(msg)
		os.Exit(1)
	}

}
