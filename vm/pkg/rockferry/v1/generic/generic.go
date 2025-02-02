package genericv1

import (
	"context"

	"github.com/eskpil/salmon/vm/controllerapi"
	"github.com/eskpil/salmon/vm/pkg/convert"
	"github.com/eskpil/salmon/vm/pkg/rockferry/resource"
	"github.com/eskpil/salmon/vm/pkg/rockferry/transport"
	"github.com/snorwin/jsonpatch"
)

type Self = resource.Resource[any]

type Interface struct {
	t *transport.Transport
}

func New(t *transport.Transport) *Interface {
	i := new(Interface)
	i.t = t
	return i
}

func (i *Interface) List(ctx context.Context, kind resource.ResourceKind, owner *resource.OwnerRef) ([]*Self, error) {
	i.t.Lock()
	defer i.t.Unlock()
	api := i.t.C()

	req := new(controllerapi.ListRequest)
	req.Kind = string(kind)
	if owner != nil {
		req.Owner.Id = owner.Id
		req.Owner.Kind = owner.Kind
	}

	response, err := api.List(ctx, req)
	if err != nil {
		return nil, err
	}

	list := make([]*Self, len(response.Resources))

	for i, unmapped := range response.Resources {
		mapped := new(Self)

		mapped.Id = unmapped.Id
		mapped.Kind = resource.ResourceKind(unmapped.Kind)
		mapped.Owner.Id = unmapped.Owner.Id
		mapped.Owner.Kind = unmapped.Owner.Kind
		mapped.Annotations = unmapped.Annotations
		mapped.Status.Phase = resource.Phase(unmapped.Status.Phase)
		mapped.Spec, err = convert.Convert[any](unmapped.Spec)

		list[i] = mapped

	}

	return list, nil
}

func (i *Interface) Watch(ctx context.Context, kind resource.ResourceKind, owner *resource.OwnerRef) (chan *Self, error) {
	i.t.Lock()
	defer i.t.Unlock()
	api := i.t.C()

	req := new(controllerapi.WatchRequest)
	req.Kind = resource.ResourceKindStorageVolume
	if owner != nil {
		req.Owner.Id = owner.Id
		req.Owner.Kind = owner.Kind
	}

	response, err := api.Watch(ctx, req)
	if err != nil {
		return nil, err
	}

	out := make(chan *Self)

	go func() {
		for {
			res, err := response.Recv()
			if err != nil {
				continue
			}

			unmapped := res.Resource
			mapped := new(Self)

			mapped.Id = unmapped.Id
			mapped.Kind = resource.ResourceKind(unmapped.Kind)

			mapped.Owner = new(resource.OwnerRef)
			mapped.Owner.Id = unmapped.Owner.Id
			mapped.Owner.Kind = unmapped.Owner.Kind

			mapped.Annotations = unmapped.Annotations
			mapped.Status.Phase = resource.Phase(unmapped.Status.Phase)
			mapped.Spec, err = convert.Convert[any](unmapped.Spec)

			out <- mapped
		}
	}()

	return out, nil
}

func (i *Interface) Patch(ctx context.Context, original Self, modified Self) error {
	i.t.Lock()
	defer i.t.Unlock()
	api := i.t.C()

	patch, err := jsonpatch.CreateJSONPatch(modified, original)
	if err != nil {
		return err
	}

	req := new(controllerapi.PatchRequest)

	req.Id = new(string)
	*req.Id = original.Id
	req.Kind = string(original.Kind)
	if original.Owner != nil {
		req.Owner = new(controllerapi.Owner)
		req.Owner.Kind = original.Owner.Kind
		req.Owner.Id = original.Owner.Id
	}
	req.Patches = patch.Raw()

	_, err = api.Patch(ctx, req)
	if err != nil {
		return err
	}

	return nil
}
