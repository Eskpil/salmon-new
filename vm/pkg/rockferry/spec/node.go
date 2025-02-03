package spec

import (
	"time"

	"github.com/eskpil/salmon/vm/pkg/rockferry/resource"
)

type NodeSpec struct {
	Topology resource.Topology `json:"topology"`

	Hostname string    `json:"hostname"`
	Kernel   string    `json:"kernel"`
	UpSince  time.Time `json:"up_since"`

	ActiveMachines uint64 `json:"active_machines"`
	TotalMachines  uint64 `json:"total_machines"`
}
