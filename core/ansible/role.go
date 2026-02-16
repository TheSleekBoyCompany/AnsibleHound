package ansible

import (
	"encoding/json"

	"github.com/Ramoreik/gopengraph/node"
)

type Role struct {
	Object
	UserOnly      bool              `json:"user_only,omitempty"`
	SummaryFields RoleSummaryFields `json:"summary_fields"`
}

type RoleSummaryFields struct {
	ResourceName            string `json:"resource_name"`
	ResourceType            string `json:"resource_type"`
	ResourceTypeDisplayName string `json:"resource_type_display_name"`
	ResourceId              int    `json:"resource_id"`
}

func (r Role) MarshalJSON() ([]byte, error) {
	type role Role
	return json.MarshalIndent((role)(r), "", "  ")
}

func (r *Role) ToBHNode() (n *node.Node) {
	return n
}

// -- Future proofing, current RBAC APIs has been deprecated, it will eventually switch to these APIs --

type RoleDefinition struct {
	Object
	Permissions []string `json:"permissions"`
	ContentType string   `json:"content_type"`
	Managed     bool     `json:"managed"`
}

func (r RoleDefinition) MarshalJSON() ([]byte, error) {
	type roleDefinition RoleDefinition
	return json.MarshalIndent((roleDefinition)(r), "", "  ")
}

type RoleUserAssignments struct {
	Object
	ContentType    string `json:"content_type"`
	ObjectId       string `json:"object_id"`
	RoleDefinition int    `json:"role_definition"`
	UserId         int    `json:"user"`
}

func (r RoleUserAssignments) MarshalJSON() ([]byte, error) {
	type roleUserAssignments RoleUserAssignments
	return json.MarshalIndent((roleUserAssignments)(r), "", "  ")
}

type RoleTeamAssignments struct {
	Object
	ContentType    string `json:"content_type"`
	ObjectId       string `json:"object_id"`
	RoleDefinition int    `json:"role_definition"`
	TeamId         int    `json:"team"`
}

func (r RoleTeamAssignments) MarshalJSON() ([]byte, error) {
	type roleTeamAssignments RoleTeamAssignments
	return json.MarshalIndent((roleTeamAssignments)(r), "", "  ")
}
