package models

const NodesKey string = "nodes"

// TODO: Use this instead of config file
type Node struct {
	Name string `json:"name"`
	Url  string `json:"url"`

	Active bool `json:"active"`

	ActiveMachines *uint64 `json:"active_machines,omitempty"`
	TotalMachines  *uint64 `json:"total_machines,omitempty"`
}
