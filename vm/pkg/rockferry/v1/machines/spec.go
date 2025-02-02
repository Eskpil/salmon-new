package machinesv1

import "github.com/eskpil/salmon/vm/pkg/rockferry/resource"

type SpecInterface struct {
	Mac   string `json:"mac"`
	Model string `json:"model"`

	Network *string `json:"network"`
	Bridge  *string `json:"bridge"`
}

type SpecDisk struct {
	Device string `json:"device"`
	Type   string `json:"type"`

	// The volume key
	Volume string `json:"volume"`
}

type Spec struct {
	Name     string            `json:"name"`
	Topology resource.Topology `json:"topology"`

	Disks      []*SpecDisk      `json:"disks"`
	Interfaces []*SpecInterface `json:"interfaces"`
}
