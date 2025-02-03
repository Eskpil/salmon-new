package spec

import "github.com/eskpil/salmon/vm/pkg/rockferry/resource"

type MachineSpecInterface struct {
	Mac   string `json:"mac"`
	Model string `json:"model"`

	Network *string `json:"network"`
	Bridge  *string `json:"bridge"`
}

type MachineSpecDiskFile struct {
	Key string `json:"key"`
}

type MachineSpecDiskNetwork struct {
	Hosts []*StoragePoolSpecSourceHost `json:"hosts"`
	Auth  StoragePoolSpecSourceAuth    `json:"auth"`

	Protocol string `json:"type"`
	Key      string `json:"key"`
}

type MachineSpecDisk struct {
	Device string `json:"device"`
	Type   string `json:"type"`

	File    *MachineSpecDiskFile    `json:"file,omitempty"`
	Network *MachineSpecDiskNetwork `json:"network,omitempty"`
}

type MachineSpec struct {
	Name     string            `json:"name"`
	Topology resource.Topology `json:"topology"`

	Disks      []*MachineSpecDisk      `json:"disks"`
	Interfaces []*MachineSpecInterface `json:"interfaces"`
}
