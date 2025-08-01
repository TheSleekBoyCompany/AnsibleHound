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

// TODO: Implement some map[int]AnsibleType object, allowing to instantly find the right ID.

func launchGathering(client http.Client, targetUrl *url.URL,
	token string, outdir string) {

	output := core.OutputJson{
		Metadata: core.Metadata{
			SourceKind: "AnsibleBase",
		},
	}
	nodes := []core.Node{}
	edges := []core.Edge{}

	// -- Gathering all nodes --

	log.Info("Gathering Users.")
	users, err := core.GatherObject[*core.User](
		client, *targetUrl, token, core.USERS_ENDPOINT,
	)
	if err != nil {
		log.Error("An error occured while gathering Users, skipping.")
		log.Error(err)
	}

	log.Info("Gathering User Roles.")
	for i, user := range users {
		userRolesEndpoint := fmt.Sprintf(core.USER_ROLES_ENDPOINT, user.ID)
		roles, err := core.GatherObject[*core.Role](client, *targetUrl, token, userRolesEndpoint)
		if err != nil {
			log.Error("An error occured while gathering Team Roles.")
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
		client, *targetUrl, token, core.HOSTS_ENDPOINT,
	)
	if err != nil {
		log.Error("An error occured while gathering Users, skipping.")
		log.Error(err)
	}
	hostNodes := core.GenerateNodes(hosts)
	nodes = append(nodes, hostNodes...)

	log.Info("Gathering Jobs.")
	jobs, err := core.GatherObject[*core.Job](
		client, *targetUrl, token, core.JOBS_ENDPOINT,
	)
	if err != nil {
		log.Error("An error occured while gathering jobs, skipping.")
		log.Error(err)
	}
	jobsNodes := core.GenerateNodes(jobs)
	nodes = append(nodes, jobsNodes...)

	log.Info("Gathering Job Templates.")
	jobTemplates, err := core.GatherObject[*core.JobTemplate](
		client, *targetUrl, token, core.JOB_TEMPLATE_ENDPOINT,
	)
	if err != nil {
		log.Error("An error occured while gathering jobs, skipping.")
		log.Error(err)
	}
	jobTemplatesNodes := core.GenerateNodes(jobTemplates)
	nodes = append(nodes, jobTemplatesNodes...)

	log.Info("Gathering Inventories.")
	inventories, err := core.GatherObject[*core.Inventory](
		client, *targetUrl, token, core.INVENTORIES_ENDPOINT,
	)
	if err != nil {
		log.Error("An error occured while gathering jobs, skipping.")
		log.Error(err)
	}
	inventoriesNodes := core.GenerateNodes(inventories)
	nodes = append(nodes, inventoriesNodes...)

	log.Info("Gathering Organizations.")
	organizations, err := core.GatherObject[*core.Organization](
		client, *targetUrl, token, core.ORGANIZATIONS_ENDPOINT,
	)
	if err != nil {
		log.Error("An error occured while gathering jobs, skipping.")
		log.Error(err)
	}
	organizationNodes := core.GenerateNodes(organizations)
	nodes = append(nodes, organizationNodes...)

	log.Info("Gathering Credentials.")
	credentials, err := core.GatherObject[*core.Credential](
		client, *targetUrl, token, core.CREDENTIALS_ENDPOINT,
	)
	if err != nil {
		log.Error("An error occured while gathering jobs, skipping.")
		log.Error(err)
	}
	credentialNodes := core.GenerateNodes(credentials)
	nodes = append(nodes, credentialNodes...)

	log.Info("Gathering Projects.")
	projects, err := core.GatherObject[*core.Project](
		client, *targetUrl, token, core.PROJECTS_ENDPOINT,
	)
	if err != nil {
		log.Error("An error occured while gathering jobs, skipping.")
		log.Error(err)
	}
	projectNodes := core.GenerateNodes(projects)
	nodes = append(nodes, projectNodes...)

	log.Info("Gathering Teams.")
	teams, err := core.GatherObject[*core.Team](
		client, *targetUrl, token, core.TEAMS_ENDPOINT,
	)
	if err != nil {
		log.Error("An error occured while gathering Teams, skipping.")
		log.Error(err)
	}

	log.Info("Gathering Team Roles.")
	for i, team := range teams {

		teamRolesEndpoint := fmt.Sprintf(core.TEAM_ROLES_ENDPOINT, team.ID)
		roles, err := core.GatherObject[*core.Role](
			client, *targetUrl, token, teamRolesEndpoint,
		)
		if err != nil {
			log.Error("An error occured while gathering Team Roles.")
			log.Error(err)
			continue
		}

		teamMembersEndpoint := fmt.Sprintf(core.TEAM_USERS_ENDPOINT, team.ID)
		members, err := core.GatherObject[*core.User](
			client, *targetUrl, token, teamMembersEndpoint,
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

	log.Info("Linking Organizations and Inventories.")
	kind := "ATContains"
	for _, inventory := range inventories {
		edge := core.GenerateEdge(kind, organizations[inventory.Organization].UUID, inventory.UUID)
		edges = append(edges, edge)
	}

	log.Info("Linking Inventories and Hosts.")
	kind = "ATContains"
	for _, host := range hosts {
		edge := core.GenerateEdge(kind, inventories[host.Inventory].UUID, host.UUID)
		edges = append(edges, edge)
	}

	log.Info("Linking Job Templates and Jobs.")
	kind = "ATContains"
	for _, job := range jobs {
		edge := core.GenerateEdge(kind, jobTemplates[job.UnifiedJobTemplate].UUID, job.UUID)
		edges = append(edges, edge)
	}

	log.Info("Linking Organizations and Job Templates.")
	kind = "ATContains"
	for _, jobTemplate := range jobTemplates {
		edge := core.GenerateEdge(kind, organizations[jobTemplate.Organization].UUID, jobTemplate.UUID)
		edges = append(edges, edge)
	}

	log.Info("Linking Organizations and Credentials.")
	kind = "ATContains"
	for _, credential := range credentials {
		if credential.Organization != 0 {
			edge := core.GenerateEdge(kind, organizations[credential.Organization].UUID, credential.UUID)
			edges = append(edges, edge)
		}
	}

	log.Info("Linking Organizations and Projects")
	kind = "ATContains"
	for _, project := range projects {
		edge := core.GenerateEdge(kind, organizations[project.Organization].UUID, project.UUID)
		edges = append(edges, edge)
	}

	log.Info("Linking Job Templates and Projects.")
	kind = "ATUses"
	for _, jobTemplate := range jobTemplates {
		edge := core.GenerateEdge(kind, jobTemplate.UUID, projects[jobTemplate.Project].UUID)
		edges = append(edges, edge)
	}

	log.Info("Linking Job Template and Inventories.")
	kind = "ATUses"
	for _, jobTemplate := range jobTemplates {
		edge := core.GenerateEdge(kind, jobTemplate.UUID, inventories[jobTemplate.Inventory].UUID)
		edges = append(edges, edge)
	}

	log.Info("Linking User Roles.")
	for _, user := range users {

		for _, role := range user.Roles {
			kind := "AT" + strings.ReplaceAll(role.Name, " ", "")

			switch role.SummaryFields.ResourceType {

			case "organization":
				edge := core.GenerateEdge(kind, user.UUID, organizations[role.SummaryFields.ResourceId].UUID)
				edges = append(edges, edge)
			case "inventory":
				edge := core.GenerateEdge(kind, user.UUID, inventories[role.SummaryFields.ResourceId].UUID)
				edges = append(edges, edge)
			case "team":
				edge := core.GenerateEdge(kind, user.UUID, teams[role.SummaryFields.ResourceId].UUID)
				edges = append(edges, edge)
			case "credential":
				edge := core.GenerateEdge(kind, user.UUID, credentials[role.SummaryFields.ResourceId].UUID)
				edges = append(edges, edge)
			case "job_template":
				edge := core.GenerateEdge(kind, user.UUID, jobTemplates[role.SummaryFields.ResourceId].UUID)
				edges = append(edges, edge)
			}
		}
	}

	log.Info("Linking Team Roles.")
	for _, team := range teams {
		for _, role := range team.Roles {
			kind := "AT" + strings.ReplaceAll(role.Name, " ", "")

			switch role.SummaryFields.ResourceType {

			case "organization":
				edge := core.GenerateEdge(kind, team.UUID, organizations[role.SummaryFields.ResourceId].UUID)
				edges = append(edges, edge)
			case "inventory":
				edge := core.GenerateEdge(kind, team.UUID, inventories[role.SummaryFields.ResourceId].UUID)
				edges = append(edges, edge)
			case "team":
				edge := core.GenerateEdge(kind, team.UUID, teams[role.SummaryFields.ResourceId].UUID)
				edges = append(edges, edge)
			case "credential":
				edge := core.GenerateEdge(kind, team.UUID, credentials[role.SummaryFields.ResourceId].UUID)
				edges = append(edges, edge)
			case "job_template":
				edge := core.GenerateEdge(kind, team.UUID, jobTemplates[role.SummaryFields.ResourceId].UUID)
				edges = append(edges, edge)
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

		outdir, _ := cmd.Flags().GetString("outdir")

		token, _ := cmd.Flags().GetString("token")
		if token == "" {
			log.Error("Empty token provided.")
			os.Exit(1)
		}

		verbose, _ := cmd.Flags().GetBool("verbose")
		if verbose {
			log.SetLevel(log.DebugLevel)
		}

		target, _ := cmd.Flags().GetString("target")
		if target == "" {
			log.Error("Empty target provided.")
			os.Exit(1)
		}

		targetUrl, err := url.Parse(target)
		if err != nil {
			log.Error("Invalid target URL specified.\n%s", err)
			os.Exit(1)
		}

		var proxyURL *url.URL
		proxy, _ := cmd.Flags().GetString("proxy")
		if proxy != "" {
			proxyURL, err = url.Parse(proxy)
			if err != nil {
				log.Errorf("Invalid proxy URL specified.\n%s", err)
				os.Exit(1)
			}
		}

		client := core.InitClient(proxyURL)

		launchGathering(client, targetUrl, token, outdir)

	},
}

func main() {

	ingestCmd.Flags().StringP("target", "u", "", "Target URL of AWX/Tower instance.")
	_ = ingestCmd.MarkFlagRequired("target")

	ingestCmd.Flags().StringP("token", "t", "", "Token to use for authentication.")
	_ = ingestCmd.MarkFlagRequired("token")

	ingestCmd.Flags().StringP("proxy", "p", "", "(optional) Configure HTTP/HTTPS proxy.")
	ingestCmd.Flags().StringP("outdir", "", "", "(optional) Output directory for the json files.")
	ingestCmd.Flags().BoolP("verbose", "v", false, "(optional) Enable debug logs.")

	if err := ingestCmd.Execute(); err != nil {
		msg := fmt.Sprintf("CLI error: %v\n", err)
		log.Error(msg)
		os.Exit(1)
	}

}
