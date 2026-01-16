package ansible

import (
	"ansible-hound/core/opengraph"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
)

var Instance AnsibleInstance

type Object struct {
	OID         string `json:"uuid,omitempty"`
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Url         string `json:"url,omitempty"`
	Type        string `json:"type,omitempty"`
	Created     string `json:"created,omitempty"`
	Modified    string `json:"modified,omitempty"`
}

func (o Object) GetOID() (uuid string) {
	return o.OID
}

func (o Object) GetID() (id int) {
	return o.ID
}

func (o *Object) InitOID(installUUID string) {
	data := fmt.Sprintf("%s_%s_%s", installUUID, strconv.Itoa(o.ID), o.Type)
	hasher := sha1.New()
	hasher.Write([]byte(data))
	hashBytes := hasher.Sum(nil)
	o.OID = hex.EncodeToString(hashBytes)
}

type AnsibleType interface {
	GetID() int
	GetOID() string
	InitOID(string)
	ToBHNode() opengraph.Node
}

type Response[T any] struct {
	Count   int `json:"count"`
	Results []T `json:"results"`
}

type AnsibleInstance struct {
	Object
	Version     string `json:"version"`
	ActiveNode  string `json:"active_node"`
	InstallUUID string `json:"install_uuid"`
}

func (i *AnsibleInstance) MarshalJSON() ([]byte, error) {
	type instance AnsibleInstance
	return json.MarshalIndent((*instance)(i), "", "  ")
}

func (i *AnsibleInstance) ToBHNode() (node opengraph.Node) {
	node.Kinds = []string{
		"ATAnsibleInstance",
	}
	node.Id = i.OID
	node.Properties = map[string]string{
		"name":         i.Name,
		"version":      i.Version,
		"active_node":  i.ActiveNode,
		"install_uuid": i.InstallUUID,
	}
	return node
}
