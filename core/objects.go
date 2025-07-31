package core

import "encoding/json"

type AnsibleTypeList interface {
	[]User | []Job | []JobTemplate | []Inventory |
		[]Project | []Organization | []Role | []Credential |
		[]Host | []Team | []RoleDefinition | []RoleTeamAssignments |
		[]RoleUserAssignments
}

type AnsibleType interface {
	User | Job | JobTemplate | Inventory |
		Project | Organization | Role | Credential |
		Host | Team | RoleDefinition | RoleTeamAssignments |
		RoleUserAssignments
}

type Object struct {
	UUID        string `json:"uuid,omitempty"`
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Url         string `json:"url,omitempty"`
	Type        string `json:"type,omitempty"`
	Created     string `json:"created,omitempty"`
	Modified    string `json:"modified,omitempty"`
}

type Response[T AnsibleType] struct {
	Count   int `json:"count"`
	Results []T `json:"results"`
}

type User struct {
	Object
	Username        string `json:"username"`
	FirstName       string `json:"first_name,omitempty"`
	LastName        string `json:"last_name,omitempty"`
	Email           string `json:"email,omitempty"`
	IsSuperUser     bool   `json:"is_superuser,omitempty"`
	IsSystemAuditor bool   `json:"is_sytem_auditor,omitempty"`
	LdapDn          string `json:"ldap_dn,omitempty"`
	LastLogin       string `json:"last_login,omitempty"`
	ExternalAccount string `json:"external_account,omitempty"`
	Roles           []Role `json:"roles"`
}

func (u *User) MarshalJSON() ([]byte, error) {
	type user User
	return json.MarshalIndent((*user)(u), "", "  ")
}

type Team struct {
	Object
	Organization int    `json:"organization,omitempty"`
	Members      []User `json:"members"`
	Roles        []Role `json:"roles"`
}

