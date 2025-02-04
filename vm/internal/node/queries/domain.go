package queries

import (
	"encoding/xml"

	"github.com/digitalocean/go-libvirt"
	"github.com/eskpil/salmon/vm/pkg/rockferry/spec"
	"github.com/eskpil/salmon/vm/pkg/virtwrap/domain"
	"github.com/google/uuid"
)

func (c *Client) CreateDomain(id string, spec *spec.MachineSpec) error {
	schema := new(domain.Schema)

	schema.Name = spec.Name
	schema.Type = "kvm"

	schema.UUID = id

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
			disk.Source.Name = d.Key
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

			disk.Source.File = d.Key

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

func (c *Client) DestroyDomain(id string) error {
	domId := uuid.MustParse(id)

	domain, err := c.v.DomainLookupByUUID(libvirt.UUID(domId))
	if err != nil {
		return err
	}

	return c.v.DomainDestroy(domain)
}
