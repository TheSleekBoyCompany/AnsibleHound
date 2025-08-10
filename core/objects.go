package core

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

var Instance AnsibleInstance

type AnsibleTypeList interface {
	[]*User | []*Job | []*JobTemplate | []*Inventory |
		[]*Project | []*Organization | []*Role | []*Credential |
		[]*Host | []*Team | []*RoleDefinition | []*RoleTeamAssignments |
		[]*RoleUserAssignments
}

type AnsibleType interface {
	GetID() int
	GetOID() string
	InitOID(string)
	ToBHNode() Node
}

type Object struct {
	OID         string `json:"uuid,omitempty"`
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Url         string `json:"url,omitempty"`
	Type        string `json:"type,omitempty"`
	Created     string `json:"created,omitempty"`
	Modified    string `json:"modified,omitempty"`
}

func (o Object) GetOID() (uuid string) {
	return o.OID
}

func (o Object) GetID() (id int) {
	return o.ID
}

func (o *Object) InitOID(installUUID string) {
	data := fmt.Sprintf("%s_%s_%s", installUUID, strconv.Itoa(o.ID), o.Type)
	hasher := sha1.New()
	hasher.Write([]byte(data))
	hashBytes := hasher.Sum(nil)
	o.OID = hex.EncodeToString(hashBytes)
}

type Response[T any] struct {
	Count   int `json:"count"`
	Results []T `json:"results"`
}

type AnsibleInstance struct {
	Object
	Version     string `json:"version"`
	ActiveNode  string `json:"active_node"`
	InstallUUID string `json:"install_uuid"`
}

func (i *AnsibleInstance) MarshalJSON() ([]byte, error) {
	type instance AnsibleInstance
	return json.MarshalIndent((*instance)(i), "", "  ")
}

func (i *AnsibleInstance) ToBHNode() (node Node) {
	node.Kinds = []string{
		"ATAnsibleInstance",
	}
	node.Id = i.OID
	node.Properties = map[string]string{
		"name":         i.Name,
		"version":      i.Version,
		"active_node":  i.ActiveNode,
		"install_uuid": i.InstallUUID,
	}
	return node
}

type User struct {
	Object
	Username        string        `json:"username"`
	FirstName       string        `json:"first_name,omitempty"`
	LastName        string        `json:"last_name,omitempty"`
	Email           string        `json:"email,omitempty"`
	IsSuperUser     bool          `json:"is_superuser,omitempty"`
	IsSystemAuditor bool          `json:"is_sytem_auditor,omitempty"`
	LdapDn          string        `json:"ldap_dn,omitempty"`
	LastLogin       string        `json:"last_login,omitempty"`
	ExternalAccount string        `json:"external_account,omitempty"`
	Roles           map[int]*Role `json:"roles"`
}

func (u *User) MarshalJSON() ([]byte, error) {
	type user User
	return json.MarshalIndent((*user)(u), "", "  ")
}

func (u *User) ToBHNode() (node Node) {
	node.Kinds = []string{
		"ATUser",
	}
	node.Id = u.OID
	node.Properties = map[string]string{
		"id":                strconv.Itoa(u.ID),
		"name":              u.Username,
		"description":       u.Description,
		"url":               u.Url,
		"firstname":         u.FirstName,
		"lastname":          u.LastName,
		"email":             u.Email,
		"is_super_user":     strconv.FormatBool(u.IsSuperUser),
		"is_system_auditor": strconv.FormatBool(u.IsSystemAuditor),
		"ldap_dn":           u.LdapDn,
		"last_login":        u.LastLogin,
		"external_account":  u.ExternalAccount,
	}
	return node
}

type Team struct {
	Object
	Organization int           `json:"organization,omitempty"`
	Members      map[int]*User `json:"members"`
	Roles        map[int]*Role `json:"roles"`
}

func (u *Team) MarshalJSON() ([]byte, error) {
	type team Team
	return json.MarshalIndent((*team)(u), "", "  ")
}

func (t *Team) ToBHNode() (node Node) {
	node.Kinds = []string{
		"ATTeam",
	}
	node.Id = t.OID
	node.Properties = map[string]string{
		"id":          strconv.Itoa(t.ID),
		"name":        t.Name,
		"description": t.Description,
		"url":         t.Url,
		"type":        t.Type,
		"created":     t.Created,
		"modified":    t.Modified,
	}
	return node
}