func (u *Team) MarshalJSON() ([]byte, error) {
	type team Team
	return json.MarshalIndent((*team)(u), "", "  ")
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

type JobTemplate struct {
	Object
	JobType                         string `json:"job_type"`
	Inventory                       int    `json:"inventory"`
	Project                         int    `json:"project"`
	Organization                    int    `json:"organization,omitempty"`
	Playbook                        string `json:"playbook"`
	SCMBranch                       string `json:"scm_branch,omitempty"`
	Limit                           string `json:"limit"`
	Verbosity                       int    `json:"verbosity"`
	ExtraVars                       string `json:"extra_vars"`
	Status                          string `json:"status,omitempty"`
	JobTags                         string `json:"job_tags,omitempty"`
	Forks                           int    `json:"forks"`
	SkipTags                        string `json:"skip_tags,omitempty"`
	StartAtTask                     string `json:"start_at_task,omitempty"`
	Timeout                         int    `json:"timeout,omitempty"`
	UseFactCache                    bool   `json:"use_fact_cache,omitempty"`
	ForceHandler                    bool   `json:"force_handlers,omitempty"`
	LastJobRun                      string `json:"last_job_run,omitempty"`
	NextJobRun                      string `json:"next_job_run,omitempty"`
	LastJobFailed                   bool   `json:"last_job_failed,omitempty"`
	ExecutionEnvironment            int    `json:"execution_environment,omitempty"`
	HostConfigKey                   string `json:"host_config_key,omitempty"`
	AskScmBranchOnLaunch            bool   `json:"ask_scm_branch_on_launch,omitempty"`
	AskDiffModeOnLaunch             bool   `json:"ask_diff_mode_on_launch,omitempty"`
	AskVariablesOnLaunch            bool   `json:"ask_variables_on_launch,omitempty"`
	AskLimitOnLaunch                bool   `json:"ask_limit_on_launch,omitempty"`
	AskTagsOnLaunch                 bool   `json:"ask_tags_on_launch,omitempty"`
	AskJobTypeOnLaunch              bool   `json:"ask_job_type_on_launch,omitempty"`
	AskVerbosityOnLaunch            bool   `json:"ask_verbosity_on_launch,omitempty"`
	AskInventoryOnLaunch            bool   `json:"ask_inventory_on_launch,omitempty"`
	AskCredentialOnLaunch           bool   `json:"ask_credential_on_launch,omitempty"`
	AskExecutionEnvironmentOnLaunch bool   `json:"ask_execution_environment_on_launch,omitempty"`
	AskLabelsOnLaunch               bool   `json:"ask_labels_on_launch,omitempty"`
	AskForksOnLaunch                bool   `json:"ask_forks_on_launch,omitempty"`
	AskJobSliceCountOnLaunch        bool   `json:"ask_job_slice_count_on_launch,omitempty"`
	AskTimeoutOnLaunch              bool   `json:"ask_timeout_on_launch,omitempty"`
	AskInstanceGroupsOnLaunch       bool   `json:"ask_instance_groups_on_launch,omitempty"`
	SurveyEnabled                   bool   `json:"survey_enabled,omitempty"`
	BecomeEnabled                   bool   `json:"become_enabled,omitempty"`
	DiffMode                        bool   `json:"diff_mode,omitempty"`
	AllowSimultaneous               bool   `json:"allow_simultaneous,omitempty"`
	CustomVirtualenv                string `json:"custom_virtualenv,omitempty"`
	JobSliceCount                   int    `json:"job_slice_count,omitempty"`
	WebhookService                  string `json:"webhook_service,omitempty"`
	WebhookCredential               int    `json:"webhook_credential,omitempty"`
	PreventInstanceGroupFallback    bool   `json:"prevent_instance_group_fallback,omitempty"`
}

func (j JobTemplate) MarshalJSON() ([]byte, error) {
	type jobTemplate JobTemplate
	return json.MarshalIndent((jobTemplate)(j), "", "  ")
}

type Job struct {
	Object
	Inventory             int                    `json:"inventory"`
	Project               int                    `json:"project"`
	Organization          int                    `json:"organization,omitempty"`
	Playbook              string                 `json:"playbook"`
	ScmBranch             string                 `json:"scm_branch,omitempty"`
	Forks                 int                    `json:"forks,omitempty"`
	Limit                 string                 `json:"limit,omitempty"`
	Verbosity             int                    `json:"verbosity,omitempty"`
	ExtraVars             string                 `json:"extra_vars,omitempty"`
	Started               string                 `json:"started,omitempty"`
	Finished              string                 `json:"finished,omitempty"`
	CanceledOn            string                 `json:"canceled_on,omitempty"`
	Elapsed               float32                `json:"elapsed,omitempty"`
	JobExplanation        string                 `json:"job_explanation,omitempty"`
	Created               string                 `json:"created,omitempty"`
	Modified              string                 `json:"modified,omitempty"`
	UnifiedJobTemplate    int                    `json:"unified_job_template"`
	LaunchType            string                 `json:"launch_type"`
	Failed                bool                   `json:"failed"`
	Status                string                 `json:"status,omitempty"`
	ExecutionEnvironment  int                    `json:"execution_environment,omitempty"`
	ExecutionNode         string                 `json:"execution_node,omitempty"`
	ControllerNode        string                 `json:"controller_node,omitempty"`
	LaunchedBy            map[string]interface{} `json:"launched_by,omitempty"`
	WorkUnitId            string                 `json:"work_unit_id,omitempty"`
	JobTags               string                 `json:"job_tags,omitempty"`
	JobType               string                 `json:"job_type,omitempty"`
	ForceHandler          bool                   `json:"force_handlers,omitempty"`
	SkipTags              string                 `json:"skip_tags,omitempty"`
	StartAtTask           string                 `json:"start_at_task,omitempty"`
	Timeout               int                    `json:"timeout,omitempty"`
	UseFactCache          bool                   `json:"use_fact_cache,omitempty"`
	PasswordNeededToStart string                 `json:"password_needed_to_start,omitempty"`
	AllowSimultaneous     bool                   `json:"allow_simultaneous,omitempty"`
	Artifacts             map[string]interface{} `json:"artifacts,omitempty"`
	ScmRevision           string                 `json:"scm_revision,omitempty"`
	InstanceGroup         int                    `json:"instance_group,omitempty"`
	DiffMode              bool                   `json:"diff_mode,omitempty"`
	JobSliceNumber        int                    `json:"job_slice_number,omitempty"`
	JobSliceCount         int                    `json:"job_slice_count,omitempty"`
	WebhookGuid           string                 `json:"webhook_guid,omitempty"`
	WebhookService        string                 `json:"webhook_service,omitempty"`
	WebhookCredential     int                    `json:"webhook_credential,omitempty"`
}

func (j Job) MarshalJSON() ([]byte, error) {
	type job Job
	return json.MarshalIndent((job)(j), "", "  ")
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
