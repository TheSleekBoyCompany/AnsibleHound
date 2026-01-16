package ansible

import (
	"ansible-hound/core/opengraph"
	"encoding/json"
	"strconv"
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

func (p *Project) ToBHNode() (node opengraph.Node) {
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
