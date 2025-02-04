package tasks

import (
	"context"
	"fmt"

	"github.com/eskpil/salmon/vm/pkg/mac"
	"github.com/eskpil/salmon/vm/pkg/rockferry"
	"github.com/eskpil/salmon/vm/pkg/rockferry/resource"
	"github.com/eskpil/salmon/vm/pkg/rockferry/spec"
	"github.com/google/uuid"
)

type CreateVirtualMachineTask struct {
	Request *rockferry.MachineRequest
}

func (t *CreateVirtualMachineTask) createVmDisks(ctx context.Context, executor *Executor, vmId string) ([]*spec.MachineSpecDisk, error) {
	disks := []*spec.MachineSpecDisk{}

	for _, disk := range t.Request.Spec.Disks {
		poolId := disk.Pool
		pools, err := executor.Rockferry.StoragePools().List(ctx, poolId, nil)
		if err != nil {
			return nil, err
		}

		pool := pools[0]

		name := uuid.NewString()
		format := "raw"

		capacity := disk.Capacity
		allocation := disk.Capacity

		// TODO: Check if volume already is created and continue
		if err := executor.Libvirt.CreateVolume(pool.Spec.Name, name, format, capacity, allocation); err != nil {
			return nil, err
		}

		volumeSpec, err := executor.Libvirt.QueryVolumeSpec(pool.Spec.Name, name)
		if err != nil {
			return nil, err
		}

		out := new(rockferry.StorageVolume)

		out.Id = fmt.Sprintf("%s/%s", pool.Id, name)
		out.Kind = resource.ResourceKindStorageVolume

		out.Annotations = map[string]string{}
		out.Annotations["vm"] = vmId
		out.Annotations["vm.name"] = t.Request.Spec.Name

		out.Spec = *volumeSpec
		out.Status.Phase = resource.PhaseCreated
		out.Owner = new(resource.OwnerRef)
		out.Owner.Id = pool.Id
		out.Owner.Kind = resource.ResourceKindStoragePool

		if err := executor.Rockferry.StorageVolumes().Create(ctx, out); err != nil {
			panic(err)
		}

		d := new(spec.MachineSpecDisk)
		d.Key = volumeSpec.Key
		d.Volume = out.Id
		if pool.Spec.Type == "rbd" {
			d.Type = "network"
			d.Device = "disk"

			d.Network = new(spec.MachineSpecDiskNetwork)
			d.Network.Protocol = pool.Spec.Type
			d.Network.Hosts = pool.Spec.Source.Hosts
			d.Network.Auth = *pool.Spec.Source.Auth
		}

		if pool.Spec.Type == "dir" {
			d.Type = "file"
			d.Device = "disk"

			d.File = new(spec.MachineSpecDiskFile)
		}

		disks = append(disks, d)
	}

	// TODO: CDROM can be network disk as well
	cdrom := new(spec.MachineSpecDisk)

	cdrom.Key = t.Request.Spec.Cdrom.Key
	// This could probably be more clean
	cdrom.File = new(spec.MachineSpecDiskFile)
	cdrom.Device = "cdrom"
	cdrom.Type = "file"

	disks = append(disks, cdrom)

	return disks, nil
}

func (t *CreateVirtualMachineTask) createNetworkInterfaces(ctx context.Context, executor *Executor) ([]*spec.MachineSpecInterface, error) {
	interfaces := make([]*spec.MachineSpecInterface, 1)

	networks, err := executor.Rockferry.Networks().List(ctx, t.Request.Spec.Network, nil)
	if err != nil {
		return nil, err
	}

	network := networks[0]

	mac, err := mac.Generate()
	if err != nil {
		return nil, err
	}

	interfaces[0] = new(spec.MachineSpecInterface)
	interfaces[0].Mac = mac
	interfaces[0].Model = "virtio"

	interfaces[0].Network = new(string)
	*interfaces[0].Network = network.Spec.Name

	interfaces[0].Bridge = new(string)
	*interfaces[0].Bridge = network.Spec.Bridge.Name

	return interfaces, nil
}

func (t *CreateVirtualMachineTask) Execute(ctx context.Context, executor *Executor) error {
	// NOTE: Used to annotate storage volumes with the vm id. This is useful for deletion.
	vmId := uuid.NewString()

	disks, err := t.createVmDisks(ctx, executor, vmId)
	if err != nil {
		return err
	}

	interfaces, err := t.createNetworkInterfaces(ctx, executor)
	if err != nil {
		return err
	}

	spec := new(spec.MachineSpec)

	spec.Name = t.Request.Spec.Name
	spec.Topology = t.Request.Spec.Topology
	spec.Disks = disks
	spec.Interfaces = interfaces

	res := new(rockferry.Machine)

	res.Annotations = map[string]string{}
	res.Annotations["machinerequest.id"] = t.Request.Id

	res.Id = vmId
	res.Kind = resource.ResourceKindMachine
	res.Owner = new(resource.OwnerRef)
	// TODO: Do not hardcode this
	res.Owner.Id = executor.NodeId
	res.Owner.Kind = resource.ResourceKindNode

	res.Status.Phase = resource.PhaseCreated

	if err := executor.Libvirt.CreateDomain(vmId, spec); err != nil {
		return err
	}

	res.Spec = *spec
	res.Status.Phase = resource.PhaseCreated

	return executor.Rockferry.Machines().Create(ctx, res)
}

func (t *CreateVirtualMachineTask) Resource() *resource.Resource[any] {
	return t.Request.Generic()
}

type DeleteVmTask struct {
	Machine *rockferry.Machine
}

func (t *DeleteVmTask) Execute(ctx context.Context, e *Executor) error {
	if err := e.Libvirt.DestroyDomain(t.Machine.Id); err != nil {
		return err
	}

	// Cleanup, yay
	for _, disk := range t.Machine.Spec.Disks {
		if disk.Volume == "" {
			continue
		}

		if err := e.Rockferry.StorageVolumes().Delete(ctx, disk.Volume); err != nil {
			fmt.Println("failed to delete storage volume", err)
			continue
		}
	}

	return e.Rockferry.MachineRequests().Delete(ctx, t.Machine.Annotations["machinerequest.id"])
}
