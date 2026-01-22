package ansible

import (
	"encoding/json"
	"strconv"

	"github.com/TheManticoreProject/gopengraph/node"
	"github.com/TheManticoreProject/gopengraph/properties"
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

func (t *Team) ToBHNode() (n *node.Node) {
	props := properties.NewProperties()
	props.SetProperty("id", strconv.Itoa(t.ID))
	props.SetProperty("name", t.Name)
	props.SetProperty("description", t.Description)
	props.SetProperty("url", t.Url)
	props.SetProperty("type", t.Type)
	props.SetProperty("created", t.Created)
	props.SetProperty("modified", t.Modified)
	n, _ = node.NewNode(t.OID, []string{"ATTeam"}, props)

	return n
}
