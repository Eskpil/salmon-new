package tasks

import (
	"context"
	"fmt"

	"github.com/eskpil/salmon/vm/pkg/rockferry"
	"github.com/eskpil/salmon/vm/pkg/rockferry/resource"
	"github.com/eskpil/salmon/vm/pkg/rockferry/status"
)

type SyncStorageVolumesTask struct {
}

func (t *SyncStorageVolumesTask) Execute(ctx context.Context, executor *Executor) error {
	fmt.Println("executing sync storage volumes task")

	volumes, err := executor.Libvirt.QueryStorageVolumes()
	if err != nil {
		return err
	}

	iface := executor.Rockferry.StorageVolumes()
	for _, local := range volumes {
		remotes, err := iface.List(ctx, local.Id, nil)
		if err != nil {
			if status.Is(err, status.ErrNoResults) {
				if err := iface.Create(ctx, local); err != nil {
					return err
				}

				continue
			}

			fmt.Println("failed to find already existing volume", err)

			continue
		}

		if len(remotes) == 0 {
			panic("remotes is zero")
		}

		// NOTE: This will make sure we do not lose any annotations on the way.
		local.Merge(remotes[0])

		if err := iface.Patch(ctx, remotes[0], local); err != nil {
			panic(err)
		}

	}

	return nil
}

type DeleteVolumeTask struct {
	Volume *rockferry.StorageVolume
}

func (t *DeleteVolumeTask) Execute(ctx context.Context, executor *Executor) error {
	return executor.Libvirt.DeleteStorageVolume(t.Volume.Spec.Key)
}

type CreateVolumeTask struct {
	Volume *rockferry.StorageVolume
}

func (t *CreateVolumeTask) Execute(ctx context.Context, executor *Executor) error {
	pools, err := executor.Rockferry.StoragePools().List(ctx, t.Volume.Owner.Id, nil)
	if err != nil {
		return err
	}
	pool := pools[0]

	name := t.Volume.Spec.Name
	format := "raw"
	capacity := t.Volume.Spec.Capacity
	allocation := t.Volume.Spec.Allocation

	if err := executor.Libvirt.CreateVolume(pool.Spec.Name, name, format, capacity, allocation); err != nil {
		return err
	}

	updatedSpec, err := executor.Libvirt.QueryVolumeSpec(pool.Spec.Name, t.Volume.Spec.Name)
	if err != nil {
		return err
	}

	modified := new(rockferry.StorageVolume)
	*modified = *t.Volume
	modified.Spec = *updatedSpec

	return executor.Rockferry.StorageVolumes().Patch(ctx, t.Volume, modified)
}

func (t *CreateVolumeTask) Resource() *resource.Resource[any] {
	return t.Volume.Generic()
}
