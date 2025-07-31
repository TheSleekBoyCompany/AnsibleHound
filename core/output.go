package core

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
)

func OutputBH_Edge(kind string, startId string, endId string) Edge {

	start := Link{
		Value:   startId,
		MatchBy: "id",
	}

	end := Link{
		Value:   endId,
		MatchBy: "id",
	}

	edge := Edge{
		Kind:  kind,
		Start: start,
		End:   end,
	}

	return edge
}

// TODO: Must be a cleaner way than this big switch case.
// PRs WELCOME *wink* *wink*

func OutputBH_Node[T AnsibleTypeList](objectLists T) []Node {

	nodes := []Node{}

	switch objects := any(objectLists).(type) {
	case []User:
		for i, user := range objects {
			node := Node{}
			node.Kinds = []string{
				"ATUser",
			}
			user.UUID = uuid.NewString()
			node.Id = user.UUID
			objects[i] = user
			node.Properties = map[string]string{
				"id":                strconv.Itoa(user.ID),
				"name":              user.Username,
				"description":       user.Description,
				"url":               user.Url,
				"firstname":         user.FirstName,
				"lastname":          user.LastName,
				"email":             user.Email,
				"is_super_user":     strconv.FormatBool(user.IsSuperUser),
				"is_system_auditor": strconv.FormatBool(user.IsSystemAuditor),
				"ldap_dn":           user.LdapDn,
				"last_login":        user.LastLogin,
				"external_account":  user.ExternalAccount,
			}
			nodes = append(nodes, node)
		}

	case []Organization:
		for i, org := range objects {
			node := Node{}
			node.Kinds = []string{
				"ATOrganization",
			}
			org.UUID = uuid.NewString()
			node.Id = org.UUID
			objects[i] = org
			node.Properties = map[string]string{
				"id":                  strconv.Itoa(org.ID),
				"name":                org.Name,
				"description":         org.Description,
				"url":                 org.Url,
				"max_hosts":           strconv.FormatInt(int64(org.MaxHosts), 10),
				"custom_virtualenv":   org.CustomVirtualenv,
				"default_environment": strconv.FormatInt(int64(org.DefaultEnvironment), 10),
			}
			nodes = append(nodes, node)
		}

	case []JobTemplate:
		for i, jobTemplate := range objects {
			node := Node{}
			node.Kinds = []string{
				"ATJobTemplate",
			}
			jobTemplate.UUID = uuid.NewString()
			node.Id = jobTemplate.UUID
			objects[i] = jobTemplate
			node.Properties = map[string]string{
				"id":                                  strconv.Itoa(jobTemplate.ID),
				"name":                                jobTemplate.Name,
				"description":                         jobTemplate.Description,
				"url":                                 jobTemplate.Url,
				"job_type":                            jobTemplate.JobType,
				"project":                             strconv.FormatInt(int64(jobTemplate.Project), 10),
				"organization":                        strconv.FormatInt(int64(jobTemplate.Organization), 10),
				"playbook":                            jobTemplate.Playbook,
				"scm_branch":                          jobTemplate.SCMBranch,
				"limit":                               jobTemplate.Limit,
				"verbosity":                           strconv.FormatInt(int64(jobTemplate.Verbosity), 10),
				"extra_vars":                          jobTemplate.ExtraVars,
				"status":                              jobTemplate.Status,
				"job_tags":                            jobTemplate.JobTags,
				"forks":                               strconv.FormatInt(int64(jobTemplate.Forks), 10),
				"skip_tags":                           jobTemplate.SkipTags,
				"start_at_task":                       jobTemplate.StartAtTask,
				"timeout":                             strconv.FormatInt(int64(jobTemplate.Timeout), 10),
				"use_fact_cache":                      strconv.FormatBool(jobTemplate.UseFactCache),
				"force_handler":                       strconv.FormatBool(jobTemplate.ForceHandler),
				"last_job_run":                        jobTemplate.LastJobRun,
				"next_job_run":                        jobTemplate.NextJobRun,
				"last_job_failed":                     strconv.FormatBool(jobTemplate.LastJobFailed),
				"execution_environment":               strconv.FormatInt(int64(jobTemplate.ExecutionEnvironment), 10),
				"host_config_key":                     jobTemplate.HostConfigKey,
				"ask_scm_branch_on_launch":            strconv.FormatBool(jobTemplate.AskScmBranchOnLaunch),
				"ask_diff_mode_on_launch":             strconv.FormatBool(jobTemplate.AskDiffModeOnLaunch),
				"ask_variables_on_launch":             strconv.FormatBool(jobTemplate.AskVariablesOnLaunch),
				"ask_limit_on_launch":                 strconv.FormatBool(jobTemplate.AskLimitOnLaunch),
				"ask_tags_on_launch":                  strconv.FormatBool(jobTemplate.AskTagsOnLaunch),
				"ask_job_type_on_launch":              strconv.FormatBool(jobTemplate.AskJobTypeOnLaunch),
				"ask_verbosity_on_launch":             strconv.FormatBool(jobTemplate.AskVerbosityOnLaunch),
				"ask_inventory_on_launch":             strconv.FormatBool(jobTemplate.AskInventoryOnLaunch),
				"ask_credential_on_launch":            strconv.FormatBool(jobTemplate.AskCredentialOnLaunch),
				"ask_execution_environment_on_launch": strconv.FormatBool(jobTemplate.AskExecutionEnvironmentOnLaunch),
				"ask_labels_on_launch":                strconv.FormatBool(jobTemplate.AskLabelsOnLaunch),
				"ask_forks_on_launch":                 strconv.FormatBool(jobTemplate.AskForksOnLaunch),
				"ask_job_slice_count_on_launch":       strconv.FormatBool(jobTemplate.AskJobSliceCountOnLaunch),
				"ask_timeout_on_launch":               strconv.FormatBool(jobTemplate.AskTimeoutOnLaunch),
				"ask_instance_groups_on_launch":       strconv.FormatBool(jobTemplate.AskInstanceGroupsOnLaunch),
				"survey_enabled":                      strconv.FormatBool(jobTemplate.SurveyEnabled),
				"become_enabled":                      strconv.FormatBool(jobTemplate.BecomeEnabled),
				"diff_mode":                           strconv.FormatBool(jobTemplate.DiffMode),
				"allow_simultaneous":                  strconv.FormatBool(jobTemplate.AllowSimultaneous),
				"custom_virtualenv":                   jobTemplate.CustomVirtualenv,
				"job_slice_count":                     strconv.FormatInt(int64(jobTemplate.JobSliceCount), 10),
				"webhook_service":                     jobTemplate.WebhookService,
				"webhook_credential":                  strconv.FormatInt(int64(jobTemplate.WebhookCredential), 10),
				"prevent_instance_group_fallback":     strconv.FormatBool(jobTemplate.PreventInstanceGroupFallback),
			}
			nodes = append(nodes, node)
		}

	case []Job:
		for i, job := range objects {
			node := Node{}
			node.Kinds = []string{
				"ATJob",
			}
			job.UUID = uuid.NewString()
			node.Id = job.UUID
			objects[i] = job
			node.Properties = map[string]string{
				"id":                       strconv.FormatInt(int64(job.ID), 10),
				"name":                     job.Name,
				"description":              job.Description,
				"url":                      job.Url,
				"type":                     job.Type,
				"created":                  job.Created,
				"modified":                 job.Modified,
				"playbook":                 job.Playbook,
				"scm_branch":               job.ScmBranch,
				"forks":                    strconv.FormatInt(int64(job.Forks), 10),
				"limit":                    job.Limit,
				"verbosity":                strconv.FormatInt(int64(job.Verbosity), 10),
				"extra_vars":               job.ExtraVars,
				"started":                  job.Started,
				"finished":                 job.Finished,
				"canceled_on":              job.CanceledOn,
				"elapsed":                  strconv.FormatFloat(float64(job.Elapsed), 'e', 10, 32),
				"job_explanation":          job.JobExplanation,
				"launch_type":              job.LaunchType,
				"unified_job_template":     strconv.FormatInt(int64(job.UnifiedJobTemplate), 10),
				"organization":             strconv.FormatInt(int64(job.Organization), 10),
				"inventory":                strconv.FormatInt(int64(job.Inventory), 10),
				"project":                  strconv.FormatInt(int64(job.Project), 10),
				"failed":                   strconv.FormatBool(job.Failed),
				"status":                   job.Status,
				"execution_environment":    strconv.FormatInt(int64(job.ExecutionEnvironment), 10),
				"execution_node":           job.ExecutionNode,
				"controller_node":          job.ControllerNode,
				"work_unit_id":             job.WorkUnitId,
				"job_tags":                 job.JobTags,
				"job_type":                 job.JobType,
				"force_handler":            strconv.FormatBool(job.ForceHandler),
				"skip_tags":                job.SkipTags,
				"start_at_task":            job.StartAtTask,
				"timeout":                  strconv.FormatInt(int64(job.Timeout), 10),
				"use_fact_cache":           strconv.FormatBool(job.UseFactCache),
				"password_needed_to_start": job.PasswordNeededToStart,
				"allow_simultaneous":       strconv.FormatBool(job.AllowSimultaneous),
				"scm_revision":             job.ScmRevision,
				"instance_group":           strconv.FormatInt(int64(job.InstanceGroup), 10),
				"diff_mode":                strconv.FormatBool(job.DiffMode),
				"job_slice_number":         strconv.FormatInt(int64(job.JobSliceNumber), 10),
				"job_slice_count":          strconv.FormatInt(int64(job.JobSliceCount), 10),
				"webhook_guid":             job.WebhookGuid,
				"webhook_service":          job.WebhookService,
				"webhook_credential":       strconv.FormatInt(int64(job.WebhookCredential), 10),
			}
			nodes = append(nodes, node)
		}

	case []Project:
		for i, project := range objects {
			node := Node{}
			node.Kinds = []string{
				"ATProject",
			}
			project.UUID = uuid.NewString()
			node.Id = project.UUID
			objects[i] = project
			node.Properties = map[string]string{
				"id":                              strconv.Itoa(project.ID),
				"name":                            project.Name,
				"description":                     project.Description,
				"url":                             project.Url,
				"organization":                    strconv.FormatInt(int64(project.Organization), 10),
				"credential":                      strconv.FormatInt(int64(project.Credential), 10),
				"timeout":                         strconv.FormatInt(int64(project.Timeout), 10),
				"status":                          project.Status,
				"local_path":                      project.LocalPath,
				"scm_type":                        project.ScmType,
				"scm_url":                         project.ScmUrl,
				"scm_branch":                      project.ScmBranch,
				"scm_ref_spec":                    project.ScmRefSpec,
				"scm_clean":                       strconv.FormatBool(project.ScmClean),
				"scm_track_submodules":            strconv.FormatBool(project.ScmTrackSubmodules),
				"scm_delete_on_update":            strconv.FormatBool(project.ScmDeleteOnUpdate),
				"scm_revision":                    project.ScmRevision,
				"last_job_run":                    project.LastJobRun,
				"next_job_run":                    project.NextJobRun,
				"last_job_failed":                 strconv.FormatBool(project.LastJobFailed),
				"scm_update_on_launch":            strconv.FormatBool(project.ScmUpdateOnLaunch),
				"scm_update_cache_timeout":        strconv.FormatInt(int64(project.ScmUpdateCacheTimeout), 10),
				"allow_override":                  strconv.FormatBool(project.AllowOverride),
				"custom_virtualenv":               project.CustomVirtualenv,
				"default_environment":             strconv.FormatInt(int64(project.DefaultEnvironment), 10),
				"signature_validation_credential": strconv.FormatInt(int64(project.SignatureValidationCredential), 10),
				"last_update_failed":              strconv.FormatBool(project.LastUpdateFailed),
				"last_update":                     project.LastUpdate,
			}
			nodes = append(nodes, node)
		}

	case []Credential:
		for i, credential := range objects {
			node := Node{}
			node.Kinds = []string{
				"ATCredential",
			}
			credential.UUID = uuid.NewString()
			node.Id = credential.UUID
			objects[i] = credential
			node.Properties = map[string]string{
				"id":              strconv.Itoa(credential.ID),
				"name":            credential.Name,
				"description":     credential.Description,
				"url":             credential.Url,
				"organization":    strconv.FormatInt(int64(credential.Organization), 10),
				"credential_type": strconv.FormatInt(int64(credential.CredentialType), 10),
				"managed":         strconv.FormatBool(credential.Managed),
				"cloud":           strconv.FormatBool(credential.Cloud),
				"kubernetes":      strconv.FormatBool(credential.Kubernetes),
			}
			nodes = append(nodes, node)
		}

	case []Inventory:
		for i, inventory := range objects {
			node := Node{}
			inventory.UUID = uuid.NewString()
			node.Id = inventory.UUID
			objects[i] = inventory
			node.Kinds = []string{"ATInventory"}
			node.Properties = map[string]string{
				"id":                              strconv.FormatInt(int64(inventory.ID), 10),
				"name":                            inventory.Name,
				"description":                     inventory.Description,
				"url":                             inventory.Url,
				"organization":                    strconv.FormatInt(int64(inventory.Organization), 10),
				"kind":                            inventory.Kind,
				"host_filter":                     inventory.HostFilter,
				"has_active_failures":             strconv.FormatBool(inventory.HasActiveFailures),
				"has_inventory_source":            strconv.FormatBool(inventory.HasInventorySources),
				"total_hosts":                     strconv.FormatInt(int64(inventory.TotalHosts), 10),
				"hosts_with_active_failures":      strconv.FormatInt(int64(inventory.HostsWithActiveFailures), 10),
				"total_groups":                    strconv.FormatInt(int64(inventory.TotalGroups), 10),
				"total_inventory_sources":         strconv.FormatInt(int64(inventory.TotalInventorySources), 10),
				"inventory_sources_with_failures": strconv.FormatInt(int64(inventory.InventorySourcesWithFailures), 10),
				"pending_deletion":                strconv.FormatBool(inventory.PendingDeletion),
				"prevent_instance_group_fallback": strconv.FormatBool(inventory.PreventInstanceGroupFallback),
			}
			nodes = append(nodes, node)
		}

	case []Host:
		for i, host := range objects {
			node := Node{}
			node.Kinds = []string{
				"ATHost",
			}
			host.UUID = uuid.NewString()
			node.Id = host.UUID
			objects[i] = host
			node.Properties = map[string]string{
				"id":                     strconv.Itoa(host.ID),
				"name":                   host.Name,
				"description":            host.Description,
				"url":                    host.Url,
				"type":                   host.Type,
				"created":                host.Created,
				"modified":               host.Modified,
				"inventory":              strconv.FormatInt(int64(host.Inventory), 10),
				"enabled":                strconv.FormatBool(host.Enabled),
				"instance_id":            host.InstanceId,
				"variables":              host.Variables,
				"has_active_failures":    strconv.FormatBool(host.HasActiveFailures),
				"last_job":               strconv.FormatInt(int64(host.LastJob), 10),
				"last_job_host_summary":  strconv.FormatInt(int64(host.LastJobHostSummary), 10),
				"ansible_facts_modified": host.AnsibleFactsModified,
			}
			nodes = append(nodes, node)
		}

	case []Team:
		for i, team := range objects {
			node := Node{}
			node.Kinds = []string{
				"ATTeam",
			}
			team.UUID = uuid.NewString()
			node.Id = team.UUID
			objects[i] = team
			node.Properties = map[string]string{
				"id":          strconv.Itoa(team.ID),
				"name":        team.Name,
				"description": team.Description,
				"url":         team.Url,
				"type":        team.Type,
				"created":     team.Created,
				"modified":    team.Modified,
			}
			nodes = append(nodes, node)
		}
	}

	return nodes

}

func CalculateName(objectType string) string {
	now := time.Now()
	epoch := now.Unix()
	return fmt.Sprintf("%d_%s.json", epoch, objectType)
}

func WriteToFile(content []byte, filePath string) error {

	log.Debug(fmt.Sprintf("Writing to file `%s`.", filePath))

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(content)
	if err != nil {
		return err
	}

	return nil
}
