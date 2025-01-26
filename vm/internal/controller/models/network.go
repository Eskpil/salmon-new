package models

const NetworksKey string = "networks/"

type Network struct {
	Uuid string `json:"uuid"`
	Name string `json:"name"`
}
