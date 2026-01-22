package main

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"

	"ansible-hound/core"
	"ansible-hound/core/ansible"
	"ansible-hound/core/opengraph"

	"github.com/Ramoreik/gopengraph/edge"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

func launchGathering(client core.AHClient, targetUrl *url.URL, outdir string, ldap core.AHLdap) {

	graph := opengraph.InitGraph()

	// -- Check if credentials are valid --

	log.Info("Authenticating on Ansible Worx/Tower instance.")
	_, err := core.AuthenticateOnAnsibleInstance(client, *targetUrl, core.ME_ENDPOINT)
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
	graph.AddNode(instanceNode)

	// -- Gathering all nodes --

	log.Info("Gathering Users.")
	users, err := core.GatherObject[*ansible.User](
		instance.InstallUUID, client, *targetUrl, core.USERS_ENDPOINT,
	)
	if err != nil {
		log.Error("An error occured while gathering Users, skipping.")
		log.Error(err)
	}

	log.Info("Gathering User Roles.")
	for i, user := range users {
		userRolesEndpoint := fmt.Sprintf(core.USER_ROLES_ENDPOINT, user.ID)
		roles, err := core.GatherObject[*ansible.Role](
			instance.InstallUUID, client, *targetUrl, userRolesEndpoint)
		if err != nil {
			log.Error("An error occured while gathering User Roles.")
			log.Error(err)
			continue
		}
		user.Roles = roles
		users[i] = user
	}
	userNodes := opengraph.GenerateNodes(users)
	opengraph.AddNodes(&graph, userNodes)

	log.Info("Gathering Hosts.")
	hosts, err := core.GatherObject[*ansible.Host](
		instance.InstallUUID, client, *targetUrl, core.HOSTS_ENDPOINT,
	)
	if err != nil {
		log.Error("An error occured while gathering Hosts, skipping.")
		log.Error(err)
	}
	hostNodes := opengraph.GenerateNodes(hosts)
	opengraph.AddNodes(&graph, hostNodes)

	log.Info("Gathering Groups.")
	groups, err := core.GatherObject[*ansible.Group](
		instance.InstallUUID, client, *targetUrl, core.GROUPS_ENDPOINT,
	)
	if err != nil {
		log.Error("An error occured while gathering Users, skipping.")
		log.Error(err)
	}
	groupsNodes := opengraph.GenerateNodes(groups)
	opengraph.AddNodes(&graph, groupsNodes)

	log.Info("Gathering Group Hosts.")
	for i, group := range groups {

		groupHostsEndpoint := fmt.Sprintf(core.GROUP_HOSTS_ENDPOINT, group.ID)
		hosts, err := core.GatherObject[*ansible.Host](
			instance.InstallUUID, client, *targetUrl, groupHostsEndpoint,
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
	jobs, err := core.GatherObject[*ansible.Job](
		instance.InstallUUID, client, *targetUrl, core.JOBS_ENDPOINT,
	)
	if err != nil {
		log.Error("An error occured while gathering Jobs, skipping.")
		log.Error(err)
	}
	jobsNodes := opengraph.GenerateNodes(jobs)
	opengraph.AddNodes(&graph, jobsNodes)

	log.Info("Gathering Job Templates.")
	jobTemplates, err := core.GatherObject[*ansible.JobTemplate](
		instance.InstallUUID, client, *targetUrl, core.JOB_TEMPLATE_ENDPOINT,
	)
	if err != nil {
		log.Error("An error occured while gathering Job Templates, skipping.")
		log.Error(err)
	}
	log.Info("Gathering Job Templates Credentials.")
	for i, jobTemplate := range jobTemplates {

		jobTemplatesCredentialsEndpoint := fmt.Sprintf(core.
			JOB_TEMPLATE_CREDENTIALS_ENDPOINT, jobTemplate.ID)
		credentials, err := core.GatherObject[*ansible.Credential](
			instance.InstallUUID, client, *targetUrl, jobTemplatesCredentialsEndpoint,
		)
		if err != nil {
			log.Error("An error occured while gathering Job Template Credentials.")
			log.Error(err)
			continue
		}

		jobTemplate.Credentials = credentials
		jobTemplates[i] = jobTemplate
	}
	jobTemplatesNodes := opengraph.GenerateNodes(jobTemplates)
	opengraph.AddNodes(&graph, jobTemplatesNodes)

	log.Info("Gathering Inventories.")
	inventories, err := core.GatherObject[*ansible.Inventory](
		instance.InstallUUID, client, *targetUrl, core.INVENTORIES_ENDPOINT,
	)
	if err != nil {
		log.Error("An error occured while gathering Inventories, skipping.")
		log.Error(err)
	}
	inventoriesNodes := opengraph.GenerateNodes(inventories)
	opengraph.AddNodes(&graph, inventoriesNodes)

	log.Info("Gathering Organizations.")
	organizations, err := core.GatherObject[*ansible.Organization](
		instance.InstallUUID, client, *targetUrl, core.ORGANIZATIONS_ENDPOINT,
	)
	if err != nil {
		log.Error("An error occured while gathering Organizations, skipping.")
		log.Error(err)
	}
	organizationNodes := opengraph.GenerateNodes(organizations)
	opengraph.AddNodes(&graph, organizationNodes)

	log.Info("Gathering Credentials.")
	credentials, err := core.GatherObject[*ansible.Credential](
		instance.InstallUUID, client, *targetUrl, core.CREDENTIALS_ENDPOINT,
	)
	if err != nil {
		log.Error("An error occured while gathering Credentials, skipping.")
		log.Error(err)
	}
	credentialNodes := opengraph.GenerateNodes(credentials)
	opengraph.AddNodes(&graph, credentialNodes)

	log.Info("Gathering Credential Types.")
	credentialTypes, err := core.GatherObject[*ansible.CredentialType](
		instance.InstallUUID, client, *targetUrl, core.CREDENTIAL_TYPES_ENDPOINT,
	)
	if err != nil {
		log.Error("An error occured while gathering Credential Types, skipping.")
		log.Error(err)
	}
	credentialTypesNodes := opengraph.GenerateNodes(credentialTypes)
	opengraph.AddNodes(&graph, credentialTypesNodes)

	log.Info("Gathering Projects.")
	projects, err := core.GatherObject[*ansible.Project](
		instance.InstallUUID, client, *targetUrl, core.PROJECTS_ENDPOINT,
	)
	if err != nil {
		log.Error("An error occured while gathering Projects, skipping.")
		log.Error(err)
	}
	projectNodes := opengraph.GenerateNodes(projects)
	opengraph.AddNodes(&graph, projectNodes)

	log.Info("Gathering Teams.")
	teams, err := core.GatherObject[*ansible.Team](
		instance.InstallUUID, client, *targetUrl, core.TEAMS_ENDPOINT,
	)
	if err != nil {
		log.Error("An error occured while gathering Teams, skipping.")
		log.Error(err)
	}

	log.Info("Gathering Team Roles.")
	for i, team := range teams {

		teamRolesEndpoint := fmt.Sprintf(core.TEAM_ROLES_ENDPOINT, team.ID)
		roles, err := core.GatherObject[*ansible.Role](
			instance.InstallUUID, client, *targetUrl, teamRolesEndpoint,
		)
		if err != nil {
			log.Error("An error occured while gathering Team Roles.")
			log.Error(err)
			continue
		}

		teamMembersEndpoint := fmt.Sprintf(core.TEAM_USERS_ENDPOINT, team.ID)
		members, err := core.GatherObject[*ansible.User](
			instance.InstallUUID, client, *targetUrl, teamMembersEndpoint,
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
	teamNodes := opengraph.GenerateNodes(teams)
	opengraph.AddNodes(&graph, teamNodes)

	// -- Creating all edges --

	log.Info("Linking Instance and Organizations.")
	edgeKind := "ATContains"
	for _, organization := range organizations {
		edge := opengraph.GenerateEdge(edgeKind, instance.OID, organization.OID)
		opengraph.AddEdge(&graph, edge)
	}

	log.Info("Linking Organizations and Inventories.")
	edgeKind = "ATContains"
	for _, inventory := range inventories {
		if core.HasAccessTo(organizations, inventory.Organization) {
			edge := opengraph.GenerateEdge(edgeKind, organizations[inventory.Organization].OID, inventory.OID)
			opengraph.AddEdge(&graph, edge)
		}
	}

	log.Info("Linking Inventories and Hosts.")
	edgeKind = "ATContains"
	for _, host := range hosts {
		if core.HasAccessTo(inventories, host.Inventory) {
			edge := opengraph.GenerateEdge(edgeKind, inventories[host.Inventory].OID, host.OID)
			opengraph.AddEdge(&graph, edge)
		}
	}

	log.Info("Linking Inventories and Groups.")
	edgeKind = "ATContains"
	for _, group := range groups {
		if core.HasAccessTo(inventories, group.Inventory) {
			edge := opengraph.GenerateEdge(edgeKind, inventories[group.Inventory].OID, group.OID)
			opengraph.AddEdge(&graph, edge)
		}
	}

	log.Info("Linking Groups and Hosts.")
	edgeKind = "ATContains"
	for _, group := range groups {
		for _, host := range group.Hosts {
			if core.HasAccessTo(hosts, host.ID) {
				edge := opengraph.GenerateEdge(edgeKind, groups[group.ID].OID, host.OID)
				opengraph.AddEdge(&graph, edge)
			}
		}
	}

	log.Info("Linking Job Templates and Jobs.")
	edgeKind = "ATContains"
	for _, job := range jobs {
		if core.HasAccessTo(jobTemplates, job.UnifiedJobTemplate) {
			edge := opengraph.GenerateEdge(edgeKind, jobTemplates[job.UnifiedJobTemplate].OID, job.OID)
			opengraph.AddEdge(&graph, edge)
		}
	}

	log.Info("Linking Organizations and Job Templates.")
	edgeKind = "ATContains"
	for _, jobTemplate := range jobTemplates {
		if core.HasAccessTo(organizations, jobTemplate.Organization) {
			edge := opengraph.GenerateEdge(edgeKind, organizations[jobTemplate.Organization].OID, jobTemplate.OID)
			opengraph.AddEdge(&graph, edge)
		}
	}

	log.Info("Linking Organizations and Credentials.")
	edgeKind = "ATContains"
	for _, credential := range credentials {
		if core.HasAccessTo(organizations, credential.Organization) {
			edge := opengraph.GenerateEdge(edgeKind, organizations[credential.Organization].OID, credential.OID)
			opengraph.AddEdge(&graph, edge)
		}
	}

	log.Info("Linking Organizations and Projects")
	edgeKind = "ATContains"
	for _, project := range projects {
		if core.HasAccessTo(organizations, project.Organization) {
			edge := opengraph.GenerateEdge(edgeKind, organizations[project.Organization].OID, project.OID)
			opengraph.AddEdge(&graph, edge)
		}
	}

	log.Info("Linking Job Templates and Projects.")
	edgeKind = "ATUses"
	for _, jobTemplate := range jobTemplates {
		if core.HasAccessTo(projects, jobTemplate.Project) {
			edge := opengraph.GenerateEdge(edgeKind, jobTemplate.OID, projects[jobTemplate.Project].OID)
			opengraph.AddEdge(&graph, edge)
		}
	}

	log.Info("Linking Job Template and Inventories.")
	edgeKind = "ATUses"
	for _, jobTemplate := range jobTemplates {
		if core.HasAccessTo(inventories, jobTemplate.Inventory) {
			edge := opengraph.GenerateEdge(edgeKind, jobTemplate.OID, inventories[jobTemplate.Inventory].OID)
			opengraph.AddEdge(&graph, edge)
		}
	}

	log.Info("Linking Credentials to Credential Type.")
	edgeKind = "ATUsesType"
	for _, credential := range credentials {
		if core.HasAccessTo(credentialTypes, credential.CredentialType) {
			edge := opengraph.GenerateEdge(edgeKind, credential.OID, credentialTypes[credential.CredentialType].OID)
			opengraph.AddEdge(&graph, edge)
		}
	}

	log.Info("Linking Job Template and Credentials.")
	edgeKind = "ATUses"
	for _, jobTemplate := range jobTemplates {
		for _, credential := range jobTemplate.Credentials {
			edge := opengraph.GenerateEdge(edgeKind, jobTemplate.OID, credential.OID)
			opengraph.AddEdge(&graph, edge)
		}
	}

	log.Info("Linking User Roles.")
	for _, user := range users {

		for _, role := range user.Roles {
			edgeKind := "AT" + strings.ReplaceAll(role.Name, " ", "")

			var edge *edge.Edge

			switch role.SummaryFields.ResourceType {

			case "organization":
				if core.HasAccessTo(organizations, role.SummaryFields.ResourceId) {
					edge = opengraph.GenerateEdge(edgeKind, user.OID, organizations[role.SummaryFields.ResourceId].OID)
				}
			case "inventory":
				if core.HasAccessTo(inventories, role.SummaryFields.ResourceId) {
					edge = opengraph.GenerateEdge(edgeKind, user.OID, inventories[role.SummaryFields.ResourceId].OID)
				}
			case "team":
				if core.HasAccessTo(teams, role.SummaryFields.ResourceId) {
					edge = opengraph.GenerateEdge(edgeKind, user.OID, teams[role.SummaryFields.ResourceId].OID)
				}
			case "credential":
				if core.HasAccessTo(credentials, role.SummaryFields.ResourceId) {
					edge = opengraph.GenerateEdge(edgeKind, user.OID, credentials[role.SummaryFields.ResourceId].OID)
				}
			case "job_template":
				if core.HasAccessTo(jobTemplates, role.SummaryFields.ResourceId) {
					edge = opengraph.GenerateEdge(edgeKind, user.OID, jobTemplates[role.SummaryFields.ResourceId].OID)
				}
			}
			if edge != nil {
				opengraph.AddEdge(&graph, edge)
			}
		}
	}

	log.Info("Linking Team Roles.")
	for _, team := range teams {
		for _, role := range team.Roles {
			edgeKind := "AT" + strings.ReplaceAll(role.Name, " ", "")

			var edge *edge.Edge

			switch role.SummaryFields.ResourceType {

			case "organization":
				if core.HasAccessTo(organizations, role.SummaryFields.ResourceId) {
					edge = opengraph.GenerateEdge(edgeKind, team.OID, organizations[role.SummaryFields.ResourceId].OID)
				}
			case "inventory":
				if core.HasAccessTo(inventories, role.SummaryFields.ResourceId) {
					edge = opengraph.GenerateEdge(edgeKind, team.OID, inventories[role.SummaryFields.ResourceId].OID)
				}
			case "team":
				if core.HasAccessTo(teams, role.SummaryFields.ResourceId) {
					edge = opengraph.GenerateEdge(edgeKind, team.OID, teams[role.SummaryFields.ResourceId].OID)
				}
			case "credential":
				if core.HasAccessTo(credentials, role.SummaryFields.ResourceId) {
					edge = opengraph.GenerateEdge(edgeKind, team.OID, credentials[role.SummaryFields.ResourceId].OID)
				}
			case "job_template":
				if core.HasAccessTo(jobTemplates, role.SummaryFields.ResourceId) {
					edge = opengraph.GenerateEdge(edgeKind, team.OID, jobTemplates[role.SummaryFields.ResourceId].OID)
				}
			}
			if edge != nil {
				opengraph.AddEdge(&graph, edge)
			}
		}
	}

	// -- Link Ansible and Active Directory --

	if (ldap != core.AHLdap{}) {

		log.Info("Linking LDAP Users.")
		edgeKind = "SyncedToAHUser"

		conn, err := core.Connect(ldap)

		if err != nil {
			log.Fatal("Connection failed:", err)
		}
		defer conn.Close()

		for _, user := range users {
			if user.ExternalAccount == core.LDAP_VALUE {
				if user.LdapDn != "" {
					ldap_dn := user.LdapDn
					objectSid, _ := core.Search(conn, ldap_dn)
					edge := opengraph.GenerateEdge(edgeKind, objectSid, user.OID)
					graph.AddEdgeWithoutValidation(edge)
				}
			}
		}
	} else {
		log.Warn("Skipping linking LDAP Users since the user used for authentication is a local user")
	}

	// -- Output final graph --

	outputJson, err := graph.ExportJSON(false)
	err = opengraph.WriteToFile(
		[]byte(outputJson),
		path.Join(outdir, opengraph.CalculateName("output")))
	if err != nil {
		log.Error(err)
	}

}

var ingestCmd = &cobra.Command{
	Use:   "collect",
	Short: "Go collector for adding Ansible WorX and Ansible Tower attack paths to BloodHound with OpenGraph ",
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
		password, _ := cmd.Flags().GetString("password")
		token, _ := cmd.Flags().GetString("token")
		if (username == "" || password == "") && token == "" {
			log.Fatal("Invalid authentication material provided.")
		}

		dc_ipAddress, _ := cmd.Flags().GetString("dc-ip")
		domain, _ := cmd.Flags().GetString("domain")
		isLDAPS, _ := cmd.Flags().GetBool("ldaps")

		if (token != "") && (dc_ipAddress != "" || domain != "") {
			log.Warn("Domain and domain controller address will not be taken into account using token authentication material")
		}

		if (token == "") && ((dc_ipAddress == "" && domain != "") || (dc_ipAddress != "" && domain == "")) {
			log.Fatal("Invalid domain name or IP provided")
		}

		if dc_ipAddress == "" && domain == "" && isLDAPS {
			log.Warn("The LDAPS parameter will not be taken into account since the domain and the domain controller address are not configured")
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

		client := core.InitClient(proxyURL, skipVerifySSL, username, password, token)

		var ldap core.AHLdap

		if dc_ipAddress != "" && domain != "" {
			ldap = core.InitLdap(dc_ipAddress, username, password, domain, isLDAPS, skipVerifySSL)
		}

		launchGathering(client, targetUrl, outdir, ldap)
	},
}

func main() {

	ingestCmd.Flags().StringP("target", "t", "", "Target URL of AWX/Tower instance.")
	_ = ingestCmd.MarkFlagRequired("target")

	ingestCmd.Flags().StringP("username", "u", "", "Username to use for authentication.")
	ingestCmd.Flags().StringP("token", "", "", "Token to use for authentication.")
	ingestCmd.Flags().StringP("password", "p", "", "Password to use for authentication.")

	ingestCmd.Flags().StringP("dc-ip", "", "", "(optional) Target IP of the domain. Required only for LDAP user")
	ingestCmd.Flags().BoolP("ldaps", "", false, "(optional) Configure LDAPS authentication on the domain controller. Required only for LDAP user")
	ingestCmd.Flags().StringP("domain", "", "", "(optional) NetBIOS domain name. Required only for LDAP user")

	ingestCmd.Flags().StringP("proxy", "", "", "(optional) Configure HTTP/HTTPS proxy.")
	ingestCmd.Flags().StringP("outdir", "", "", "(optional) Output directory for the json files.")
	ingestCmd.Flags().BoolP("verbose", "v", false, "(optional) Enable debug logs.")
	ingestCmd.Flags().BoolP("skip-verify-ssl", "k", false, "(optional) Skips SSL/TLS verification for HTTP and LDAP.")

	if err := ingestCmd.Execute(); err != nil {
		msg := fmt.Sprintf("CLI error: %v\n", err)
		log.Error(msg)
		os.Exit(1)
	}

}
