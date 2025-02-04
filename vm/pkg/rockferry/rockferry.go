package rockferry

import (
	"github.com/eskpil/salmon/vm/controllerapi"
	"github.com/eskpil/salmon/vm/pkg/rockferry/resource"
	"github.com/eskpil/salmon/vm/pkg/rockferry/spec"
	"github.com/eskpil/salmon/vm/pkg/rockferry/transport"
)

type WatchAction = int

const (
	WatchActionPut WatchAction = iota
	WatchActionDelete
	WatchActionAll
)

type MachineRequest = resource.Resource[spec.MachineRequestSpec]
type StorageVolume = resource.Resource[spec.StorageVolumeSpec]
type StoragePool = resource.Resource[spec.StoragePoolSpec]
type Node = resource.Resource[spec.NodeSpec]
type Network = resource.Resource[spec.NetworkSpec]
type Machine = resource.Resource[spec.MachineSpec]

type Client struct {
	c *controllerapi.ControllerApiClient
	t *transport.Transport

	nodesv1            *Interface[spec.NodeSpec]
	storagevolumesv1   *Interface[spec.StorageVolumeSpec]
	machinesv1         *Interface[spec.MachineSpec]
	machinesrequestsv1 *Interface[spec.MachineRequestSpec]
	networksv1         *Interface[spec.NetworkSpec]
	storagepoolsv1     *Interface[spec.StoragePoolSpec]
}

func New(url string) (*Client, error) {
	transport, err := transport.New(url)
	if err != nil {
		return nil, err
	}

	return &Client{
		nodesv1:            NewInterface[spec.NodeSpec](resource.ResourceKindNode, transport),
		storagevolumesv1:   NewInterface[spec.StorageVolumeSpec](resource.ResourceKindStorageVolume, transport),
		machinesv1:         NewInterface[spec.MachineSpec](resource.ResourceKindMachine, transport),
		machinesrequestsv1: NewInterface[spec.MachineRequestSpec](resource.ResourceKindMachineRequest, transport),
		networksv1:         NewInterface[spec.NetworkSpec](resource.ResourceKindNetwork, transport),
		storagepoolsv1:     NewInterface[spec.StoragePoolSpec](resource.ResourceKindStoragePool, transport),

		t: transport,
	}, nil
}

func (c *Client) Nodes() *Interface[spec.NodeSpec] {
	return c.nodesv1
}

func (c *Client) StorageVolumes() *Interface[spec.StorageVolumeSpec] {
	return c.storagevolumesv1
}

func (c *Client) Generic(kind resource.ResourceKind) *Interface[any] {
	return NewInterface[any](kind, c.t)
}

func (c *Client) Machines() *Interface[spec.MachineSpec] {
	return c.machinesv1
}

func (c *Client) MachineRequests() *Interface[spec.MachineRequestSpec] {
	return c.machinesrequestsv1
}

func (c *Client) Networks() *Interface[spec.NetworkSpec] {
	return c.networksv1
}

func (c *Client) StoragePools() *Interface[spec.StoragePoolSpec] {
	return c.storagepoolsv1
}
