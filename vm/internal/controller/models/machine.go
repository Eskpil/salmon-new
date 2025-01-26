package models

const MachinesKey = "machines"

type Machine struct {
	Uuid string `json:"uuid"`
	Name string `json:"name"`
	Node string `json:"node"`
}
