package cron

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/eskpil/salmon/vm/internal/controller/models"
	"github.com/eskpil/salmon/vm/nodeapi"
	"github.com/google/uuid"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func SyncWithNodes() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// TODO: Enable some kind of config
	// TODO: Avoid multiple db connections
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		panic(err)
	}

	path := fmt.Sprintf("%s/%s", models.RootKey, models.ResourceKindNode)
	res, err := cli.Get(ctx, path, clientv3.WithPrefix())
	if err != nil {
		fmt.Println("failed to list nodes")
		return
	}

	for _, kv := range res.Kvs {
		resource := new(models.Resource)
		if err := json.Unmarshal(kv.Value, resource); err != nil {
			fmt.Println("failed to unmarshal node")
			continue
		}

		if err := syncNode(ctx, cli, resource); err != nil {
			fmt.Println("failed to resync node")
			continue
		}

	}

	return
}

func syncNode(ctx context.Context, db *clientv3.Client, resource *models.Resource) error {
	cc, err := grpc.NewClient(resource.Annotations["node.url"], grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	api := nodeapi.NewNodeApiClient(cc)

	res, err := api.Ping(ctx, new(nodeapi.PingRequest))
	if err != nil {
		fmt.Println("failed to ping node")
		return err
	}

	resource.Spec = res.Node
	path := fmt.Sprintf("%s/%s/%s", models.RootKey, models.ResourceKindNode, resource.Id)
	if _, err := db.Put(ctx, path, string(resource.Marshal())); err != nil {
		fmt.Println("failed to update node spec")
		return err
	}

	if err := syncNodeNetworks(ctx, api, db, resource); err != nil {
		fmt.Println("failed to sync node networks")
		return err
	}

	if err := syncNodeStorage(ctx, api, db, resource); err != nil {
		fmt.Println("failed to sync node storage")
		return err
	}

	if err := syncNodeMachines(ctx, api, db, resource); err != nil {
		fmt.Println("failed to sync node machines")
		return err
	}

	return nil
}

func syncNodeNetworks(ctx context.Context, api nodeapi.NodeApiClient, db *clientv3.Client, node *models.Resource) error {
	networks, err := api.ListNetworks(ctx, new(nodeapi.ListNetworksRequest))
	if err != nil {
		return err
	}

	for _, network := range networks.List {
		path := fmt.Sprintf("%s/%s/%s", models.RootKey, models.ResourceKindNetwork, network.Uuid)
		_, err := db.Get(ctx, path)
		if err != nil {
			fmt.Println("failed to check if network exists, continuing")
			continue
		}

		resource := new(models.Resource)

		resource.Kind = models.ResourceKindNetwork
		resource.Id = network.Uuid
		resource.Owner = new(models.OwnerRef)
		resource.Owner.Id = node.Id
		resource.Owner.Kind = models.ResourceKindNode
		resource.Spec = network

		// If the network already exists it will be updated and revision count updated
		bytes, err := json.Marshal(resource)
		if err != nil {
			fmt.Println("failed to marshal network to json, continuing")
			continue
		}
		_, err = db.Put(ctx, path, string(bytes))
		if err != nil {
			fmt.Println("failed to put network to database, contiunuing")
		}
	}

	return nil
}

func syncNodeStorage(ctx context.Context, api nodeapi.NodeApiClient, db *clientv3.Client, node *models.Resource) error {
	pools, err := api.ListStoragePools(ctx, new(nodeapi.ListStoragePoolsRequest))
	if err != nil {
		return err
	}

	for _, pool := range pools.List {
		path := fmt.Sprintf("%s/%s/%s", models.RootKey, models.ResourceKindStoragePool, pool.Uuid)
		_, err := db.Get(ctx, path)
		if err != nil {
			fmt.Println("failed to check if pool exists, continuing")
			continue
		}

		resource := new(models.Resource)

		resource.Kind = models.ResourceKindNetwork
		resource.Id = uuid.NewString()
		resource.Owner = new(models.OwnerRef)
		resource.Owner.Id = node.Id
		resource.Owner.Kind = models.ResourceKindNode
		resource.Spec = pool

		// If the storage pool already exists it will be updated and revision count updated
		bytes, err := json.Marshal(resource)
		if err != nil {
			fmt.Println("failed to marshal pool to json, continuing")
			continue
		}
		_, err = db.Put(ctx, path, string(bytes))
		if err != nil {
			fmt.Println("failed to put pool to database, contiunuing")
		}
	}

	volumes, err := api.ListStorageVolumes(ctx, new(nodeapi.ListStorageVolumesRequest))
	if err != nil {
		return err
	}

	for _, volume := range volumes.List {
		path := fmt.Sprintf("%s/%s/%s/%s", models.RootKey, models.ResourceKindStorageVolume, volume.Pool, volume.Name)
		_, err := db.Get(ctx, path)
		if err != nil {
			fmt.Println("failed to check if volume exists, continuing")
			continue
		}

		resource := new(models.Resource)

		resource.Kind = models.ResourceKindNetwork
		resource.Id = uuid.NewString()
		resource.Owner = new(models.OwnerRef)
		resource.Owner.Id = volume.Pool
		resource.Owner.Kind = models.ResourceKindStoragePool
		resource.Spec = volume

		// If the storage volume already exists it will be updated and revision count updated
		bytes, err := json.Marshal(resource)
		if err != nil {
			fmt.Println("failed to marshal volume to json, continuing")
			continue
		}
		_, err = db.Put(ctx, path, string(bytes))
		if err != nil {
			fmt.Println("failed to put volume to database, contiunuing")
		}
	}

	return nil
}

func syncNodeMachines(ctx context.Context, api nodeapi.NodeApiClient, db *clientv3.Client, node *models.Resource) error {
	machines, err := api.ListMachines(ctx, new(nodeapi.ListMachinesRequest))
	if err != nil {
		return err
	}

	for _, machine := range machines.List {
		path := fmt.Sprintf("%s/%s/%s", models.RootKey, models.ResourceKindMachine, machine.Uuid)
		_, err := db.Get(ctx, path)
		if err != nil {
			fmt.Println("failed to check if machine exists, continuing")
			continue
		}

		resource := new(models.Resource)

		resource.Kind = models.ResourceKindNetwork
		resource.Id = uuid.NewString()
		resource.Owner = new(models.OwnerRef)
		resource.Owner.Id = node.Id
		resource.Owner.Kind = models.ResourceKindNode
		resource.Spec = machine

		// If the storage volume already exists it will be updated and revision count updated
		bytes, err := json.Marshal(resource)
		if err != nil {
			fmt.Println("failed to marshal volume to json, continuing")
			continue
		}
		_, err = db.Put(ctx, path, string(bytes))
		if err != nil {
			fmt.Println("failed to put volume to database, contiunuing")
		}
	}

	return nil
}
