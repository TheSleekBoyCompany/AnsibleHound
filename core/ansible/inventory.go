package ansible

import (
	"encoding/json"
	"strconv"

	"github.com/Ramoreik/gopengraph/node"
	"github.com/Ramoreik/gopengraph/properties"
)

type Inventory struct {
	Object
	Organization                 int    `json:"organization"`
	Kind                         string `json:"kind,omitempty"`
	HostFilter                   string `json:"host_filder,omitempty"`
	Variables                    string `json:"variables,omitempty"`
	HasActiveFailures            bool   `json:"has_active_failures,omitempty"`
	TotalHosts                   int    `json:"total_hosts,omitempty"`
	HostsWithActiveFailures      int    `json:"host_with_active_failures,omitempty"`
	TotalGroups                  int    `json:"total_groups,omitempty"`
	HasInventorySources          bool   `json:"has_inventory_sources,omitempty"`
	TotalInventorySources        int    `json:"total_inventory_sources,omitempty"`
	InventorySourcesWithFailures int    `json:"inventory_sources_with_failures,omitempty"`
	PendingDeletion              bool   `json:"pending_deletion,omitempty"`
	PreventInstanceGroupFallback bool   `json:"prevent_instance_group_fallback,omitempty"`
}

func (i Inventory) MarshalJSON() ([]byte, error) {
	type inventory Inventory
	return json.MarshalIndent((inventory)(i), "", "  ")
}

func (i *Inventory) ToBHNode() (n *node.Node) {
	props := properties.NewProperties()
	props.SetProperty("id", strconv.FormatInt(int64(i.ID), 10))
	props.SetProperty("name", i.Name)
	props.SetProperty("description", i.Description)
	props.SetProperty("url", i.Url)
	props.SetProperty("organization", strconv.FormatInt(int64(i.Organization), 10))
	props.SetProperty("kind", i.Kind)
	props.SetProperty("host_filter", i.HostFilter)
	props.SetProperty("has_active_failures", strconv.FormatBool(i.HasActiveFailures))
	props.SetProperty("has_inventory_source", strconv.FormatBool(i.HasInventorySources))
	props.SetProperty("total_hosts", strconv.FormatInt(int64(i.TotalHosts), 10))
	props.SetProperty("hosts_with_active_failures", strconv.FormatInt(int64(i.HostsWithActiveFailures), 10))
	props.SetProperty("total_groups", strconv.FormatInt(int64(i.TotalGroups), 10))
	props.SetProperty("total_inventory_sources", strconv.FormatInt(int64(i.TotalInventorySources), 10))
	props.SetProperty("inventory_sources_with_failures", strconv.FormatInt(int64(i.InventorySourcesWithFailures), 10))
	props.SetProperty("pending_deletion", strconv.FormatBool(i.PendingDeletion))
	props.SetProperty("prevent_instance_group_fallback", strconv.FormatBool(i.PreventInstanceGroupFallback))
	n, _ = node.NewNode(i.OID, []string{"ATInventory"}, props)

	return n
}

type Host struct {
	Object
	Inventory            int    `json:"inventory,omitempty"`
	Enabled              bool   `json:"enabled,omitempty"`
	InstanceId           string `json:"instance_id,omitempty"`
	Variables            string `json:"variables,omitempty"`
	HasActiveFailures    bool   `json:"has_active_failures,omitempty"`
	LastJob              int    `json:"last_job,omitempty"`
	LastJobHostSummary   int    `json:"last_job_host_summary,omitempty"`
	AnsibleFactsModified string `json:"ansible_facts_modified,omitempty"`
}

func (i Host) MarshalJSON() ([]byte, error) {
	type host Host
	return json.MarshalIndent((host)(i), "", "  ")
}

func (h *Host) ToBHNode() (n *node.Node) {
	props := properties.NewProperties()
	props.SetProperty("id", strconv.Itoa(h.ID))
	props.SetProperty("name", h.Name)
	props.SetProperty("description", h.Description)
	props.SetProperty("url", h.Url)
	props.SetProperty("type", h.Type)
	props.SetProperty("created", h.Created)
	props.SetProperty("modified", h.Modified)
	props.SetProperty("inventory", strconv.FormatInt(int64(h.Inventory), 10))
	props.SetProperty("enabled", strconv.FormatBool(h.Enabled))
	props.SetProperty("instance_id", h.InstanceId)
	props.SetProperty("variables", h.Variables)
	props.SetProperty("has_active_failures", strconv.FormatBool(h.HasActiveFailures))
	props.SetProperty("last_job", strconv.FormatInt(int64(h.LastJob), 10))
	props.SetProperty("last_job_host_summary", strconv.FormatInt(int64(h.LastJobHostSummary), 10))
	props.SetProperty("ansible_facts_modified", h.AnsibleFactsModified)
	n, _ = node.NewNode(h.OID, []string{"ATHost"}, props)

	return n
}

type Group struct {
	Object
	Inventory int           `json:"inventory,omitempty"`
	Variables string        `json:"variables,omitempty"`
	Hosts     map[int]*Host `json:"hosts,omitempty"`
}

func (i Group) MarshalJSON() ([]byte, error) {
	type group Group
	return json.MarshalIndent((group)(i), "", "  ")
}

func (g *Group) ToBHNode() (n *node.Node) {
	props := properties.NewProperties()
	props.SetProperty("id", strconv.Itoa(g.ID))
	props.SetProperty("name", g.Name)
	props.SetProperty("description", g.Description)
	props.SetProperty("url", g.Url)
	props.SetProperty("type", g.Type)
	props.SetProperty("created", g.Created)
	props.SetProperty("modified", g.Modified)
	props.SetProperty("inventory", strconv.FormatInt(int64(g.Inventory), 10))
	props.SetProperty("variables", g.Variables)
	n, _ = node.NewNode(g.OID, []string{"ATGroup"}, props)

	return n
}
