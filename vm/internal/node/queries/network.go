package queries

import (
	"encoding/xml"

	"github.com/digitalocean/go-libvirt"
	"github.com/eskpil/salmon/vm/pkg/rockferry"
	"github.com/eskpil/salmon/vm/pkg/rockferry/spec"
	"github.com/eskpil/salmon/vm/pkg/virtwrap/network"
)

func listAllNetworks(c *libvirt.Libvirt) ([]libvirt.Network, error) {
	networks, _, err := c.ConnectListAllNetworks(100, 1|2)
	return networks, err
}

func completeNetwork(c *libvirt.Libvirt, unmappedNetwork libvirt.Network) (*rockferry.Network, error) {
	xmlSchema, err := c.NetworkGetXMLDesc(unmappedNetwork, 0)
	if err != nil {
		return nil, err
	}

	var schema network.Schema
	if err := xml.Unmarshal([]byte(xmlSchema), &schema); err != nil {
		return nil, err
	}

	mapped := new(rockferry.Network)

	mapped.Spec.Name = schema.Name
	if schema.Ipv6 == "yes" {
		mapped.Spec.Ipv6 = true
	} else {
		mapped.Spec.Ipv6 = false
	}

	mapped.Id = schema.Uuid
	mapped.Spec.Mtu = uint64(schema.Mtu.Size)

	mapped.Spec.Bridge.Name = schema.Bridge.Name
	mapped.Spec.Bridge.Stp = schema.Bridge.Stp
	mapped.Spec.Bridge.Delay = schema.Bridge.Delay

	mapped.Spec.Forward.Dev = schema.Forward.Dev
	mapped.Spec.Forward.Mode = schema.Forward.Mode

	if schema.Forward.Nat != nil {
		mapped.Spec.Forward.Nat = new(spec.NetworkSpecForwardNat)

		mapped.Spec.Forward.Nat.AddressStart = schema.Forward.Nat.Address.Start
		mapped.Spec.Forward.Nat.AddressEnd = schema.Forward.Nat.Address.End

		mapped.Spec.Forward.Nat.PortEnd = schema.Forward.Nat.Port.End
		mapped.Spec.Forward.Nat.PortStart = schema.Forward.Nat.Port.Start
	}

	return mapped, nil
}

func (c *Client) ListAllNetworks() ([]*rockferry.Network, error) {
	unmappedNetworks, err := listAllNetworks(c.v)
	if err != nil {
		return nil, err
	}

	networks := make([]*rockferry.Network, len(unmappedNetworks))

	for i, unmappedNetwork := range unmappedNetworks {
		network, err := completeNetwork(c.v, unmappedNetwork)
		if err != nil {
			continue
		}

		networks[i] = network
	}

	return networks, nil
}
