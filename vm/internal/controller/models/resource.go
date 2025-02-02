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

type Phase string

const (
	PhaseRequested = "requested"
	PhaseCreating  = "creating"
	PhaseCreated   = "created"
)

type Status struct {
	Phase Phase `json:"phase"`
}

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
	Status      Status            `json:"status"`
}

// Should probobly propagte erros
func (r *Resource) Marshal() []byte {
	bytes, _ := json.Marshal(r)
	return bytes
}
