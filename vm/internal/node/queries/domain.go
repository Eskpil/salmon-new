package queries

import (
	"encoding/json"
	"encoding/xml"
	"fmt"

	"github.com/digitalocean/go-libvirt"
	"github.com/eskpil/salmon/vm/nodeapi"
	machinesv1 "github.com/eskpil/salmon/vm/pkg/rockferry/v1/machines"
	"github.com/eskpil/salmon/vm/pkg/virtwrap/domain"
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

func (c *Client) CreateDomain(spec *machinesv1.Spec) error {
	schema := new(domain.Schema)

	schema.Name = spec.Name
	schema.Type = "kvm"

	schema.Memory.Unit = "bytes"
	schema.Memory.Value = spec.Topology.Memory

	schema.VCPU = new(domain.VCPU)
	schema.VCPU.CPUs = uint32(spec.Topology.Cores) * uint32(spec.Topology.Threads)
	schema.VCPU.Placement = "static"

	schema.CPU.Topology = new(domain.CPUTopology)
	schema.CPU.Topology.Cores = uint32(spec.Topology.Cores)
	schema.CPU.Topology.Threads = uint32(spec.Topology.Threads)
	schema.CPU.Topology.Sockets = 1
	schema.CPU.Mode = "host-passthrough"

	schema.Features = new(domain.Features)
	schema.Features.ACPI = new(domain.FeatureEnabled)
	schema.Features.APIC = new(domain.FeatureEnabled)

	schema.Devices.Emulator = "/usr/bin/qemu-system-x86_64"

	schema.OS.Type.Arch = "x86_64"
	schema.OS.Type.Machine = "pc-q35-7.2"
	schema.OS.Type.OS = "hvm"

	schema.OS.BootOrder = append(schema.OS.BootOrder, domain.Boot{Dev: "hd"})
	schema.OS.BootOrder = append(schema.OS.BootOrder, domain.Boot{Dev: "cdrom"})

	for _, d := range spec.Disks {
		disk := new(domain.Disk)

		if d.Type == "network" {
			disk.Type = "network"
			disk.Device = d.Device

			disk.Driver = new(domain.DiskDriver)
			disk.Driver.Name = "qemu"
			disk.Driver.Type = "raw"

			disk.Auth = new(domain.DiskAuth)

			disk.Auth.Username = d.Network.Auth.Username
			disk.Auth.Secret = new(domain.DiskSecret)
			disk.Auth.Secret.Type = d.Network.Auth.Type
			disk.Auth.Secret.UUID = d.Network.Auth.Secret

			disk.Source.Protocol = d.Network.Protocol
			disk.Source.Name = d.Network.Key
			disk.Source.Host = new(domain.DiskSourceHost)
			disk.Source.Host.Name = d.Network.Hosts[0].Name
			disk.Source.Host.Port = d.Network.Hosts[0].Port

			disk.Target.Bus = "virtio"
			// TODO: Create a function which returns unique device names
			disk.Target.Device = "vda"
		}

		if d.Type == "file" {
			disk.Type = "file"
			disk.Device = d.Device

			disk.Source.File = d.File.Key

			disk.Driver = new(domain.DiskDriver)
			disk.Driver.Name = "qemu"
			disk.Driver.Type = "raw"

			disk.Target.Bus = "sata"
			disk.Target.Device = "sda"

		}

		schema.Devices.Disks = append(schema.Devices.Disks, *disk)
	}

	for _, i := range spec.Interfaces {
		iface := new(domain.Interface)

		iface.MAC = new(domain.MAC)
		iface.MAC.MAC = i.Mac
		iface.Type = "network"
		iface.Source.Network = "bridged-network"
		iface.Model = new(domain.Model)
		iface.Model.Type = "virtio"

		schema.Devices.Interfaces = append(schema.Devices.Interfaces, *iface)
	}

	vnc := new(domain.Graphics)

	vnc.Type = "vnc"
	vnc.AutoPort = "yes"
	vnc.Passwd.Value = "123"
	vnc.Listen = new(domain.GraphicsListen)

	vnc.Listen.Type = "address"
	vnc.Listen.Address = "0.0.0.0"

	schema.Devices.Graphics = append(schema.Devices.Graphics, *vnc)

	bytes, err := xml.Marshal(schema)
	if err != nil {
		panic(err)
	}

	returned, err := c.v.DomainCreateXML(string(bytes), 0)
	if err != nil {
		panic(err)
	}

	_ = returned

	return nil
}
