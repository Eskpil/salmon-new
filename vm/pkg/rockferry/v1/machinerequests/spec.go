package machinerequestsv1

import "github.com/eskpil/salmon/vm/pkg/rockferry/resource"

type SpecDisk struct {
	Pool       string `json:"pool"`
	Capacity   uint64 `json:"capacity"`
	Allocation uint64 `json:"allocation"`
}

type SpecCdrom struct {
	Key string `json:"key"`
}

type Spec struct {
	Name     string            `json:"name"`
	Topology resource.Topology `json:"topology"`
	Network  string            `json:"network"`
	Disks    []*SpecDisk       `json:"disks"`
	Cdrom    *SpecCdrom        `json:"cdrom"`
}
