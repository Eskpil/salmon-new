package tasks

import (
	"context"
	"fmt"
	"os"
	"runtime"

	"github.com/eskpil/salmon/vm/pkg/rockferry"
	"github.com/eskpil/salmon/vm/pkg/uname"
)

type SyncNodeTask struct{}

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

	// TODO: Do not do this and figure out memory
	modified.Spec.Topology.Cores = uint64(runtime.NumCPU()) / 2
	modified.Spec.Topology.Threads = 1

	modified.Spec.Hostname, _ = os.Hostname()

	uname, _ := uname.New()
	modified.Spec.Kernel = fmt.Sprintf("%s %s %s", uname.Sysname(), uname.Machine(), uname.KernelRelease())

	// TODO: Should be patch, but caused error on controller
	return e.Rockferry.Nodes().Create(ctx, modified)
}
