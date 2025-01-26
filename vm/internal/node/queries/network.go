package queries

import (
	"encoding/json"
	"encoding/xml"

	"github.com/digitalocean/go-libvirt"
	"github.com/eskpil/salmon/vm/internal/node/virtwrap/network"
	"github.com/eskpil/salmon/vm/nodeapi"
)

func listAllNetworks(c *libvirt.Libvirt) ([]libvirt.Network, error) {
	networks, _, err := c.ConnectListAllNetworks(100, 1|2)
	return networks, err
}

func completeNetwork(c *libvirt.Libvirt, unmappedNetwork libvirt.Network) (*nodeapi.Network, error) {
	xmlSchema, err := c.NetworkGetXMLDesc(unmappedNetwork, 0)
	if err != nil {
		return nil, err
	}

	var schema network.Schema
	if err := xml.Unmarshal([]byte(xmlSchema), &schema); err != nil {
		return nil, err
	}

	schemaJson, err := json.Marshal(schema)
	if err != nil {
		return nil, err
	}

	mapped := new(nodeapi.Network)

	mapped.Name = schema.Name
	if schema.Ipv6 == "yes" {
		mapped.Ipv6 = true
	} else {
		mapped.Ipv6 = false
	}

	mapped.Bridge = new(nodeapi.NetworkBridge)
	mapped.Forward = new(nodeapi.NetworkForward)

	mapped.Uuid = schema.Uuid
	mapped.Mtu = uint64(schema.Mtu.Size)
	mapped.Schema = schemaJson

	mapped.Bridge.Name = schema.Bridge.Name
	mapped.Bridge.Stp = schema.Bridge.Stp
	mapped.Bridge.Delay = schema.Bridge.Delay

	mapped.Forward.Dev = schema.Forward.Dev
	mapped.Forward.Mode = schema.Forward.Mode

	if schema.Forward.Nat != nil {
		mapped.Forward.Nat = new(nodeapi.NetworkForwardNat)

		mapped.Forward.Nat.AddressStart = schema.Forward.Nat.Address.Start
		mapped.Forward.Nat.AddressEnd = schema.Forward.Nat.Address.End

		mapped.Forward.Nat.PortEnd = schema.Forward.Nat.Port.End
		mapped.Forward.Nat.PortStart = schema.Forward.Nat.Port.Start
	}

	return mapped, nil
}

func (c *Client) ListAllNetworks() ([]*nodeapi.Network, error) {
	unmappedNetworks, err := listAllNetworks(c.v)
	if err != nil {
		return nil, err
	}

	networks := make([]*nodeapi.Network, len(unmappedNetworks))

	for i, unmappedNetwork := range unmappedNetworks {
		network, err := completeNetwork(c.v, unmappedNetwork)
		if err != nil {
			continue
		}

		networks[i] = network
	}

	return networks, nil
}
