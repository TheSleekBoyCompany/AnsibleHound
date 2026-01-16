package ansible

import (
	"ansible-hound/core/opengraph"
	"encoding/json"
	"strconv"
)

type JobTemplate struct {
	Object
	JobType                         string              `json:"job_type"`
	Inventory                       int                 `json:"inventory"`
	Project                         int                 `json:"project"`
	Organization                    int                 `json:"organization,omitempty"`
	Playbook                        string              `json:"playbook"`
	SCMBranch                       string              `json:"scm_branch,omitempty"`
	Limit                           string              `json:"limit"`
	Verbosity                       int                 `json:"verbosity"`
	Credentials                     map[int]*Credential `json:"credentials"`
	ExtraVars                       string              `json:"extra_vars"`
	Status                          string              `json:"status,omitempty"`
	JobTags                         string              `json:"job_tags,omitempty"`
	Forks                           int                 `json:"forks"`
	SkipTags                        string              `json:"skip_tags,omitempty"`
	StartAtTask                     string              `json:"start_at_task,omitempty"`
	Timeout                         int                 `json:"timeout,omitempty"`
	UseFactCache                    bool                `json:"use_fact_cache,omitempty"`
	ForceHandler                    bool                `json:"force_handlers,omitempty"`
	LastJobRun                      string              `json:"last_job_run,omitempty"`
	NextJobRun                      string              `json:"next_job_run,omitempty"`
	LastJobFailed                   bool                `json:"last_job_failed,omitempty"`
	ExecutionEnvironment            int                 `json:"execution_environment,omitempty"`
	HostConfigKey                   string              `json:"host_config_key,omitempty"`
	AskScmBranchOnLaunch            bool                `json:"ask_scm_branch_on_launch,omitempty"`
	AskDiffModeOnLaunch             bool                `json:"ask_diff_mode_on_launch,omitempty"`
	AskVariablesOnLaunch            bool                `json:"ask_variables_on_launch,omitempty"`
	AskLimitOnLaunch                bool                `json:"ask_limit_on_launch,omitempty"`
	AskTagsOnLaunch                 bool                `json:"ask_tags_on_launch,omitempty"`
	AskJobTypeOnLaunch              bool                `json:"ask_job_type_on_launch,omitempty"`
	AskVerbosityOnLaunch            bool                `json:"ask_verbosity_on_launch,omitempty"`
	AskInventoryOnLaunch            bool                `json:"ask_inventory_on_launch,omitempty"`
	AskCredentialOnLaunch           bool                `json:"ask_credential_on_launch,omitempty"`
	AskExecutionEnvironmentOnLaunch bool                `json:"ask_execution_environment_on_launch,omitempty"`
	AskLabelsOnLaunch               bool                `json:"ask_labels_on_launch,omitempty"`
	AskForksOnLaunch                bool                `json:"ask_forks_on_launch,omitempty"`
	AskJobSliceCountOnLaunch        bool                `json:"ask_job_slice_count_on_launch,omitempty"`
	AskTimeoutOnLaunch              bool                `json:"ask_timeout_on_launch,omitempty"`
	AskInstanceGroupsOnLaunch       bool                `json:"ask_instance_groups_on_launch,omitempty"`
	SurveyEnabled                   bool                `json:"survey_enabled,omitempty"`
	BecomeEnabled                   bool                `json:"become_enabled,omitempty"`
	DiffMode                        bool                `json:"diff_mode,omitempty"`
	AllowSimultaneous               bool                `json:"allow_simultaneous,omitempty"`
	CustomVirtualenv                string              `json:"custom_virtualenv,omitempty"`
	JobSliceCount                   int                 `json:"job_slice_count,omitempty"`
	WebhookService                  string              `json:"webhook_service,omitempty"`
	WebhookCredential               int                 `json:"webhook_credential,omitempty"`
	PreventInstanceGroupFallback    bool                `json:"prevent_instance_group_fallback,omitempty"`
}

func (j JobTemplate) MarshalJSON() ([]byte, error) {
	type jobTemplate JobTemplate
	return json.MarshalIndent((jobTemplate)(j), "", "  ")
}

