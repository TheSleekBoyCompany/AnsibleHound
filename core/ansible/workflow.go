package ansible

import (
	"encoding/json"
	"strconv"

	"github.com/TheManticoreProject/gopengraph/node"
	"github.com/TheManticoreProject/gopengraph/properties"
)

type WorkflowJobTemplate struct {
	Object
	LastJobRun           string `json:"last_job_run,omitempty"`
	LastJobFailed        bool   `json:"last_job_failed,omitempty"`
	NextJobRun           string `json:"next_job_run,omitempty"`
	Status               string `json:"status,omitempty"`
	ExtraVars            string `json:"extra_vars"`
	Organization         int    `json:"organization,omitempty"`
	SurveyEnabled        bool   `json:"survey_enabled,omitempty"`
	AllowSimultaneous    bool   `json:"allow_simultaneous,omitempty"`
	AskVariablesOnLaunch bool   `json:"ask_variables_on_launch,omitempty"`
	Inventory            int    `json:"inventory"`
	Limit                string `json:"limit"`
	SCMBranch            string `json:"scm_branch,omitempty"`
	AskInventoryOnLaunch bool   `json:"ask_inventory_on_launch,omitempty"`
	AskScmBranchOnLaunch bool   `json:"ask_scm_branch_on_launch,omitempty"`
	AskLimitOnLaunch     bool   `json:"ask_limit_on_launch,omitempty"`
	WebhookService       string `json:"webhook_service,omitempty"`
	WebhookCredential    int    `json:"webhook_credential,omitempty"`
	AskLabelsOnLaunch    bool   `json:"ask_labels_on_launch,omitempty"`
	AskSkipTagsOnLaunch  bool   `json:"ask_skip_tags_on_launch,omitempty"`
	AskTagsOnLaunch      bool   `json:"ask_tags_on_launch,omitempty"`
	SkipTags             string `json:"skip_tags,omitempty"`
	JobTags              string `json:"job_tags,omitempty"`
}

func (w WorkflowJobTemplate) MarshalJSON() ([]byte, error) {
	type wkf WorkflowJobTemplate
	return json.MarshalIndent((wkf)(w), "", "  ")
}

func (w *WorkflowJobTemplate) ToBHNode() (n *node.Node) {
	props := properties.NewProperties()
	props.SetProperty("id", strconv.Itoa(w.ID))
	props.SetProperty("name", w.Name)
	props.SetProperty("description", w.Description)
	props.SetProperty("url", w.Url)
	props.SetProperty("organization", strconv.FormatInt(int64(w.Organization), 10))
	props.SetProperty("extra_vars", w.ExtraVars)
	props.SetProperty("status", w.Status)
	props.SetProperty("last_job_run", w.LastJobRun)
	props.SetProperty("next_job_run", w.NextJobRun)
	props.SetProperty("last_job_failed", strconv.FormatBool(w.LastJobFailed))
	props.SetProperty("survey_enabled", strconv.FormatBool(w.SurveyEnabled))
	props.SetProperty("allow_simultaneous", strconv.FormatBool(w.AllowSimultaneous))
	props.SetProperty("ask_variables_on_launch", strconv.FormatBool(w.AskVariablesOnLaunch))
	props.SetProperty("inventory", strconv.FormatInt(int64(w.Inventory), 10))
	props.SetProperty("limit", w.Limit)
	props.SetProperty("scm_branch", w.SCMBranch)
	props.SetProperty("ask_inventory_on_launch", strconv.FormatBool(w.AskInventoryOnLaunch))
	props.SetProperty("ask_scm_branch_on_launch", strconv.FormatBool(w.AskScmBranchOnLaunch))
	props.SetProperty("ask_limit_on_launch", strconv.FormatBool(w.AskLimitOnLaunch))
	props.SetProperty("webhook_service", w.WebhookService)
	props.SetProperty("webhook_credential", strconv.FormatInt(int64(w.WebhookCredential), 10))
	props.SetProperty("ask_labels_on_launch", strconv.FormatBool(w.AskLabelsOnLaunch))
	props.SetProperty("ask_skip_tags_on_launch", strconv.FormatBool(w.AskSkipTagsOnLaunch))
	props.SetProperty("ask_tags_on_launch", strconv.FormatBool(w.AskTagsOnLaunch))
	props.SetProperty("job_tags", w.JobTags)
	props.SetProperty("skip_tags", w.SkipTags)
	n, _ = node.NewNode(w.OID, []string{"ATWorkflowJobTemplate"}, props)

	return n
}

type WorkflowJobTemplateNode struct {
	Object
	Inventory              int    `json:"inventory,omitempty"`
	SCMBranch              string `json:"scm_branch,omitempty"`
	JobType                string `json:"job_type,omitempty"`
	SkipTags               string `json:"skip_tags,omitempty"`
	JobTags                string `json:"job_tags,omitempty"`
	Limit                  string `json:"limit"`
	DiffMode               bool   `json:"diff_mode,omitempty"`
	Verbosity              int    `json:"verbosity"`
	ExecutionEnvironment   int    `json:"execution_environment,omitempty"`
	Forks                  int    `json:"forks"`
	JobSliceCount          int    `json:"job_slice_count,omitempty"`
	Timeout                int    `json:"timeout,omitempty"`
	WorkflowJobTemplate    int    `json:"workflow_job_template"`
	UnifiedJobTemplate     int    `json:"unified_job_template"`
	SuccessNodes           []int  `json:"success_nodes"`
	FailureNodes           []int  `json:"failure_nodes"`
	AlwaysNodes            []int  `json:"always_nodes"`
	AllParentsMustConverge bool   `json:"all_parents_must_converge"`
}

func (w WorkflowJobTemplateNode) MarshalJSON() ([]byte, error) {
	type wkf WorkflowJobTemplateNode
	return json.MarshalIndent((wkf)(w), "", "  ")
}

func (w *WorkflowJobTemplateNode) ToBHNode() (n *node.Node) {
	props := properties.NewProperties()
	props.SetProperty("id", strconv.Itoa(w.ID))
	props.SetProperty("name", w.Name)
	props.SetProperty("description", w.Description)
	props.SetProperty("url", w.Url)
	props.SetProperty("inventory", w.Inventory)
	props.SetProperty("scm_branch", w.SCMBranch)
	props.SetProperty("job_type", w.JobType)
	props.SetProperty("skip_tags", w.SkipTags)
	props.SetProperty("job_tags", w.JobTags)
	props.SetProperty("limit", w.Limit)
	props.SetProperty("diff_mode", strconv.FormatBool(w.DiffMode))
	props.SetProperty("verbosity", strconv.FormatInt(int64(w.Verbosity), 10))
	props.SetProperty("execution_environment", strconv.FormatInt(int64(w.ExecutionEnvironment), 10))
	props.SetProperty("forks", strconv.FormatInt(int64(w.Forks), 10))
	props.SetProperty("job_slice_count", strconv.FormatInt(int64(w.JobSliceCount), 10))
	props.SetProperty("timeout", strconv.FormatInt(int64(w.Timeout), 10))
	props.SetProperty("workflow_job_template", strconv.FormatInt(int64(w.WorkflowJobTemplate), 10))
	props.SetProperty("unified_job_template", strconv.FormatInt(int64(w.UnifiedJobTemplate), 10))
	props.SetProperty("all_parents_must_converge", strconv.FormatBool(w.AllParentsMustConverge))
	n, _ = node.NewNode(w.OID, []string{"ATWorkflowJobTemplateNode"}, props)

	return n
}
