package resource

import "google.golang.org/protobuf/types/known/structpb"

type ResourceKind = string

const (
	ResourceKindAll            = ""
	ResourceKindNode           = "node"
	ResourceKindStoragePool    = "storagepool"
	ResourceKindStorageVolume  = "storagevolume"
	ResourceKindNetwork        = "network"
	ResourceKindMachine        = "machine"
	ResourceKindMachineRequest = "machinerequest"
)

type Phase string

const (
	PhaseRequested = "requested"
	PhaseCreating  = "creating"
	PhaseErrored   = "errored"
	PhaseCreated   = "created"
)

// Used by node and machine
type Topology struct {
	Sockets uint64 `json:"sockets"`
	Cores   uint64 `json:"cores"`
	Threads uint64 `json:"threads"`
	Memory  uint64 `json:"memory"`
}

type Status struct {
	Phase Phase   `json:"phase"`
	Error *string `json:"error"`
}

type OwnerRef struct {
	// The resource type, such as node
	Kind string `json:"kind"`
	Id   string `json:"id"`
}

type Resource[T any] struct {
	Id          string            `json:"id"`
	Kind        ResourceKind      `json:"kind"`
	Annotations map[string]string `json:"annotations"`
	Owner       *OwnerRef         `json:"owner,omitempty"`
	Spec        T                 `json:"spec"`
	Status      Status            `json:"status"`

	RawSpec *structpb.Struct `json:"-"`
}

func (r *Resource[T]) Merge(with *Resource[T]) {
	// TODO: More fields possibily?
	for k, v := range with.Annotations {
		r.Annotations[k] = v
	}
}

func (r *Resource[T]) Generic() *Resource[any] {
	var spec interface{}
	spec = r.Spec

	return &Resource[any]{
		Id:          r.Id,
		Kind:        r.Kind,
		Annotations: r.Annotations,
		Owner:       r.Owner,
		Spec:        &spec, // Store spec as interface{}
		Status:      r.Status,
	}
}
