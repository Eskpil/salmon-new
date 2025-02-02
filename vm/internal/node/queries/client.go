package queries

import (
	"errors"
	"net/url"

	"github.com/digitalocean/go-libvirt"
	"github.com/eskpil/salmon/vm/nodeapi"
)

type Client struct {
	v *libvirt.Libvirt

	pools   []*nodeapi.StoragePool
	volumes []*nodeapi.StorageVolume
}

func NewClient() (*Client, error) {
	uri, _ := url.Parse(string(libvirt.QEMUSystem))
	l, err := libvirt.ConnectToURI(uri)
	if err != nil {
		return nil, err
	}
	client := &Client{v: l}

	return client, client.preload()
}

func (c *Client) preload() error {
	return c.preloadStorage()
}

func (c *Client) findVolumeByKey(key string) (*nodeapi.StorageVolume, error) {
	for _, vol := range c.volumes {
		if vol.Key == key {
			return vol, nil
		}
	}

	return nil, errors.New("could not find volume with key")
}

func (c *Client) findPoolByUuid(uuid string) (*nodeapi.StoragePool, error) {
	for _, pool := range c.pools {
		if pool.Uuid == uuid {
			return pool, nil
		}
	}

	return nil, errors.New("could not find pool with key")
}

func (c *Client) Destroy() {
	c.v.Disconnect()
}
