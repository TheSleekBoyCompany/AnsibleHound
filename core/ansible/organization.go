package ansible

import (
	"encoding/json"
	"strconv"

	"github.com/Ramoreik/gopengraph/node"
	"github.com/Ramoreik/gopengraph/properties"
)

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

func (o *Organization) ToBHNode() (n *node.Node) {
	props := properties.NewProperties()
	props.SetProperty("id", strconv.Itoa(o.ID))
	props.SetProperty("name", o.Name)
	props.SetProperty("description", o.Description)
	props.SetProperty("url", o.Url)
	props.SetProperty("max_hosts", strconv.FormatInt(int64(o.MaxHosts), 10))
	props.SetProperty("custom_virtualenv", o.CustomVirtualenv)
	props.SetProperty("default_environment", strconv.FormatInt(int64(o.DefaultEnvironment), 10))
	n, _ = node.NewNode(o.OID, []string{"ATOrganization"}, props)

	return n
}
