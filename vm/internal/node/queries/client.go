package queries

import (
	"net/url"

	"github.com/digitalocean/go-libvirt"
)

type Client struct {
	v *libvirt.Libvirt
}

func NewClient() (*Client, error) {
	uri, _ := url.Parse(string(libvirt.QEMUSystem))
	l, err := libvirt.ConnectToURI(uri)
	if err != nil {
		return nil, err
	}
	client := &Client{v: l}

	return client, nil
}

func (c *Client) Destroy() {
	c.v.Disconnect()
}
