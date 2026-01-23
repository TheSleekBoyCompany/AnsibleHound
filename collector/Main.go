package main

import (
	"fmt"
	"net/url"
	"os"

	"ansible-hound/core/gather"
	"ansible-hound/core/opengraph"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

func launch(client gather.AHClient, targetUrl *url.URL,
	outdir string, ldap gather.AHLdap) {

	graph := opengraph.InitGraph()

	// -- Check if credentials are valid --

	log.Info("Authenticating on Ansible Worx/Tower instance.")
	err := gather.ValidateCredentials(client, *targetUrl)
	if err != nil {
		log.Fatal("Unable to authenticate on Ansible WorX/Tower instance due to invalid credentials.")
	}

	// -- Gathering Ansible Instance information --

	log.Info("Gathering Ansible Worx/Tower instance information.")
	instance, err := gather.GatherAnsibleInstance(client, *targetUrl)
	if err != nil {
		log.Fatalf("Unable to gather Ansible WorX/Tower information (%s).", targetUrl)
	}
	instance.Name = targetUrl.Host
	instanceNode := instance.ToBHNode()
	graph.AddNode(instanceNode)

	// -- Gathering all nodes --

	users, err := gather.GatherUsers(client, instance.InstallUUID, *targetUrl)
	if err == nil {
		userNodes := opengraph.GenerateNodes(users)
		opengraph.AddNodes(&graph, userNodes)
	}

	hosts, err := gather.GatherHosts(client, instance.InstallUUID, *targetUrl)
	if err == nil {
		hostNodes := opengraph.GenerateNodes(hosts)
		opengraph.AddNodes(&graph, hostNodes)
	}

	groups, err := gather.GatherGroups(client, instance.InstallUUID, *targetUrl)
	if err == nil {
		groupsNodes := opengraph.GenerateNodes(groups)
		opengraph.AddNodes(&graph, groupsNodes)
	}

	jobs, err := gather.GatherJobs(client, instance.InstallUUID, *targetUrl)
	if err == nil {
		jobsNodes := opengraph.GenerateNodes(jobs)
		opengraph.AddNodes(&graph, jobsNodes)
	}

	jobTemplates, err := gather.GatherJobTemplates(client, instance.InstallUUID, *targetUrl)
	if err == nil {
		jobTemplatesNodes := opengraph.GenerateNodes(jobTemplates)
		opengraph.AddNodes(&graph, jobTemplatesNodes)
	}

	inventories, err := gather.GatherInventories(client, instance.InstallUUID, *targetUrl)
	if err == nil {
		inventoriesNodes := opengraph.GenerateNodes(inventories)
		opengraph.AddNodes(&graph, inventoriesNodes)
	}

	organizations, err := gather.GatherOrganizations(client, instance.InstallUUID, *targetUrl)
	if err == nil {
		organizationNodes := opengraph.GenerateNodes(organizations)
		opengraph.AddNodes(&graph, organizationNodes)
	}

	credentials, err := gather.GatherCredentials(client, instance.InstallUUID, *targetUrl)
	if err == nil {
		credentialNodes := opengraph.GenerateNodes(credentials)
		opengraph.AddNodes(&graph, credentialNodes)
	}

	credentialTypes, err := gather.GatherCredentialTypes(client, instance.InstallUUID, *targetUrl)
	if err == nil {
		credentialTypesNodes := opengraph.GenerateNodes(credentialTypes)
		opengraph.AddNodes(&graph, credentialTypesNodes)
	}

	projects, err := gather.GatherProjects(client, instance.InstallUUID, *targetUrl)
	if err == nil {
		projectNodes := opengraph.GenerateNodes(projects)
		opengraph.AddNodes(&graph, projectNodes)
	}

	teams, err := gather.GatherTeams(client, instance.InstallUUID, *targetUrl)
	if err == nil {
		teamNodes := opengraph.GenerateNodes(teams)
		opengraph.AddNodes(&graph, teamNodes)
	}

	// -- Creating all edges --

	opengraph.LinkOrganization(&graph,
		instance.OID, organizations, inventories,
		jobTemplates, credentials, projects)

	opengraph.LinkInventory(&graph, inventories,
		hosts, groups)

	opengraph.LinkJobTemplates(&graph, jobTemplates, jobs,
		projects, inventories, credentials, credentialTypes)

	opengraph.LinkUserRoles(&graph, users, organizations,
		inventories, teams, credentials,
		jobTemplates)

	opengraph.LinkTeamRoles(&graph, users, organizations,
		inventories, teams, credentials,
		jobTemplates)

	opengraph.LinkAdministrativeRights(&graph, users, jobTemplates, credentials,
		inventories, projects, organizations, teams)

	// -- Link Ansible and Active Directory --

	opengraph.LinkAD(&graph, ldap, users)

	// -- Output final graph --

	opengraph.Output(&graph, outdir)

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

		client := gather.InitClient(proxyURL, skipVerifySSL, username, password, token)

		var ldap gather.AHLdap

		if dc_ipAddress != "" && domain != "" {
			ldap = gather.InitLdap(dc_ipAddress, username, password, domain, isLDAPS, skipVerifySSL)
		}

		launch(client, targetUrl, outdir, ldap)
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