type Role struct {
	Object
	UserOnly      bool              `json:"user_only,omitempty"`
	SummaryFields RoleSummaryFields `json:"summary_fields"`
}

type RoleSummaryFields struct {
	ResourceName            string `json:"resource_name"`
	ResourceType            string `json:"resource_type"`
	ResourceTypeDisplayName string `json:"resource_type_display_name"`
	ResourceId              int    `json:"resource_id"`
}

func (r Role) MarshalJSON() ([]byte, error) {
	type role Role
	return json.MarshalIndent((role)(r), "", "  ")
}

func (r *Role) ToBHNode() (node Node) {
	return Node{}
}

type Organization struct {
	Object
	MaxHosts           int    `json:"max_hosts,omitempty"`
	CustomVirtualenv   string `json:"custom_virtualenv,omitempty"`
	DefaultEnvironment int    `json:"default_environment,omitempty"`
}

func (o Organization) MarshalJSON() ([]byte, error) {
	type organization Organization
	return json.MarshalIndent((organization)(o), "", "  ")
}

func (o *Organization) ToBHNode() (node Node) {
	node.Kinds = []string{
		"ATOrganization",
	}
	node.Id = o.OID
	node.Properties = map[string]string{
		"id":                  strconv.Itoa(o.ID),
		"name":                o.Name,
		"description":         o.Description,
		"url":                 o.Url,
		"max_hosts":           strconv.FormatInt(int64(o.MaxHosts), 10),
		"custom_virtualenv":   o.CustomVirtualenv,
		"default_environment": strconv.FormatInt(int64(o.DefaultEnvironment), 10),
	}
	return node
}

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

