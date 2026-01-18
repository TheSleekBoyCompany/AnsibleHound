package ansible

import (
	"encoding/json"
	"strconv"

	"github.com/TheManticoreProject/gopengraph/node"
	"github.com/TheManticoreProject/gopengraph/properties"
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

func (u *User) ToBHNode() (n *node.Node) {
	props := properties.NewProperties()
	props.SetProperty("id", strconv.Itoa(u.ID))
	props.SetProperty("name", u.Username)
	props.SetProperty("description", u.Description)
	props.SetProperty("url", u.Url)
	props.SetProperty("firstname", u.FirstName)
	props.SetProperty("lastname", u.LastName)
	props.SetProperty("email", u.Email)
	props.SetProperty("is_super_user", strconv.FormatBool(u.IsSuperUser))
	props.SetProperty("is_system_auditor", strconv.FormatBool(u.IsSystemAuditor))
	props.SetProperty("ldap_dn", u.LdapDn)
	props.SetProperty("last_login", u.LastLogin)
	props.SetProperty("external_account", u.ExternalAccount)
	n, _ = node.NewNode(u.OID, []string{"ATUser"}, props)

	return n
}
