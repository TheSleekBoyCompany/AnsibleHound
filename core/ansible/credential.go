package ansible

import (
	"ansible-hound/core/opengraph"
	"encoding/json"
	"strconv"
	"strings"
)

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

func (c *Credential) ToBHNode() (node opengraph.Node) {
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

func (ct *CredentialType) ToBHNode() (node opengraph.Node) {
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
