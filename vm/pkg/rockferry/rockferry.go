package rockferry

import (
	"github.com/eskpil/salmon/vm/controllerapi"
	"github.com/eskpil/salmon/vm/pkg/rockferry/transport"
	genericv1 "github.com/eskpil/salmon/vm/pkg/rockferry/v1/generic"
	machinerequestsv1 "github.com/eskpil/salmon/vm/pkg/rockferry/v1/machinerequests"
	machinesv1 "github.com/eskpil/salmon/vm/pkg/rockferry/v1/machines"
	networksv1 "github.com/eskpil/salmon/vm/pkg/rockferry/v1/networks"
	nodesv1 "github.com/eskpil/salmon/vm/pkg/rockferry/v1/nodes"
	storagepoolsv1 "github.com/eskpil/salmon/vm/pkg/rockferry/v1/storagepools"
	storagevolumesv1 "github.com/eskpil/salmon/vm/pkg/rockferry/v1/storagevolumes"
)

// TODO: Discover a solution with generics to align with DRY (do not repeat yourself) and KISS (keep it simple stupid)
// 		 should be relativly simple, make Interface take a generic type for the Spec and pass in the resource kind as a
// 		 argument.

type Client struct {
	c *controllerapi.ControllerApiClient

	genericv1         *genericv1.Interface
	nodesv1           *nodesv1.Interface
	storagevolumesv1  *storagevolumesv1.Interface
	machinesv1        *machinesv1.Interface
	machinerequestsv1 *machinerequestsv1.Interface
	networksv1        *networksv1.Interface
	storagepoolsv1    *storagepoolsv1.Interface
}

func New() (*Client, error) {
	transport, err := transport.New("10.100.0.102:9090")
	if err != nil {
		return nil, err
	}

	return &Client{
		nodesv1:           nodesv1.New(transport),
		storagevolumesv1:  storagevolumesv1.New(transport),
		genericv1:         genericv1.New(transport),
		machinesv1:        machinesv1.New(transport),
		machinerequestsv1: machinerequestsv1.New(transport),
		networksv1:        networksv1.New(transport),
		storagepoolsv1:    storagepoolsv1.New(transport),
	}, nil
}

func (c *Client) Nodes() *nodesv1.Interface {
	return c.nodesv1
}

func (c *Client) StorageVolumes() *storagevolumesv1.Interface {
	return c.storagevolumesv1
}

func (c *Client) Generic() *genericv1.Interface {
	return c.genericv1
}

func (c *Client) Machines() *machinesv1.Interface {
	return c.machinesv1
}

func (c *Client) MachineRequests() *machinerequestsv1.Interface {
	return c.machinerequestsv1
}

func (c *Client) Networks() *networksv1.Interface {
	return c.networksv1
}

func (c *Client) StoragePools() *storagepoolsv1.Interface {
	return c.storagepoolsv1
}
