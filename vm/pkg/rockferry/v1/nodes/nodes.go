package nodesv1

import (
	"context"

	"github.com/eskpil/salmon/vm/pkg/convert"
	"github.com/eskpil/salmon/vm/pkg/rockferry/resource"
	"github.com/eskpil/salmon/vm/pkg/rockferry/transport"
)

type Spec struct {
}

type Self = resource.Resource[*Spec]

type Interface struct {
	t *transport.Transport
}

func fix(unmapped *resource.Resource[any]) *Self {
	mapped := new(Self)
	mapped.Id = unmapped.Id
	mapped.Owner = unmapped.Owner
	mapped.Kind = unmapped.Kind
	mapped.Annotations = unmapped.Annotations
	mapped.Status = unmapped.Status
	mapped.Spec, _ = convert.Convert[Spec](unmapped.RawSpec)

	return mapped
}

func New(t *transport.Transport) *Interface {
	i := new(Interface)
	i.t = t
	return i
}

func (i *Interface) List(ctx context.Context, id string, owner *resource.OwnerRef) ([]*Self, error) {
	in, err := i.t.List(ctx, resource.ResourceKindStorageVolume, id, owner)
	if err != nil {
		return nil, err
	}

	out := make([]*Self, len(in))

	for i, unmapped := range in {
		out[i] = fix(unmapped)
	}

	return out, nil
}

func (i *Interface) Watch(ctx context.Context, id string, owner *resource.OwnerRef) (chan *Self, error) {
	in, err := i.t.Watch(ctx, resource.ResourceKindStorageVolume, id, owner)
	if err != nil {
		return nil, err
	}

	out := make(chan *Self)

	go func() {
		for {
			out <- fix(<-in)
		}
	}()

	return out, nil
}

func (i *Interface) Patch(ctx context.Context, original *Self, modified *Self) error {
	return i.t.Patch(ctx, original.Generic(), modified.Generic())
}
