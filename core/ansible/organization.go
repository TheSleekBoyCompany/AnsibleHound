package ansible

import (
	"ansible-hound/core/opengraph"
	"encoding/json"
	"strconv"
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

func (o *Organization) ToBHNode() (node opengraph.Node) {
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
