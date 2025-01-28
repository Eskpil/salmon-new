package models

import "encoding/json"

const RootKey = "rockferry"

type ResourceKind string

const (
	ResourceKindNode          = "node"
	ResourceKindStoragePool   = "storagepool"
	ResourceKindStorageVolume = "storagevolume"
	ResourceKindNetwork       = "network"
	ResourceKindMachine       = "machine"
)

type OwnerRef struct {
	// The resource type, such as node
	Kind string `json:"kind"`
	Id   string `json:"id"`
}

type Resource struct {
	Id          string            `json:"id"`
	Kind        string            `json:"kind"`
	Annotations map[string]string `json:"annotations"`
	Owner       *OwnerRef         `json:"owner,omitempty"`
	Spec        interface{}       `json:"spec"`
}

// Should probobly propagte erros
func (r *Resource) Marshal() []byte {
	bytes, _ := json.Marshal(r)
	return bytes
}
