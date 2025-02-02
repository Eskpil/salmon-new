package rockferry

import (
	"github.com/eskpil/salmon/vm/controllerapi"
	"github.com/eskpil/salmon/vm/pkg/rockferry/transport"
	genericv1 "github.com/eskpil/salmon/vm/pkg/rockferry/v1/generic"
	machinesv1 "github.com/eskpil/salmon/vm/pkg/rockferry/v1/machines"
	nodesv1 "github.com/eskpil/salmon/vm/pkg/rockferry/v1/nodes"
	storagevolumesv1 "github.com/eskpil/salmon/vm/pkg/rockferry/v1/storagevolumes"
)

type Client struct {
	c *controllerapi.ControllerApiClient

	genericv1        *genericv1.Interface
	nodesv1          *nodesv1.Interface
	storagevolumesv1 *storagevolumesv1.Interface
	machinesv1       *machinesv1.Interface
}

func New() (*Client, error) {
	transport, err := transport.New("10.100.0.102:9090")
	if err != nil {
		return nil, err
	}

	return &Client{
		nodesv1:          nodesv1.New(transport),
		storagevolumesv1: storagevolumesv1.New(transport),
		genericv1:        genericv1.New(transport),
		machinesv1:       machinesv1.New(transport),
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
