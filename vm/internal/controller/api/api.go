package api

import (
	"time"

	"github.com/eskpil/salmon/vm/controllerapi"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type Controller struct {
	controllerapi.UnimplementedControllerApiServer
	Db *clientv3.Client
}

func New() (Controller, error) {
	controller := new(Controller)

	// TODO: Enable some kind of config
	// TODO: Avoid multiple db connections
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})

	// TODO: Avoid panic
	if err != nil {
		panic(err)
	}

	controller.Db = cli

	return *controller, nil
}
