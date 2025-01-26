package nodeconnection

import (
	"github.com/eskpil/salmon/vm/internal/controller/config"
	"github.com/eskpil/salmon/vm/nodeapi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Connection struct {
	Config *config.Node
	Client nodeapi.NodeApiClient
}

func New(node *config.Node) (*Connection, error) {
	cc, err := grpc.NewClient(node.Url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := nodeapi.NewNodeApiClient(cc)
	return &Connection{Client: client, Config: node}, nil
}
