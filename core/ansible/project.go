package ansible

import (
	"encoding/json"
	"strconv"

	"github.com/Ramoreik/gopengraph/node"
	"github.com/Ramoreik/gopengraph/properties"
)

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

func (p *Project) ToBHNode() (n *node.Node) {
	props := properties.NewProperties()
	props.SetProperty("id", strconv.Itoa(p.ID))
	props.SetProperty("name", p.Name)
	props.SetProperty("description", p.Description)
	props.SetProperty("url", p.Url)
	props.SetProperty("organization", strconv.FormatInt(int64(p.Organization), 10))
	props.SetProperty("credential", strconv.FormatInt(int64(p.Credential), 10))
	props.SetProperty("timeout", strconv.FormatInt(int64(p.Timeout), 10))
	props.SetProperty("status", p.Status)
	props.SetProperty("local_path", p.LocalPath)
	props.SetProperty("scm_type", p.ScmType)
	props.SetProperty("scm_url", p.ScmUrl)
	props.SetProperty("scm_branch", p.ScmBranch)
	props.SetProperty("scm_ref_spec", p.ScmRefSpec)
	props.SetProperty("scm_clean", strconv.FormatBool(p.ScmClean))
	props.SetProperty("scm_track_submodules", strconv.FormatBool(p.ScmTrackSubmodules))
	props.SetProperty("scm_delete_on_update", strconv.FormatBool(p.ScmDeleteOnUpdate))
	props.SetProperty("scm_revision", p.ScmRevision)
	props.SetProperty("last_job_run", p.LastJobRun)
	props.SetProperty("next_job_run", p.NextJobRun)
	props.SetProperty("last_job_failed", strconv.FormatBool(p.LastJobFailed))
	props.SetProperty("scm_update_on_launch", strconv.FormatBool(p.ScmUpdateOnLaunch))
	props.SetProperty("scm_update_cache_timeout", strconv.FormatInt(int64(p.ScmUpdateCacheTimeout), 10))
	props.SetProperty("allow_override", strconv.FormatBool(p.AllowOverride))
	props.SetProperty("custom_virtualenv", p.CustomVirtualenv)
	props.SetProperty("default_environment", strconv.FormatInt(int64(p.DefaultEnvironment), 10))
	props.SetProperty("signature_validation_credential", strconv.FormatInt(int64(p.SignatureValidationCredential), 10))
	props.SetProperty("last_update_failed", strconv.FormatBool(p.LastUpdateFailed))
	props.SetProperty("last_update", p.LastUpdate)
	n, _ = node.NewNode(p.OID, []string{"ATProject"}, props)

	return n
}
