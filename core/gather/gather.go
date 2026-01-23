package gather

import (
	"ansible-hound/core/ansible"
	"fmt"
	"net/url"

	"github.com/charmbracelet/log"
)

func ValidateCredentials(client AHClient, targetUrl url.URL) (err error) {
	_, err = AuthenticateOnAnsibleInstance(client, targetUrl, ME_ENDPOINT)
	return err
}

func GatherUsers(client AHClient, installUUID string,
	targetUrl url.URL) (users map[int]*ansible.User, err error) {

	log.Info("Gathering Users.")
	users, err = GatherObject[*ansible.User](
		installUUID, client, targetUrl, USERS_ENDPOINT,
	)
	if err != nil {
		log.Error("An error occured while gathering Users, skipping.")
		log.Error(err)
		return nil, err
	}

	log.Info("Gathering User Roles.")
	for i, user := range users {
		userRolesEndpoint := fmt.Sprintf(USER_ROLES_ENDPOINT, user.ID)
		roles, err := GatherObject[*ansible.Role](
			installUUID, client, targetUrl, userRolesEndpoint)
		if err != nil {
			log.Error("An error occured while gathering User Roles.")
			log.Error(err)
			continue
		}
		user.Roles = roles
		users[i] = user
	}
	return users, err
}

func GatherHosts(client AHClient, installUUID string,
	targetUrl url.URL) (hosts map[int]*ansible.Host, err error) {

	log.Info("Gathering Hosts.")
	hosts, err = GatherObject[*ansible.Host](
		installUUID, client, targetUrl, HOSTS_ENDPOINT,
	)
	if err != nil {
		log.Error("An error occured while gathering Hosts, skipping.")
		log.Error(err)
	}
	return hosts, err
}

func GatherGroups(client AHClient, installUUID string,
	targetUrl url.URL) (groups map[int]*ansible.Group, err error) {

	log.Info("Gathering Groups.")
	groups, err = GatherObject[*ansible.Group](
		installUUID, client, targetUrl, GROUPS_ENDPOINT,
	)
	if err != nil {
		log.Error("An error occured while gathering Groups, skipping.")
		log.Error(err)
	}
	log.Info("Gathering Group Hosts.")
	for i, group := range groups {

		groupHostsEndpoint := fmt.Sprintf(GROUP_HOSTS_ENDPOINT, group.ID)
		hosts, err := GatherObject[*ansible.Host](
			installUUID, client, targetUrl, groupHostsEndpoint,
		)
		if err != nil {
			log.Error("An error occured while gathering Group Hosts, skipping Group.")
			log.Error(err)
			continue
		}
		group.Hosts = hosts
		groups[i] = group
	}
	return groups, err
}

func GatherJobs(client AHClient, installUUID string,
	targetUrl url.URL) (jobs map[int]*ansible.Job, err error) {

	log.Info("Gathering Jobs.")
	jobs, err = GatherObject[*ansible.Job](
		installUUID, client, targetUrl, JOBS_ENDPOINT,
	)
	if err != nil {
		log.Error("An error occured while gathering Jobs, skipping.")
		log.Error(err)
	}
	return jobs, err
}

func GatherJobTemplates(client AHClient, installUUID string,
	targetUrl url.URL) (jobTemplates map[int]*ansible.JobTemplate, err error) {

	log.Info("Gathering Job Templates.")
	jobTemplates, err = GatherObject[*ansible.JobTemplate](
		installUUID, client, targetUrl, JOB_TEMPLATE_ENDPOINT,
	)
	if err != nil {
		log.Error("An error occured while gathering Job Templates, skipping.")
		log.Error(err)
	}

	log.Info("Gathering Job Templates Credentials.")
	for i, jobTemplate := range jobTemplates {

		jobTemplatesCredentialsEndpoint := fmt.Sprintf(
			JOB_TEMPLATE_CREDENTIALS_ENDPOINT, jobTemplate.ID)
		credentials, err := GatherObject[*ansible.Credential](
			installUUID, client, targetUrl, jobTemplatesCredentialsEndpoint,
		)
		if err != nil {
			log.Error("An error occured while gathering Job Template Credentials.")
			log.Error(err)
			continue
		}

		jobTemplate.Credentials = credentials
		jobTemplates[i] = jobTemplate
	}

	return jobTemplates, err

}

