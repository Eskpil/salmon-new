package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os"
	"runtime"

	"github.com/eskpil/salmon/vm/internal/node/queries"
	"github.com/eskpil/salmon/vm/nodeapi"
	"github.com/eskpil/salmon/vm/pkg/uname"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Node struct {
	qc *queries.Client
	nodeapi.UnimplementedNodeApiServer
}

func (n Node) ListMachines(context.Context, *nodeapi.ListMachinesRequest) (*nodeapi.MachineList, error) {
	machines, err := n.qc.QueryMachines()
	if err != nil {
		return nil, err
	}

	list := new(nodeapi.MachineList)
	list.List = machines
	return list, nil
}

func (n Node) ListNetworks(context.Context, *nodeapi.ListNetworksRequest) (*nodeapi.NetworkList, error) {
	networks, err := n.qc.ListAllNetworks()
	if err != nil {
		return nil, err
	}

	list := new(nodeapi.NetworkList)
	list.List = networks
	return list, err
}

func (n Node) ListStoragePools(context.Context, *nodeapi.ListStoragePoolsRequest) (*nodeapi.StoragePoolList, error) {
	pools, err := n.qc.QueryStoragePools()
	if err != nil {
		return nil, err
	}

	list := new(nodeapi.StoragePoolList)
	list.List = pools
	return list, nil
}

func (n Node) ListStorageVolumes(context.Context, *nodeapi.ListStorageVolumesRequest) (*nodeapi.StorageVolumeList, error) {
	volumes, err := n.qc.QueryStorageVolumes()
	if err != nil {
		return nil, err
	}

	list := new(nodeapi.StorageVolumeList)
	list.List = volumes
	return list, nil
}

func (n Node) CreateVolume(ctx context.Context, req *nodeapi.CreateVolumeRequest) (*nodeapi.CreateVolumeResponse, error) {
	if req.Allocation <= 0 {
		return nil, errors.New("allocation must be greater than 0")
	}

	err := n.qc.CreateVolume(req.Pool, req.Name, req.Format, int(req.Allocation))
	if err != nil {
		return nil, err
	}

	res := new(nodeapi.CreateVolumeResponse)
	return res, nil
}

func (n Node) Ping(ctx context.Context, req *nodeapi.PingRequest) (*nodeapi.PingResponse, error) {
	res := new(nodeapi.PingResponse)
	node := new(nodeapi.Node)

	// TODO: Use live data
	node.ActiveMachines = 2
	node.TotalMachines = 10

	node.Topology = new(nodeapi.Topology)

	// TODO: Do not do this, figure out memory
	node.Topology.Cores = uint64(runtime.NumCPU()) / 2
	node.Topology.Threads = 2

	node.Hostname, _ = os.Hostname()

	uname, _ := uname.New()
	node.Kernel = fmt.Sprintf("%s %s %s", uname.Sysname(), uname.Machine(), uname.KernelRelease())

	res.Node = node

	return res, nil
}

func old_main() {
	node := Node{}

	client, err := queries.NewClient()
	if err != nil {
		fmt.Println("failed to initialize client", err)
		return
	}

	node.qc = client

	listener, err := net.Listen("tcp", "0.0.0.0:9090")
	if err != nil {
		slog.Error("could not create listener", slog.Any("err", err))
		return
	}

	server := grpc.NewServer()
	nodeapi.RegisterNodeApiServer(server, node)

	reflection.Register(server)

	fmt.Println("starting up")

	if err := server.Serve(listener); err != nil {
		slog.Error("could not serve requests", slog.Any("err", err))
	}
}
