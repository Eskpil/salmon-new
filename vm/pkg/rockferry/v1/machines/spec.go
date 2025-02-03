package machinesv1

import (
	"github.com/eskpil/salmon/vm/pkg/rockferry/resource"
	storagepoolsv1 "github.com/eskpil/salmon/vm/pkg/rockferry/v1/storagepools"
)

type SpecInterface struct {
	Mac   string `json:"mac"`
	Model string `json:"model"`

	Network *string `json:"network"`
	Bridge  *string `json:"bridge"`
}

type SpecDiskFile struct {
	Key string `json:"key"`
}

type SpecDiskNetwork struct {
	Hosts []*storagepoolsv1.SpecSourceHost `json:"hosts"`
	Auth  *storagepoolsv1.SpecSourceAuth   `json:"auth"`

	Protocol string `json:"type"`
	Key      string `json:"key"`
}

type SpecDisk struct {
	Device string `json:"device"`
	Type   string `json:"type"`

	File    *SpecDiskFile    `json:"file,omitempty"`
	Network *SpecDiskNetwork `json:"network,omitempty"`
}

type Spec struct {
	Name     string            `json:"name"`
	Topology resource.Topology `json:"topology"`

	Disks      []*SpecDisk      `json:"disks"`
	Interfaces []*SpecInterface `json:"interfaces"`
}
