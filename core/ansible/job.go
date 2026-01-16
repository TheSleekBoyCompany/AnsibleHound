package ansible

import (
	"ansible-hound/core/opengraph"
	"encoding/json"
	"strconv"
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

func (j *Job) ToBHNode() (node opengraph.Node) {
	node.Kinds = []string{
		"ATJob",
	}
	node.Id = j.OID
	node.Properties = map[string]string{
		"id":                       strconv.FormatInt(int64(j.ID), 10),
		"name":                     j.Name,
		"description":              j.Description,
		"url":                      j.Url,
		"type":                     j.Type,
		"created":                  j.Created,
		"modified":                 j.Modified,
		"playbook":                 j.Playbook,
		"scm_branch":               j.ScmBranch,
		"forks":                    strconv.FormatInt(int64(j.Forks), 10),
		"limit":                    j.Limit,
		"verbosity":                strconv.FormatInt(int64(j.Verbosity), 10),
		"extra_vars":               j.ExtraVars,
		"started":                  j.Started,
		"finished":                 j.Finished,
		"canceled_on":              j.CanceledOn,
		"elapsed":                  strconv.FormatFloat(float64(j.Elapsed), 'e', 10, 32),
		"job_explanation":          j.JobExplanation,
		"launch_type":              j.LaunchType,
		"unified_job_template":     strconv.FormatInt(int64(j.UnifiedJobTemplate), 10),
		"organization":             strconv.FormatInt(int64(j.Organization), 10),
		"inventory":                strconv.FormatInt(int64(j.Inventory), 10),
		"project":                  strconv.FormatInt(int64(j.Project), 10),
		"failed":                   strconv.FormatBool(j.Failed),
		"status":                   j.Status,
		"execution_environment":    strconv.FormatInt(int64(j.ExecutionEnvironment), 10),
		"execution_node":           j.ExecutionNode,
		"controller_node":          j.ControllerNode,
		"work_unit_id":             j.WorkUnitId,
		"job_tags":                 j.JobTags,
		"job_type":                 j.JobType,
		"force_handler":            strconv.FormatBool(j.ForceHandler),
		"skip_tags":                j.SkipTags,
		"start_at_task":            j.StartAtTask,
		"timeout":                  strconv.FormatInt(int64(j.Timeout), 10),
		"use_fact_cache":           strconv.FormatBool(j.UseFactCache),
		"password_needed_to_start": j.PasswordNeededToStart,
		"allow_simultaneous":       strconv.FormatBool(j.AllowSimultaneous),
		"scm_revision":             j.ScmRevision,
		"instance_group":           strconv.FormatInt(int64(j.InstanceGroup), 10),
		"diff_mode":                strconv.FormatBool(j.DiffMode),
		"job_slice_number":         strconv.FormatInt(int64(j.JobSliceNumber), 10),
		"job_slice_count":          strconv.FormatInt(int64(j.JobSliceCount), 10),
		"webhook_guid":             j.WebhookGuid,
		"webhook_service":          j.WebhookService,
		"webhook_credential":       strconv.FormatInt(int64(j.WebhookCredential), 10),
	}
	return node
}
