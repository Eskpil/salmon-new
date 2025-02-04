package transport

import (
	"context"
	"fmt"
	"time"

	"github.com/eskpil/salmon/vm/controllerapi"
	"github.com/eskpil/salmon/vm/pkg/convert"
	"github.com/eskpil/salmon/vm/pkg/rockferry/resource"
	rstatus "github.com/eskpil/salmon/vm/pkg/rockferry/status"
	"github.com/snorwin/jsonpatch"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type Transport struct {
	client controllerapi.ControllerApiClient
}

func New(url string) (*Transport, error) {
	t := new(Transport)

	cc, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	t.client = controllerapi.NewControllerApiClient(cc)

	return t, nil
}

func (t *Transport) C() controllerapi.ControllerApiClient {
	return t.client
}

func (t *Transport) Watch(ctx context.Context, action int, kind resource.ResourceKind, id string, owner *resource.OwnerRef) (chan *resource.Resource[interface{}], error) {
	api := t.C()

	// Create the initial watch request
	req := new(controllerapi.WatchRequest)
	req.Kind = kind
	if owner != nil {
		req.Owner.Id = owner.Id
		req.Owner.Kind = owner.Kind
	}

	req.Action = controllerapi.WatchAction(action)

	// Create the channel to send updates
	out := make(chan *resource.Resource[any])

	// Function to start watching and handle reconnection
	var watch func() error
	watch = func() error {
		// Call the API to start watching
		response, err := api.Watch(ctx, req)
		if err != nil {
			return err
		}

		// Process incoming resources and handle reconnection
		go func() {
			for {
				select {
				default:
					res, err := response.Recv()
					if err != nil {
						// Log the error (can be enhanced with more detailed logging)
						fmt.Printf("Watch receive error: %v. Reconnecting...\n", err)

						// Retry after a delay (backoff logic can be implemented here)
						// You may add exponential backoff, but here we'll retry immediately for simplicity
						retryDelay := 2 * time.Second // You can increase this over multiple retries
						select {
						case <-time.After(retryDelay): // Retry after waiting for the delay
							// Attempt to restart the watch connection
							if err := watch(); err != nil {
								fmt.Printf("Failed to reconnect watch: %v\n", err)
								return
							}
						case <-ctx.Done():
							close(out)
							return
						}
						return
					}

					// Process the resource
					unmapped := res.Resource
					mapped := new(resource.Resource[any])

					// Map the resource
					mapped.Id = unmapped.Id
					mapped.Kind = resource.ResourceKind(unmapped.Kind)

					mapped.Owner = new(resource.OwnerRef)
					mapped.Owner.Id = unmapped.Owner.Id
					mapped.Owner.Kind = unmapped.Owner.Kind

					mapped.Annotations = unmapped.Annotations
					mapped.Status.Phase = resource.Phase(unmapped.Status.Phase)

					mapped.RawSpec = unmapped.Spec

					// Send the mapped resource to the channel
					select {
					case out <- mapped:
					case <-ctx.Done():
						close(out)
						return
					}
				}
			}
		}()
		return nil
	}

	// Start the watch and handle the error if it fails
	if err := watch(); err != nil {
		return nil, err
	}

	return out, nil
}

func (t *Transport) Patch(ctx context.Context, original *resource.Resource[any], modified *resource.Resource[any]) error {
	api := t.C()

	patch, err := jsonpatch.CreateJSONPatch(modified, original)
	if err != nil {
		return err
	}

	// NOTE: If the resource is not update it, why bother the controller
	if 0 >= len(patch.Raw()) {
		return nil
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

func (t *Transport) List(ctx context.Context, kind resource.ResourceKind, id string, owner *resource.OwnerRef) ([]*resource.Resource[any], error) {
	api := t.C()

	req := new(controllerapi.ListRequest)
	if id != "" {
		req.Id = new(string)
		*req.Id = id
	}
	req.Kind = string(kind)
	if owner != nil {
		req.Owner.Id = owner.Id
		req.Owner.Kind = owner.Kind
	}

	response, err := api.List(ctx, req)
	if err != nil {
		if s, ok := status.FromError(err); ok && s != nil {
			if s.Code() == codes.NotFound {
				return nil, rstatus.NewError(rstatus.ErrNoResults, "no results")
			}
		}

		return nil, err
	}

	list := make([]*resource.Resource[any], len(response.Resources))

	for i, unmapped := range response.Resources {
		mapped := new(resource.Resource[any])

		mapped.Id = unmapped.Id
		mapped.Kind = resource.ResourceKind(unmapped.Kind)
		if unmapped.Owner != nil {
			mapped.Owner = new(resource.OwnerRef)
			mapped.Owner.Id = unmapped.Owner.Id
			mapped.Owner.Kind = unmapped.Owner.Kind
		}
		mapped.Annotations = unmapped.Annotations
		mapped.Status.Phase = resource.Phase(unmapped.Status.Phase)
		mapped.RawSpec = unmapped.GetSpec()

		list[i] = mapped

	}

	return list, nil
}

func (t *Transport) Create(ctx context.Context, in *resource.Resource[any]) error {
	api := t.C()

	req := new(controllerapi.CreateRequest)
	req.Resource = new(controllerapi.Resource)

	req.Resource.Id = in.Id

	req.Resource.Kind = string(in.Kind)
	req.Resource.Annotations = in.Annotations

	if in.Owner != nil {
		req.Resource.Owner = new(controllerapi.Owner)
		req.Resource.Owner.Id = in.Owner.Id
		req.Resource.Owner.Kind = in.Owner.Kind
	}

	spec, err := convert.Outgoing(&in.Spec)
	if err != nil {
		return err
	}

	req.Resource.Spec = spec
	req.Resource.Status = new(controllerapi.Status)
	req.Resource.Status.Phase = string(in.Status.Phase)

	_, err = api.Create(ctx, req)
	return err
}

func (t *Transport) Delete(ctx context.Context, kind resource.ResourceKind, id string) error {
	api := t.C()

	req := new(controllerapi.DeleteRequest)
	req.Kind = kind
	req.Id = id

	_, err := api.Delete(ctx, req)
	return err
}
