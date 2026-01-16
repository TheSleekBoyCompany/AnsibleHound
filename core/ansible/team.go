package ansible

import (
	"ansible-hound/core/opengraph"
	"encoding/json"
	"strconv"
)

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

func (t *Team) ToBHNode() (node opengraph.Node) {
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