func (j *JobTemplate) ToBHNode() (node Node) {
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

func (j *Job) ToBHNode() (node Node) {
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

type Project struct {
	Object
	Organization                  int    `json:"organization"`
	Status                        string `json:"status,omitempty"`
	LocalPath                     string `json:"local_path,omitempty"`
	ScmType                       string `json:"scm_type,omitempty"`
	ScmUrl                        string `json:"scm_url,omitempty"`
	ScmBranch                     string `json:"scm_branch,omitempty"`
	ScmRefSpec                    string `json:"scm_refspec,omitempty"`
	ScmClean                      bool   `json:"scm_clean,omitempty"`
	ScmTrackSubmodules            bool   `json:"scm_track_submodules,omitempty"`
	ScmDeleteOnUpdate             bool   `json:"scm_delete_on_update,omitempty"`
	Credential                    int    `json:"credential,omitempty"`
	Timeout                       int    `json:"timeout,omitempty"`
	ScmRevision                   string `json:"scm_revision,omitempty"`
	LastJobRun                    string `json:"last_job_run,omitempty"`
	NextJobRun                    string `json:"next_job_run,omitempty"`
	LastJobFailed                 bool   `json:"last_job_failed,omitempty"`
	ScmUpdateOnLaunch             bool   `json:"scm_update_on_launch,omitempty"`
	ScmUpdateCacheTimeout         int    `json:"scm_update_cache_timeout"`
	AllowOverride                 bool   `json:"allow_override,omitempty"`
	CustomVirtualenv              string `json:"custom_virtualenv,omitempty"`
	DefaultEnvironment            int    `json:"default_environment,omitempty"`
	SignatureValidationCredential int    `json:"signature_validation_credential,omitempty"`
	LastUpdateFailed              bool   `json:"last_update_failed,omitempty"`
	LastUpdate                    string `json:"last_updated,omitempty"`
}

func (p Project) MarshalJSON() ([]byte, error) {
	type project Project
	return json.MarshalIndent((project)(p), "", "  ")
}

func (p *Project) ToBHNode() (node Node) {
	node.Kinds = []string{
		"ATProject",
	}
	node.Id = p.OID
	node.Properties = map[string]string{
		"id":                              strconv.Itoa(p.ID),
		"name":                            p.Name,
		"description":                     p.Description,
		"url":                             p.Url,
		"organization":                    strconv.FormatInt(int64(p.Organization), 10),
		"credential":                      strconv.FormatInt(int64(p.Credential), 10),
		"timeout":                         strconv.FormatInt(int64(p.Timeout), 10),
		"status":                          p.Status,
		"local_path":                      p.LocalPath,
		"scm_type":                        p.ScmType,
		"scm_url":                         p.ScmUrl,
		"scm_branch":                      p.ScmBranch,
		"scm_ref_spec":                    p.ScmRefSpec,
		"scm_clean":                       strconv.FormatBool(p.ScmClean),
		"scm_track_submodules":            strconv.FormatBool(p.ScmTrackSubmodules),
		"scm_delete_on_update":            strconv.FormatBool(p.ScmDeleteOnUpdate),
		"scm_revision":                    p.ScmRevision,
		"last_job_run":                    p.LastJobRun,
		"next_job_run":                    p.NextJobRun,
		"last_job_failed":                 strconv.FormatBool(p.LastJobFailed),
		"scm_update_on_launch":            strconv.FormatBool(p.ScmUpdateOnLaunch),
		"scm_update_cache_timeout":        strconv.FormatInt(int64(p.ScmUpdateCacheTimeout), 10),
		"allow_override":                  strconv.FormatBool(p.AllowOverride),
		"custom_virtualenv":               p.CustomVirtualenv,
		"default_environment":             strconv.FormatInt(int64(p.DefaultEnvironment), 10),
		"signature_validation_credential": strconv.FormatInt(int64(p.SignatureValidationCredential), 10),
		"last_update_failed":              strconv.FormatBool(p.LastUpdateFailed),
		"last_update":                     p.LastUpdate,
	}
	return node
}

type Credential struct {
	Object
	Organization   int            `json:"organization"`
	CredentialType int            `json:"credential_type,omitempty"`
	Managed        bool           `json:"managed,omitempty"`
	Inputs         map[string]any `json:"inputs,omitempty"`
	Cloud          bool           `json:"cloud,omitempty"`
	Kubernetes     bool           `json:"kubernetes,omitempty"`
	Kind           string         `json:"kind,omitempty"`
}

func (c Credential) MarshalJSON() ([]byte, error) {
	type credential Credential
	return json.MarshalIndent((credential)(c), "", "  ")
}

func (c *Credential) ToBHNode() (node Node) {
	node.Kinds = []string{
		"ATCredential",
	}
	node.Id = c.OID
	node.Properties = map[string]string{
		"id":              strconv.Itoa(c.ID),
		"name":            c.Name,
		"description":     c.Description,
		"url":             c.Url,
		"organization":    strconv.FormatInt(int64(c.Organization), 10),
		"credential_type": strconv.FormatInt(int64(c.CredentialType), 10),
		"managed":         strconv.FormatBool(c.Managed),
		"cloud":           strconv.FormatBool(c.Cloud),
		"kubernetes":      strconv.FormatBool(c.Kubernetes),
	}
	return node
}

type CredentialType struct {
	Object
	Managed    bool           `json:"managed,omitempty"`
	Inputs     map[string]any `json:"inputs,omitempty"`
	Injectors  map[string]any `json:"injectors,omitempty"`
	Cloud      bool           `json:"cloud,omitempty"`
	Kubernetes bool           `json:"kubernetes,omitempty"`
	Namespace  string         `json:"namespace,omitempty"`
	Kind       string         `json:"kind,omitempty"`
}

func (ct CredentialType) MarshalJSON() ([]byte, error) {
	type credentialType CredentialType
	return json.MarshalIndent((credentialType)(ct), "", "  ")
}

func (ct *CredentialType) ToBHNode() (node Node) {
	node.Kinds = []string{
		"ATCredentialType",
	}
	node.Id = ct.OID

	node.Properties = map[string]string{
		"id":          strconv.Itoa(ct.ID),
		"name":        ct.Name,
		"description": ct.Description,
		"url":         ct.Url,
		"namespace":   ct.Namespace,
		"managed":     strconv.FormatBool(ct.Managed),
		"cloud":       strconv.FormatBool(ct.Cloud),
		"kubernetes":  strconv.FormatBool(ct.Kubernetes),
	}

	var ok bool
	// TODO: Create object to manage `inputs`, the any is getting annoying to manage.
	// Most of these fields are documented and hardcoded.
	if _, ok = ct.Inputs["fields"]; ok {
		fields := ct.Inputs["fields"].([]any)
		for _, field := range fields {

			// There is type checking for these values on AWX/Tower's side.
			field := field.(map[string]any)
			id := field["id"].(string) // ID and Label are needed fields and should never be missing empty.

			if _, ok = field["help_text"]; ok {
				node.Properties["field_"+id+"_help"] = field["help_text"].(string)
			}
			if _, ok = field["type"]; ok {
				node.Properties["field_"+id+"_type"] = field["type"].(string)
			}
			if _, ok = field["secret"]; ok {
				node.Properties["field_"+id+"_secret"] = strconv.FormatBool(field["secret"].(bool))
			}
			if _, ok = field["help_text"]; ok {
				node.Properties["field_"+id+"_help"] = field["help_text"].(string)
			}
		}

		var requiredAny []any
		var required []string
		if _, ok = ct.Inputs["required"]; ok {
			requiredAny = ct.Inputs["required"].([]any)
			for _, entry := range requiredAny {
				entry := entry.(string)
				required = append(required, entry)
			}
		}
		node.Properties["fields_required"] = strings.Join(required, ", ")

	}

	if _, ok = ct.Injectors["file"]; ok {
		fileInjectors := ct.Injectors["file"].(map[string]any)
		for k, v := range fileInjectors {
			node.Properties["injector_file_"+k] = v.(string)
		}
	}
	if _, ok = ct.Injectors["extra_vars"]; ok {
		fileInjectors := ct.Injectors["extra_vars"].(map[string]any)
		for k, v := range fileInjectors {
			node.Properties["injector_extra_vars_"+k] = v.(string)
		}
	}
	if _, ok = ct.Injectors["env"]; ok {
		fileInjectors := ct.Injectors["env"].(map[string]any)
		for k, v := range fileInjectors {
			node.Properties["injector_env_"+k] = v.(string)
		}
	}
	return node
}

type Inventory struct {
	Object
	Organization                 int    `json:"organization"`
	Kind                         string `json:"kind,omitempty"`
	HostFilter                   string `json:"host_filder,omitempty"`
	Variables                    string `json:"variables,omitempty"`
	HasActiveFailures            bool   `json:"has_active_failures,omitempty"`
	TotalHosts                   int    `json:"total_hosts,omitempty"`
	HostsWithActiveFailures      int    `json:"host_with_active_failures,omitempty"`
	TotalGroups                  int    `json:"total_groups,omitempty"`
	HasInventorySources          bool   `json:"has_inventory_sources,omitempty"`
	TotalInventorySources        int    `json:"total_inventory_sources,omitempty"`
	InventorySourcesWithFailures int    `json:"inventory_sources_with_failures,omitempty"`
	PendingDeletion              bool   `json:"pending_deletion,omitempty"`
	PreventInstanceGroupFallback bool   `json:"prevent_instance_group_fallback,omitempty"`
}

func (i Inventory) MarshalJSON() ([]byte, error) {
	type inventory Inventory
	return json.MarshalIndent((inventory)(i), "", "  ")
}

func (i *Inventory) ToBHNode() (node Node) {
	node.Id = i.OID
	node.Kinds = []string{"ATInventory"}
	node.Properties = map[string]string{
		"id":                              strconv.FormatInt(int64(i.ID), 10),
		"name":                            i.Name,
		"description":                     i.Description,
		"url":                             i.Url,
		"organization":                    strconv.FormatInt(int64(i.Organization), 10),
		"kind":                            i.Kind,
		"host_filter":                     i.HostFilter,
		"has_active_failures":             strconv.FormatBool(i.HasActiveFailures),
		"has_inventory_source":            strconv.FormatBool(i.HasInventorySources),
		"total_hosts":                     strconv.FormatInt(int64(i.TotalHosts), 10),
		"hosts_with_active_failures":      strconv.FormatInt(int64(i.HostsWithActiveFailures), 10),
		"total_groups":                    strconv.FormatInt(int64(i.TotalGroups), 10),
		"total_inventory_sources":         strconv.FormatInt(int64(i.TotalInventorySources), 10),
		"inventory_sources_with_failures": strconv.FormatInt(int64(i.InventorySourcesWithFailures), 10),
		"pending_deletion":                strconv.FormatBool(i.PendingDeletion),
		"prevent_instance_group_fallback": strconv.FormatBool(i.PreventInstanceGroupFallback),
	}
	return node
}

type Host struct {
	Object
	Inventory            int    `json:"inventory,omitempty"`
	Enabled              bool   `json:"enabled,omitempty"`
	InstanceId           string `json:"instance_id,omitempty"`
	Variables            string `json:"variables,omitempty"`
	HasActiveFailures    bool   `json:"has_active_failures,omitempty"`
	LastJob              int    `json:"last_job,omitempty"`
	LastJobHostSummary   int    `json:"last_job_host_summary,omitempty"`
	AnsibleFactsModified string `json:"ansible_facts_modified,omitempty"`
}

func (i Host) MarshalJSON() ([]byte, error) {
	type host Host
	return json.MarshalIndent((host)(i), "", "  ")
}

func (h *Host) ToBHNode() (node Node) {
	node.Kinds = []string{
		"ATHost",
	}
	node.Id = h.OID
	node.Properties = map[string]string{
		"id":                     strconv.Itoa(h.ID),
		"name":                   h.Name,
		"description":            h.Description,
		"url":                    h.Url,
		"type":                   h.Type,
		"created":                h.Created,
		"modified":               h.Modified,
		"inventory":              strconv.FormatInt(int64(h.Inventory), 10),
		"enabled":                strconv.FormatBool(h.Enabled),
		"instance_id":            h.InstanceId,
		"variables":              h.Variables,
		"has_active_failures":    strconv.FormatBool(h.HasActiveFailures),
		"last_job":               strconv.FormatInt(int64(h.LastJob), 10),
		"last_job_host_summary":  strconv.FormatInt(int64(h.LastJobHostSummary), 10),
		"ansible_facts_modified": h.AnsibleFactsModified,
	}
	return node
}

// -- Future proofing, curernt RBAC APIs has been deprecated, it will eventually switch to these APIs --

type RoleDefinition struct {
	Object
	Permissions []string `json:"permissions"`
	ContentType string   `json:"content_type"`
	Managed     bool     `json:"managed"`
}

func (r RoleDefinition) MarshalJSON() ([]byte, error) {
	type roleDefinition RoleDefinition
	return json.MarshalIndent((roleDefinition)(r), "", "  ")
}

type RoleUserAssignments struct {
	Object
	ContentType    string `json:"content_type"`
	ObjectId       string `json:"object_id"`
	RoleDefinition int    `json:"role_definition"`
	UserId         int    `json:"user"`
}

func (r RoleUserAssignments) MarshalJSON() ([]byte, error) {
	type roleUserAssignments RoleUserAssignments
	return json.MarshalIndent((roleUserAssignments)(r), "", "  ")
}

type RoleTeamAssignments struct {
	Object
	ContentType    string `json:"content_type"`
	ObjectId       string `json:"object_id"`
	RoleDefinition int    `json:"role_definition"`
	TeamId         int    `json:"team"`
}

func (r RoleTeamAssignments) MarshalJSON() ([]byte, error) {
	type roleTeamAssignments RoleTeamAssignments
	return json.MarshalIndent((roleTeamAssignments)(r), "", "  ")
}

// -- Gestion du json de sortie --

type OutputJson struct {
	Metadata Metadata `json:"metadata"`
	Graph    Graph    `json:"graph"`
}

type Metadata struct {
	SourceKind string `json:"source_kind,omitempty"`
}

type Graph struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
}

type Node struct {
	Id         string            `json:"id"`
	Kinds      []string          `json:"kinds,omitempty"`
	Properties map[string]string `json:"properties"`
}

type Edge struct {
	Kind  string `json:"kind"`
	Start Link   `json:"start"`
	End   Link   `json:"end"`
}

type Link struct {
	Value   string `json:"value"`
	MatchBy string `json:"match_by"`
}