func GatherInventories(client AHClient, installUUID string,
	targetUrl url.URL) (inventories map[int]*ansible.Inventory, err error) {

	log.Info("Gathering Inventories.")
	inventories, err = GatherObject[*ansible.Inventory](
		installUUID, client, targetUrl, INVENTORIES_ENDPOINT,
	)
	if err != nil {
		log.Error("An error occured while gathering Inventories, skipping.")
		log.Error(err)
	}

	return inventories, err
}

func GatherOrganizations(client AHClient, installUUID string,
	targetUrl url.URL) (organizations map[int]*ansible.Organization, err error) {

	log.Info("Gathering Organizations.")
	organizations, err = GatherObject[*ansible.Organization](
		installUUID, client, targetUrl, ORGANIZATIONS_ENDPOINT,
	)
	if err != nil {
		log.Error("An error occured while gathering Organizations, skipping.")
		log.Error(err)
	}

	return organizations, err

}

func GatherCredentials(client AHClient, installUUID string,
	targetUrl url.URL) (credentials map[int]*ansible.Credential, err error) {

	log.Info("Gathering Credentials.")
	credentials, err = GatherObject[*ansible.Credential](
		installUUID, client, targetUrl, CREDENTIALS_ENDPOINT,
	)
	if err != nil {
		log.Error("An error occured while gathering Credentials, skipping.")
		log.Error(err)
	}
	return credentials, err
}

func GatherCredentialTypes(client AHClient, installUUID string,
	targetUrl url.URL) (credentialTypes map[int]*ansible.CredentialType, err error) {

	log.Info("Gathering Credential Types.")
	credentialTypes, err = GatherObject[*ansible.CredentialType](
		installUUID, client, targetUrl, CREDENTIAL_TYPES_ENDPOINT,
	)
	if err != nil {
		log.Error("An error occured while gathering Credential Types, skipping.")
		log.Error(err)
	}

	return credentialTypes, err
}

func GatherProjects(client AHClient, installUUID string,
	targetUrl url.URL) (projects map[int]*ansible.Project, err error) {

	log.Info("Gathering Projects.")
	projects, err = GatherObject[*ansible.Project](
		installUUID, client, targetUrl, PROJECTS_ENDPOINT,
	)
	if err != nil {
		log.Error("An error occured while gathering Projects, skipping.")
		log.Error(err)
	}

	return projects, err
}

func GatherTeams(client AHClient, installUUID string,
	targetUrl url.URL) (teams map[int]*ansible.Team, err error) {

	log.Info("Gathering Teams.")
	teams, err = GatherObject[*ansible.Team](
		installUUID, client, targetUrl, TEAMS_ENDPOINT,
	)
	if err != nil {
		log.Error("An error occured while gathering Teams, skipping.")
		log.Error(err)
	}

	log.Info("Gathering Team Roles.")
	for i, team := range teams {

		teamRolesEndpoint := fmt.Sprintf(TEAM_ROLES_ENDPOINT, team.ID)
		roles, err := GatherObject[*ansible.Role](
			installUUID, client, targetUrl, teamRolesEndpoint,
		)
		if err != nil {
			log.Error("An error occured while gathering Team Roles.")
			log.Error(err)
			continue
		}

		teamMembersEndpoint := fmt.Sprintf(TEAM_USERS_ENDPOINT, team.ID)
		members, err := GatherObject[*ansible.User](
			installUUID, client, targetUrl, teamMembersEndpoint,
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

	return teams, err
}
