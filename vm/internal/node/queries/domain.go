package queries

import (
	"encoding/json"
	"encoding/xml"
	"fmt"

	"github.com/digitalocean/go-libvirt"
	"github.com/eskpil/salmon/vm/internal/node/virtwrap/domain"
	"github.com/eskpil/salmon/vm/nodeapi"
)

func (c *Client) listAllDomains() ([]libvirt.Domain, error) {
	domains, _, err := c.v.ConnectListAllDomains(100, 1|2)
	if err != nil {
		return nil, err
	}
	return domains, nil
}

func (c *Client) resolveDomainInterfaces(machine *nodeapi.Machine, schema domain.Schema) {
	// TODO: This is probably not fool-proof, as interfaces can be for more than networking.
	machine.Interfaces = make([]*nodeapi.Interface, len(schema.Devices.Interfaces))

	for i, iface := range schema.Devices.Interfaces {
		machine.Interfaces[i] = new(nodeapi.Interface)

		machine.Interfaces[i].Model = iface.Model.Type
		machine.Interfaces[i].Mac = iface.MAC.MAC

		if iface.Source.Network != "" {
			machine.Interfaces[i].Network = new(string)
			*machine.Interfaces[i].Network = iface.Source.Network
		}
		if iface.Source.Bridge != "" {
			machine.Interfaces[i].Bridge = new(string)
			*machine.Interfaces[i].Bridge = iface.Source.Bridge
		}
	}
}

func (c *Client) resolveDomainDisks(machine *nodeapi.Machine, schema domain.Schema) {
	machine.Disks = make([]*nodeapi.Disk, len(schema.Devices.Disks))

	for i, disk := range schema.Devices.Disks {
		machine.Disks[i] = new(nodeapi.Disk)

		machine.Disks[i].Type = disk.Type
		machine.Disks[i].Device = disk.Device

		// TODO: Support more, (rbd nfs etc)
		if disk.Type == "file" && disk.Source.File != "" {
			volume, err := c.findVolumeByKey(disk.Source.File)
			if err != nil {
				panic(err)
			}
			machine.Disks[i].Volume = fmt.Sprintf("%s;%s", volume.Pool, volume.Key)

		}
	}
}

func (c *Client) completeDomain(dom libvirt.Domain) (*nodeapi.Machine, error) {
	_, _, memory, cores, _, err := c.v.DomainGetInfo(dom)
	if err != nil {
		return nil, err
	}

	xmlSchema, err := c.v.DomainGetXMLDesc(dom, 0)
	if err != nil {
		return nil, err
	}

	var schema domain.Schema
	if err := xml.Unmarshal([]byte(xmlSchema), &schema); err != nil {
		return nil, err
	}

	schemaJson, err := json.Marshal(schema)
	if err != nil {
		return nil, err
	}

	machine := new(nodeapi.Machine)

	topology := new(nodeapi.Topology)

	topology.Cores = uint64(cores)
	topology.Threads = uint64(1)
	topology.Memory = memory

	machine.Uuid = schema.UUID
	machine.Schema = schemaJson
	machine.Name = dom.Name
	machine.Topology = topology

	c.resolveDomainInterfaces(machine, schema)
	c.resolveDomainDisks(machine, schema)

	return machine, nil
}

func (c *Client) QueryMachines() ([]*nodeapi.Machine, error) {
	domains, err := c.listAllDomains()
	if err != nil {
		return nil, err
	}

	machines := make([]*nodeapi.Machine, len(domains))

	for i, domain := range domains {
		machine, err := c.completeDomain(domain)
		if err != nil {
			return nil, err
		}

		machines[i] = machine
	}

	return machines, nil
}
