package tasks

import (
	"context"

	"github.com/eskpil/salmon/vm/pkg/rockferry/resource"
	storagevolumesv1 "github.com/eskpil/salmon/vm/pkg/rockferry/v1/storagevolumes"
)

type RefreshVolumeTask struct {
	Volume *storagevolumesv1.Self
}

type CreateVolumeTask struct {
	Volume *storagevolumesv1.Self
}

func (t *CreateVolumeTask) Execute(ctx context.Context, executor *Executor) error {
	pool := t.Volume.Owner.Id
	name := t.Volume.Spec.Name
	format := "raw"
	capacity := t.Volume.Spec.Capacity
	allocation := t.Volume.Spec.Allocation

	if err := executor.Libvirt.CreateVolume(pool, name, format, capacity, allocation); err != nil {
		return err
	}

	updatedSpec, err := executor.Libvirt.QueryVolumeSpec(t.Volume.Owner.Id, t.Volume.Spec.Name)
	if err != nil {
		return err
	}

	modified := new(storagevolumesv1.Self)
	*modified = *t.Volume
	modified.Spec = updatedSpec

	return executor.Rockferry.StorageVolumes().Patch(ctx, t.Volume, modified)
}

func (t *CreateVolumeTask) Resource() *resource.Resource[any] {
	return t.Volume.Generic()
}

func (t *RefreshVolumeTask) Execute(ctx context.Context, executor *Executor) error {

	return nil

}