func (j *JobTemplate) ToBHNode() (node opengraph.Node) {
	node.Kinds = []string{
		"ATJobTemplate",
	}
	node.Id = j.OID
	node.Properties = map[string]string{
		"id":                                  strconv.Itoa(j.ID),
		"name":                                j.Name,
		"description":                         j.Description,
		"url":                                 j.Url,
		"job_type":                            j.JobType,
		"project":                             strconv.FormatInt(int64(j.Project), 10),
		"organization":                        strconv.FormatInt(int64(j.Organization), 10),
		"playbook":                            j.Playbook,
		"scm_branch":                          j.SCMBranch,
		"limit":                               j.Limit,
		"verbosity":                           strconv.FormatInt(int64(j.Verbosity), 10),
		"extra_vars":                          j.ExtraVars,
		"status":                              j.Status,
		"job_tags":                            j.JobTags,
		"forks":                               strconv.FormatInt(int64(j.Forks), 10),
		"skip_tags":                           j.SkipTags,
		"start_at_task":                       j.StartAtTask,
		"timeout":                             strconv.FormatInt(int64(j.Timeout), 10),
		"use_fact_cache":                      strconv.FormatBool(j.UseFactCache),
		"force_handler":                       strconv.FormatBool(j.ForceHandler),
		"last_job_run":                        j.LastJobRun,
		"next_job_run":                        j.NextJobRun,
		"last_job_failed":                     strconv.FormatBool(j.LastJobFailed),
		"execution_environment":               strconv.FormatInt(int64(j.ExecutionEnvironment), 10),
		"host_config_key":                     j.HostConfigKey,
		"ask_scm_branch_on_launch":            strconv.FormatBool(j.AskScmBranchOnLaunch),
		"ask_diff_mode_on_launch":             strconv.FormatBool(j.AskDiffModeOnLaunch),
		"ask_variables_on_launch":             strconv.FormatBool(j.AskVariablesOnLaunch),
		"ask_limit_on_launch":                 strconv.FormatBool(j.AskLimitOnLaunch),
		"ask_tags_on_launch":                  strconv.FormatBool(j.AskTagsOnLaunch),
		"ask_job_type_on_launch":              strconv.FormatBool(j.AskJobTypeOnLaunch),
		"ask_verbosity_on_launch":             strconv.FormatBool(j.AskVerbosityOnLaunch),
		"ask_inventory_on_launch":             strconv.FormatBool(j.AskInventoryOnLaunch),
		"ask_credential_on_launch":            strconv.FormatBool(j.AskCredentialOnLaunch),
		"ask_execution_environment_on_launch": strconv.FormatBool(j.AskExecutionEnvironmentOnLaunch),
		"ask_labels_on_launch":                strconv.FormatBool(j.AskLabelsOnLaunch),
		"ask_forks_on_launch":                 strconv.FormatBool(j.AskForksOnLaunch),
		"ask_job_slice_count_on_launch":       strconv.FormatBool(j.AskJobSliceCountOnLaunch),
		"ask_timeout_on_launch":               strconv.FormatBool(j.AskTimeoutOnLaunch),
		"ask_instance_groups_on_launch":       strconv.FormatBool(j.AskInstanceGroupsOnLaunch),
		"survey_enabled":                      strconv.FormatBool(j.SurveyEnabled),
		"become_enabled":                      strconv.FormatBool(j.BecomeEnabled),
		"diff_mode":                           strconv.FormatBool(j.DiffMode),
		"allow_simultaneous":                  strconv.FormatBool(j.AllowSimultaneous),
		"custom_virtualenv":                   j.CustomVirtualenv,
		"job_slice_count":                     strconv.FormatInt(int64(j.JobSliceCount), 10),
		"webhook_service":                     j.WebhookService,
		"webhook_credential":                  strconv.FormatInt(int64(j.WebhookCredential), 10),
		"prevent_instance_group_fallback":     strconv.FormatBool(j.PreventInstanceGroupFallback),
	}
	return node
}
