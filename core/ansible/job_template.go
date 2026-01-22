package ansible

import (
	"encoding/json"
	"strconv"

	"github.com/Ramoreik/gopengraph/node"
	"github.com/Ramoreik/gopengraph/properties"
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

func (j *JobTemplate) ToBHNode() (n *node.Node) {
	props := properties.NewProperties()
	props.SetProperty("id", strconv.Itoa(j.ID))
	props.SetProperty("name", j.Name)
	props.SetProperty("description", j.Description)
	props.SetProperty("url", j.Url)
	props.SetProperty("job_type", j.JobType)
	props.SetProperty("project", strconv.FormatInt(int64(j.Project), 10))
	props.SetProperty("organization", strconv.FormatInt(int64(j.Organization), 10))
	props.SetProperty("playbook", j.Playbook)
	props.SetProperty("scm_branch", j.SCMBranch)
	props.SetProperty("limit", j.Limit)
	props.SetProperty("verbosity", strconv.FormatInt(int64(j.Verbosity), 10))
	props.SetProperty("extra_vars", j.ExtraVars)
	props.SetProperty("status", j.Status)
	props.SetProperty("job_tags", j.JobTags)
	props.SetProperty("forks", strconv.FormatInt(int64(j.Forks), 10))
	props.SetProperty("skip_tags", j.SkipTags)
	props.SetProperty("start_at_task", j.StartAtTask)
	props.SetProperty("timeout", strconv.FormatInt(int64(j.Timeout), 10))
	props.SetProperty("use_fact_cache", strconv.FormatBool(j.UseFactCache))
	props.SetProperty("force_handler", strconv.FormatBool(j.ForceHandler))
	props.SetProperty("last_job_run", j.LastJobRun)
	props.SetProperty("next_job_run", j.NextJobRun)
	props.SetProperty("last_job_failed", strconv.FormatBool(j.LastJobFailed))
	props.SetProperty("execution_environment", strconv.FormatInt(int64(j.ExecutionEnvironment), 10))
	props.SetProperty("host_config_key", j.HostConfigKey)
	props.SetProperty("ask_scm_branch_on_launch", strconv.FormatBool(j.AskScmBranchOnLaunch))
	props.SetProperty("ask_diff_mode_on_launch", strconv.FormatBool(j.AskDiffModeOnLaunch))
	props.SetProperty("ask_variables_on_launch", strconv.FormatBool(j.AskVariablesOnLaunch))
	props.SetProperty("ask_limit_on_launch", strconv.FormatBool(j.AskLimitOnLaunch))
	props.SetProperty("ask_tags_on_launch", strconv.FormatBool(j.AskTagsOnLaunch))
	props.SetProperty("ask_job_type_on_launch", strconv.FormatBool(j.AskJobTypeOnLaunch))
	props.SetProperty("ask_verbosity_on_launch", strconv.FormatBool(j.AskVerbosityOnLaunch))
	props.SetProperty("ask_inventory_on_launch", strconv.FormatBool(j.AskInventoryOnLaunch))
	props.SetProperty("ask_credential_on_launch", strconv.FormatBool(j.AskCredentialOnLaunch))
	props.SetProperty("ask_execution_environment_on_launch", strconv.FormatBool(j.AskExecutionEnvironmentOnLaunch))
	props.SetProperty("ask_labels_on_launch", strconv.FormatBool(j.AskLabelsOnLaunch))
	props.SetProperty("ask_forks_on_launch", strconv.FormatBool(j.AskForksOnLaunch))
	props.SetProperty("ask_job_slice_count_on_launch", strconv.FormatBool(j.AskJobSliceCountOnLaunch))
	props.SetProperty("ask_timeout_on_launch", strconv.FormatBool(j.AskTimeoutOnLaunch))
	props.SetProperty("ask_instance_groups_on_launch", strconv.FormatBool(j.AskInstanceGroupsOnLaunch))
	props.SetProperty("survey_enabled", strconv.FormatBool(j.SurveyEnabled))
	props.SetProperty("become_enabled", strconv.FormatBool(j.BecomeEnabled))
	props.SetProperty("diff_mode", strconv.FormatBool(j.DiffMode))
	props.SetProperty("allow_simultaneous", strconv.FormatBool(j.AllowSimultaneous))
	props.SetProperty("custom_virtualenv", j.CustomVirtualenv)
	props.SetProperty("job_slice_count", strconv.FormatInt(int64(j.JobSliceCount), 10))
	props.SetProperty("webhook_service", j.WebhookService)
	props.SetProperty("webhook_credential", strconv.FormatInt(int64(j.WebhookCredential), 10))
	props.SetProperty("prevent_instance_group_fallback", strconv.FormatBool(j.PreventInstanceGroupFallback))
	n, _ = node.NewNode(j.OID, []string{"ATJobTemplate"}, props)

	return n
}
