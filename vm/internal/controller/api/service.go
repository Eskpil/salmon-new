package api

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/eskpil/salmon/vm/controllerapi"
	"github.com/eskpil/salmon/vm/internal/controller/models"
	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/google/uuid"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c Controller) Watch(req *controllerapi.WatchRequest, res grpc.ServerStreamingServer[controllerapi.WatchResponse]) error {
	ctx := res.Context()

	var opts []clientv3.OpOption
	path := ""
	if req.Id == nil {
		opts = append(opts, clientv3.WithPrefix())
		path = fmt.Sprintf("%s/%s/", models.RootKey, req.Kind)
	} else {
		path = fmt.Sprintf("%s/%s/%s", models.RootKey, req.Kind, *req.Id)
	}

	// TODO: Avoid this hack
	if req.Kind == models.ResourceKindStorageVolume && req.Id != nil {
		path = fmt.Sprintf("%s/%s/%s", models.RootKey, req.Kind, req.Owner.Id)
	}

	if req.Action == controllerapi.WatchAction_DELETE {
		opts = append(opts, clientv3.WithPrevKV())
	}

	channel := c.Db.Watch(ctx, path, opts...)

	for {
		w := <-channel
		if w.Canceled {
			break
		}

		if err := w.Err(); err != nil {
			panic(err)
		}

		for _, event := range w.Events {
			// NOTE: Skip unwanted events
			if int(event.Type) != int(req.Action) && req.Action != controllerapi.WatchAction_ALL {
				continue
			}

			if req.Action == controllerapi.WatchAction_DELETE {
				event.Kv = event.PrevKv
			}

			resource := new(controllerapi.Resource)
			if err := json.Unmarshal(event.Kv.Value, resource); err != nil {
				panic(err)
			}

			if req.Owner != nil && req.Owner.Id != "" && req.Owner.Kind != "" && resource.Owner != nil {
				if req.Owner.Id != resource.Owner.Id && req.Owner.Kind != resource.Owner.Kind {
					continue
				}
			}

			response := new(controllerapi.WatchResponse)
			response.Resource = resource
			if err := res.Send(response); err != nil {
				panic(err)
			}

		}
	}

	return nil
}

func (c Controller) List(ctx context.Context, req *controllerapi.ListRequest) (*controllerapi.ListResponse, error) {
	var opts []clientv3.OpOption

	path := ""
	if req.Id == nil {
		opts = append(opts, clientv3.WithPrefix())
		path = fmt.Sprintf("%s/%s/", models.RootKey, req.Kind)
	} else {
		path = fmt.Sprintf("%s/%s/%s", models.RootKey, req.Kind, *req.Id)
	}

	// TODO: Avoid this hack
	if req.Kind == models.ResourceKindStorageVolume && req.Owner != nil && req.Id == nil {
		path = fmt.Sprintf("%s/%s/%s", models.RootKey, req.Kind, req.Owner.Id)
	}

	res, err := c.Db.Get(ctx, path, opts...)
	if err != nil {
		fmt.Println("failed to fetch resources", err)
		return nil, status.Errorf(codes.Internal, "something wrong happend")
	}

	if 0 >= len(res.Kvs) {
		return nil, status.Errorf(codes.NotFound, "resource not found")
	}

	response := new(controllerapi.ListResponse)

	for _, kv := range res.Kvs {
		resource := new(controllerapi.Resource)
		if err := json.Unmarshal(kv.Value, resource); err != nil {
			fmt.Println("unable to unmarshal resource", err)
			return nil, status.Errorf(codes.Internal, "something wrong happend")
		}

		if req.Owner != nil && req.Owner.Id != "" && req.Owner.Kind != "" && resource.Owner != nil {
			if req.Owner.Id != resource.Owner.Id && req.Owner.Kind != resource.Owner.Kind {
				continue
			}
		}

		response.Resources = append(response.Resources, resource)

	}

	return response, nil
}

func (c Controller) Patch(ctx context.Context, req *controllerapi.PatchRequest) (*controllerapi.PatchResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if req.Id == nil {
		return nil, status.Errorf(codes.InvalidArgument, "resource id must be specified")
	}

	path := fmt.Sprintf("%s/%s/%s", models.RootKey, req.Kind, *req.Id)

	res, err := c.Db.Get(ctx, path)
	if err != nil {
		fmt.Println("failed to fetch resource", err)
		return nil, status.Errorf(codes.Internal, "something wrong happend")
	}

	if 0 >= len(res.Kvs) {
		return nil, status.Errorf(codes.NotFound, "resource not found")
	}

	if 2 <= len(res.Kvs) {
		panic("more than 1 response")
	}

	original := res.Kvs[0].Value

	patch, err := jsonpatch.DecodePatch(req.Patches)
	if err != nil {
		panic(err)
	}

	modified, err := patch.Apply(original)
	if err != nil {
		panic(err)
	}

	_, err = c.Db.Put(ctx, path, string(modified))
	if err != nil {
		panic(err)
	}

	return new(controllerapi.PatchResponse), nil
}

func (c Controller) Create(ctx context.Context, input *controllerapi.CreateRequest) (*controllerapi.CreateResponse, error) {
	if input.Resource.Id == "" {
		input.Resource.Id = uuid.NewString()
	}

	path := fmt.Sprintf("%s/%s/%s", models.RootKey, input.Resource.Kind, input.Resource.Id)

	bytes, err := json.Marshal(input.Resource)
	if err != nil {
		panic(err)
	}

	if _, err := c.Db.Put(ctx, path, string(bytes)); err != nil {
		fmt.Println("failed to insert resource", err)
		return nil, status.Errorf(codes.Internal, "something wrong happend")
	}

	return new(controllerapi.CreateResponse), nil
}

func (c Controller) Delete(ctx context.Context, req *controllerapi.DeleteRequest) (*controllerapi.DeleteResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	path := fmt.Sprintf("%s/%s/%s", models.RootKey, req.Kind, req.Id)

	_, err := c.Db.Delete(ctx, path)
	if err != nil {
		fmt.Println("failed to delete resource")
		return nil, status.Errorf(codes.Internal, "something wrong happend")
	}

	return new(controllerapi.DeleteResponse), nil
}
