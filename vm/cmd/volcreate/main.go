package main

import (
	"context"

	"github.com/eskpil/salmon/vm/nodeapi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cc, err := grpc.NewClient("10.100.0.101:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	client := nodeapi.NewNodeApiClient(cc)

	ctx := context.Background()

	req := new(nodeapi.CreateVolumeRequest)

	req.Pool = "6ddb6928-dc10-44fc-b7a3-e4632b2eef76"
	req.Name = "test-rdb123"
	req.Format = "raw"
	req.Allocation = 8 * 1024 * 1024

	_, err = client.CreateVolume(ctx, req)
	if err != nil {
		panic(err)
	}

}
