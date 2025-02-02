package tasks

import (
	"context"
	"fmt"

	"github.com/eskpil/salmon/vm/pkg/mac"
	"github.com/eskpil/salmon/vm/pkg/rockferry/resource"
	machinerequestsv1 "github.com/eskpil/salmon/vm/pkg/rockferry/v1/machinerequests"
	machinesv1 "github.com/eskpil/salmon/vm/pkg/rockferry/v1/machines"
	storagevolumesv1 "github.com/eskpil/salmon/vm/pkg/rockferry/v1/storagevolumes"
	"github.com/google/uuid"
)

type CreateVirtualMachineTask struct {
	Request *machinerequestsv1.Self
}

func (t *CreateVirtualMachineTask) createVmDisks(ctx context.Context, executor *Executor) ([]*machinesv1.SpecDisk, error) {
	fmt.Println(t.Request.Spec)

	disks := make([]*machinesv1.SpecDisk, len(t.Request.Spec.Disks))

	for i, disk := range t.Request.Spec.Disks {
		pool := disk.Pool
		name := uuid.NewString()
		format := "raw"

		capacity := disk.Capacity
		allocation := disk.Capacity

		// TODO: Check if volume already is created and continue
		if err := executor.Libvirt.CreateVolume(pool, name, format, capacity, allocation); err != nil {
			return nil, err
		}

		spec, err := executor.Libvirt.QueryVolumeSpec(disk.Pool, name)
		if err != nil {
			return nil, err
		}

		out := new(storagevolumesv1.Self)

		out.Kind = resource.ResourceKindStorageVolume
		out.Spec = spec
		out.Status.Phase = resource.PhaseCreated
		out.Owner = new(resource.OwnerRef)
		out.Owner.Id = pool
		out.Owner.Kind = resource.ResourceKindStoragePool

		if err := executor.Rockferry.StorageVolumes().Create(ctx, out); err != nil {
			panic(err)
		}

		disks[i] = new(machinesv1.SpecDisk)
		disks[i].Volume = spec.Key
		disks[i].Device = "disk"
	}

	// This could probably be more clean
	disks[len(disks)-1] = new(machinesv1.SpecDisk)
	disks[len(disks)-1].Volume = t.Request.Spec.Cdrom.Key
	disks[len(disks)-1].Device = "cdrom"
	disks[len(disks)-1].Type = "file"

	return disks, nil
}

func (t *CreateVirtualMachineTask) createNetworkInterfaces(ctx context.Context, executor *Executor) ([]*machinesv1.SpecInterface, error) {
	interfaces := make([]*machinesv1.SpecInterface, 1)

	networks, err := executor.Rockferry.Networks().List(ctx, t.Request.Spec.Network, nil)
	if err != nil {
		return nil, err
	}

	network := networks[0]

	mac, err := mac.Generate()
	if err != nil {
		return nil, err
	}

	interfaces[0] = new(machinesv1.SpecInterface)
	interfaces[0].Mac = mac
	interfaces[0].Model = "virtio"

	interfaces[0].Network = new(string)
	*interfaces[0].Network = network.Spec.Name

	interfaces[0].Bridge = new(string)
	*interfaces[0].Bridge = network.Spec.Bridge.Name

	return interfaces, nil
}

func (t *CreateVirtualMachineTask) Execute(ctx context.Context, executor *Executor) error {
	disks, err := t.createVmDisks(ctx, executor)
	if err != nil {
		return err
	}

	interfaces, err := t.createNetworkInterfaces(ctx, executor)
	if err != nil {
		return err
	}

	spec := new(machinesv1.Spec)

	spec.Name = t.Request.Spec.Name
	spec.Topology = t.Request.Spec.Topology
	spec.Disks = disks
	spec.Interfaces = interfaces

	return executor.Libvirt.CreateDomain(spec)
}

func (t *CreateVirtualMachineTask) Resource() *resource.Resource[any] {
	return t.Request.Generic()
}
