package ansible

import (
	"ansible-hound/core/opengraph"
	"encoding/json"
	"strconv"
)

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

func (u *User) ToBHNode() (node opengraph.Node) {
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
