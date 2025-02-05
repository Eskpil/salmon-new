package tasks

import (
	"context"
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/eskpil/salmon/vm/pkg/rockferry"
	"github.com/eskpil/salmon/vm/pkg/uname"
	"github.com/shirou/gopsutil/v3/cpu"
)

type SyncNodeTask struct{}

func readNodeCpu() (uint64, uint64, uint64, error) {
	info, err := cpu.Info()
	if err != nil {
		return 0, 0, 0, err
	}

	// Determine sockets and threads per core
	coreMap := make(map[string]int)
	threads := len(info)

	for _, cpu := range info {
		coreID := cpu.CoreID
		physicalID := cpu.PhysicalID // Socket ID

		// Count unique physical cores
		coreKey := physicalID + "-" + coreID
		if _, exists := coreMap[coreKey]; !exists {
			coreMap[coreKey] = 1
			threads++
		}
	}

	// Get number of sockets
	socketMap := make(map[string]bool)
	for _, cpu := range info {
		socketMap[cpu.PhysicalID] = true
	}

	return uint64(len(socketMap)), uint64(len(coreMap)), uint64(threads), nil
}

func (t *SyncNodeTask) Execute(ctx context.Context, e *Executor) error {
	fmt.Println("executing sync node task")

	nodes, err := e.Rockferry.Nodes().List(ctx, e.NodeId, nil)
	if err != nil {
		return err
	}

	original := nodes[0]

	modified := new(rockferry.Node)
	*modified = *original

	modified.Spec.Hostname, _ = os.Hostname()

	modified.Spec.ActiveMachines = 2
	modified.Spec.TotalMachines = 10

	sockets, cores, threads, err := readNodeCpu()
	if err != nil {
		return err
	}

	modified.Spec.Topology.Sockets = sockets
	modified.Spec.Topology.Cores = cores
	modified.Spec.Topology.Threads = threads

	var info syscall.Sysinfo_t
	err = syscall.Sysinfo(&info)
	if err != nil {
		return err
	}

	modified.Spec.UpSince = time.Now().Add(-time.Duration(info.Uptime) * time.Second)

	modified.Spec.Topology.Memory = info.Totalram

	uname, _ := uname.New()
	modified.Spec.Kernel = fmt.Sprintf("%s %s %s", uname.Sysname(), uname.Machine(), uname.KernelRelease())

	// TODO: Should be patch, but caused error on controller
	return e.Rockferry.Nodes().Create(ctx, modified)
}
