package ansible

import (
	"encoding/json"
	"strconv"

	"github.com/TheManticoreProject/gopengraph/node"
	"github.com/TheManticoreProject/gopengraph/properties"
)

type Job struct {
	Object
	Inventory             int            `json:"inventory"`
	Project               int            `json:"project"`
	Organization          int            `json:"organization,omitempty"`
	Playbook              string         `json:"playbook"`
	ScmBranch             string         `json:"scm_branch,omitempty"`
	Forks                 int            `json:"forks,omitempty"`
	Limit                 string         `json:"limit,omitempty"`
	Verbosity             int            `json:"verbosity,omitempty"`
	ExtraVars             string         `json:"extra_vars,omitempty"`
	Started               string         `json:"started,omitempty"`
	Finished              string         `json:"finished,omitempty"`
	CanceledOn            string         `json:"canceled_on,omitempty"`
	Elapsed               float32        `json:"elapsed,omitempty"`
	JobExplanation        string         `json:"job_explanation,omitempty"`
	Created               string         `json:"created,omitempty"`
	Modified              string         `json:"modified,omitempty"`
	UnifiedJobTemplate    int            `json:"unified_job_template"`
	LaunchType            string         `json:"launch_type"`
	Failed                bool           `json:"failed"`
	Status                string         `json:"status,omitempty"`
	ExecutionEnvironment  int            `json:"execution_environment,omitempty"`
	ExecutionNode         string         `json:"execution_node,omitempty"`
	ControllerNode        string         `json:"controller_node,omitempty"`
	LaunchedBy            map[string]any `json:"launched_by,omitempty"`
	WorkUnitId            string         `json:"work_unit_id,omitempty"`
	JobTags               string         `json:"job_tags,omitempty"`
	JobType               string         `json:"job_type,omitempty"`
	ForceHandler          bool           `json:"force_handlers,omitempty"`
	SkipTags              string         `json:"skip_tags,omitempty"`
	StartAtTask           string         `json:"start_at_task,omitempty"`
	Timeout               int            `json:"timeout,omitempty"`
	UseFactCache          bool           `json:"use_fact_cache,omitempty"`
	PasswordNeededToStart string         `json:"password_needed_to_start,omitempty"`
	AllowSimultaneous     bool           `json:"allow_simultaneous,omitempty"`
	Artifacts             map[string]any `json:"artifacts,omitempty"`
	ScmRevision           string         `json:"scm_revision,omitempty"`
	InstanceGroup         int            `json:"instance_group,omitempty"`
	DiffMode              bool           `json:"diff_mode,omitempty"`
	JobSliceNumber        int            `json:"job_slice_number,omitempty"`
	JobSliceCount         int            `json:"job_slice_count,omitempty"`
	WebhookGuid           string         `json:"webhook_guid,omitempty"`
	WebhookService        string         `json:"webhook_service,omitempty"`
	WebhookCredential     int            `json:"webhook_credential,omitempty"`
}

func (j Job) MarshalJSON() ([]byte, error) {
	type job Job
	return json.MarshalIndent((job)(j), "", "  ")
}

func (j *Job) ToBHNode() (n *node.Node) {
	props := properties.NewProperties()
	props.SetProperty("id", strconv.FormatInt(int64(j.ID), 10))
	props.SetProperty("name", j.Name)
	props.SetProperty("description", j.Description)
	props.SetProperty("url", j.Url)
	props.SetProperty("type", j.Type)
	props.SetProperty("created", j.Created)
	props.SetProperty("modified", j.Modified)
	props.SetProperty("playbook", j.Playbook)
	props.SetProperty("scm_branch", j.ScmBranch)
	props.SetProperty("forks", strconv.FormatInt(int64(j.Forks), 10))
	props.SetProperty("limit", j.Limit)
	props.SetProperty("verbosity", strconv.FormatInt(int64(j.Verbosity), 10))
	props.SetProperty("extra_vars", j.ExtraVars)
	props.SetProperty("started", j.Started)
	props.SetProperty("finished", j.Finished)
	props.SetProperty("canceled_on", j.CanceledOn)
	props.SetProperty("elapsed", strconv.FormatFloat(float64(j.Elapsed), 'e', 10, 32))
	props.SetProperty("job_explanation", j.JobExplanation)
	props.SetProperty("launch_type", j.LaunchType)
	props.SetProperty("unified_job_template", strconv.FormatInt(int64(j.UnifiedJobTemplate), 10))
	props.SetProperty("organization", strconv.FormatInt(int64(j.Organization), 10))
	props.SetProperty("inventory", strconv.FormatInt(int64(j.Inventory), 10))
	props.SetProperty("project", strconv.FormatInt(int64(j.Project), 10))
	props.SetProperty("failed", strconv.FormatBool(j.Failed))
	props.SetProperty("status", j.Status)
	props.SetProperty("execution_environment", strconv.FormatInt(int64(j.ExecutionEnvironment), 10))
	props.SetProperty("execution_node", j.ExecutionNode)
	props.SetProperty("controller_node", j.ControllerNode)
	props.SetProperty("work_unit_id", j.WorkUnitId)
	props.SetProperty("job_tags", j.JobTags)
	props.SetProperty("job_type", j.JobType)
	props.SetProperty("force_handler", strconv.FormatBool(j.ForceHandler))
	props.SetProperty("skip_tags", j.SkipTags)
	props.SetProperty("start_at_task", j.StartAtTask)
	props.SetProperty("timeout", strconv.FormatInt(int64(j.Timeout), 10))
	props.SetProperty("use_fact_cache", strconv.FormatBool(j.UseFactCache))
	props.SetProperty("password_needed_to_start", j.PasswordNeededToStart)
	props.SetProperty("allow_simultaneous", strconv.FormatBool(j.AllowSimultaneous))
	props.SetProperty("scm_revision", j.ScmRevision)
	props.SetProperty("instance_group", strconv.FormatInt(int64(j.InstanceGroup), 10))
	props.SetProperty("diff_mode", strconv.FormatBool(j.DiffMode))
	props.SetProperty("job_slice_number", strconv.FormatInt(int64(j.JobSliceNumber), 10))
	props.SetProperty("job_slice_count", strconv.FormatInt(int64(j.JobSliceCount), 10))
	props.SetProperty("webhook_guid", j.WebhookGuid)
	props.SetProperty("webhook_service", j.WebhookService)
	props.SetProperty("webhook_credential", strconv.FormatInt(int64(j.WebhookCredential), 10))
	n, _ = node.NewNode(j.OID, []string{"ATJob"}, props)

	return n
}
