package genericv1

import (
	"context"

	"github.com/eskpil/salmon/vm/pkg/rockferry/resource"
	"github.com/eskpil/salmon/vm/pkg/rockferry/transport"
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

func (i *Interface) List(ctx context.Context, kind resource.ResourceKind, id string, owner *resource.OwnerRef) ([]*Self, error) {
	return i.t.List(ctx, kind, id, owner)
}

func (i *Interface) Watch(ctx context.Context, kind resource.ResourceKind, id string, owner *resource.OwnerRef) (chan *Self, error) {
	return i.t.Watch(ctx, kind, id, owner)
}

func (i *Interface) Patch(ctx context.Context, original *Self, modified *Self) error {
	return i.t.Patch(ctx, original, modified)
}
