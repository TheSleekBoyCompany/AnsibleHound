package ansible

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/Ramoreik/gopengraph/node"
	"github.com/Ramoreik/gopengraph/properties"
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

func (c *Credential) ToBHNode() (n *node.Node) {
	props := properties.NewProperties()
	props.SetProperty("id", strconv.Itoa(c.ID))
	props.SetProperty("name", c.Name)
	props.SetProperty("description", c.Description)
	props.SetProperty("url", c.Url)
	props.SetProperty("organization", strconv.FormatInt(int64(c.Organization), 10))
	props.SetProperty("credential_type", strconv.FormatInt(int64(c.CredentialType), 10))
	props.SetProperty("managed", strconv.FormatBool(c.Managed))
	props.SetProperty("cloud", strconv.FormatBool(c.Cloud))
	props.SetProperty("kubernetes", strconv.FormatBool(c.Kubernetes))

	if c.CredentialType == 1 { // Credential is of type Machine

		var ok bool
		if _, ok = c.Inputs["username"]; ok {
			username := c.Inputs["username"].(string)
			props.SetProperty("username", username)
		}

		var sshKeyDefined bool
		var password bool
		_, sshKeyDefined = c.Inputs["ssh_key_data"]
		_, password = c.Inputs["password"]
		if password && sshKeyDefined {
			props.SetProperty("machine_credential_type", "both")
		}
		if password && !sshKeyDefined {
			props.SetProperty("machine_credential_type", "password")
		}
		if !password && sshKeyDefined {
			props.SetProperty("machine_credential_type", "ssh")
		}
	}

	n, _ = node.NewNode(c.OID, []string{"ATCredential"}, props)

	return n
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

func (ct *CredentialType) ToBHNode() (n *node.Node) {
	props := properties.NewProperties()
	props.SetProperty("id", strconv.Itoa(ct.ID))
	props.SetProperty("name", ct.Name)
	props.SetProperty("description", ct.Description)
	props.SetProperty("url", ct.Url)
	props.SetProperty("namespace", ct.Namespace)
	props.SetProperty("managed", strconv.FormatBool(ct.Managed))
	props.SetProperty("cloud", strconv.FormatBool(ct.Cloud))
	props.SetProperty("kubernetes", strconv.FormatBool(ct.Kubernetes))

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
				props.SetProperty("field_"+id+"_help", field["help_text"].(string))
			}
			if _, ok = field["type"]; ok {
				props.SetProperty("field_"+id+"_type", field["type"].(string))
			}
			if _, ok = field["secret"]; ok {
				props.SetProperty("field_"+id+"_secret", strconv.FormatBool(field["secret"].(bool)))
			}
			if _, ok = field["help_text"]; ok {
				props.SetProperty("field_"+id+"_help", field["help_text"].(string))
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
		props.SetProperty("fields_required", strings.Join(required, ", "))

	}

	if _, ok = ct.Injectors["file"]; ok {
		fileInjectors := ct.Injectors["file"].(map[string]any)
		for k, v := range fileInjectors {
			props.SetProperty("injector_file_"+k, v.(string))
		}
	}
	if _, ok = ct.Injectors["extra_vars"]; ok {
		props.SetProperty("injector_extra_vars", true)
		fileInjectors := ct.Injectors["extra_vars"].(map[string]any)
		for k, v := range fileInjectors {
			props.SetProperty("injector_extra_vars_"+k, v.(string))
		}
	}
	if _, ok = ct.Injectors["env"]; ok {
		props.SetProperty("injector_env", true)
		fileInjectors := ct.Injectors["env"].(map[string]any)
		for k, v := range fileInjectors {
			props.SetProperty("injector_env_"+k, v.(string))
		}
	}
	n, _ = node.NewNode(ct.OID, []string{"ATCredentialType"}, props)

	return n
}
