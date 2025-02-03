package spec

import "github.com/eskpil/salmon/vm/pkg/rockferry/resource"

type MachineRequestSpecDisk struct {
	Pool       string `json:"pool"`
	Capacity   uint64 `json:"capacity"`
	Allocation uint64 `json:"allocation"`
}

type MachineRequestSpecCdrom struct {
	Key string `json:"key"`
}

type MachineRequestSpec struct {
	Name     string                    `json:"name"`
	Topology resource.Topology         `json:"topology"`
	Network  string                    `json:"network"`
	Disks    []*MachineRequestSpecDisk `json:"disks"`
	Cdrom    *MachineRequestSpecCdrom  `json:"cdrom"`
}
