package ansible

import (
	"ansible-hound/core/opengraph"
	"encoding/json"
	"strconv"
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

func (i *Inventory) ToBHNode() (node opengraph.Node) {
	node.Id = i.OID
	node.Kinds = []string{"ATInventory"}
	node.Properties = map[string]string{
		"id":                              strconv.FormatInt(int64(i.ID), 10),
		"name":                            i.Name,
		"description":                     i.Description,
		"url":                             i.Url,
		"organization":                    strconv.FormatInt(int64(i.Organization), 10),
		"kind":                            i.Kind,
		"host_filter":                     i.HostFilter,
		"has_active_failures":             strconv.FormatBool(i.HasActiveFailures),
		"has_inventory_source":            strconv.FormatBool(i.HasInventorySources),
		"total_hosts":                     strconv.FormatInt(int64(i.TotalHosts), 10),
		"hosts_with_active_failures":      strconv.FormatInt(int64(i.HostsWithActiveFailures), 10),
		"total_groups":                    strconv.FormatInt(int64(i.TotalGroups), 10),
		"total_inventory_sources":         strconv.FormatInt(int64(i.TotalInventorySources), 10),
		"inventory_sources_with_failures": strconv.FormatInt(int64(i.InventorySourcesWithFailures), 10),
		"pending_deletion":                strconv.FormatBool(i.PendingDeletion),
		"prevent_instance_group_fallback": strconv.FormatBool(i.PreventInstanceGroupFallback),
	}
	return node
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

func (h *Host) ToBHNode() (node opengraph.Node) {
	node.Kinds = []string{
		"ATHost",
	}
	node.Id = h.OID
	node.Properties = map[string]string{
		"id":                     strconv.Itoa(h.ID),
		"name":                   h.Name,
		"description":            h.Description,
		"url":                    h.Url,
		"type":                   h.Type,
		"created":                h.Created,
		"modified":               h.Modified,
		"inventory":              strconv.FormatInt(int64(h.Inventory), 10),
		"enabled":                strconv.FormatBool(h.Enabled),
		"instance_id":            h.InstanceId,
		"variables":              h.Variables,
		"has_active_failures":    strconv.FormatBool(h.HasActiveFailures),
		"last_job":               strconv.FormatInt(int64(h.LastJob), 10),
		"last_job_host_summary":  strconv.FormatInt(int64(h.LastJobHostSummary), 10),
		"ansible_facts_modified": h.AnsibleFactsModified,
	}
	return node
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

func (g *Group) ToBHNode() (node opengraph.Node) {
	node.Kinds = []string{
		"ATGroup",
	}
	node.Id = g.OID
	node.Properties = map[string]string{
		"id":          strconv.Itoa(g.ID),
		"name":        g.Name,
		"description": g.Description,
		"url":         g.Url,
		"type":        g.Type,
		"created":     g.Created,
		"modified":    g.Modified,
		"inventory":   strconv.FormatInt(int64(g.Inventory), 10),
		"variables":   g.Variables,
	}
	return node
}
