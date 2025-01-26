package state

import (
	"github.com/eskpil/salmon/vm/internal/controller/config"
	"github.com/eskpil/salmon/vm/internal/controller/nodeconnection"
)

type State struct {
	NodeConnections []*nodeconnection.Connection
}

func New(config *config.Config) (*State, error) {
	s := new(State)
	for _, n := range config.Nodes {
		connection, err := nodeconnection.New(&n)
		if err != nil {
			return nil, err
		}
		s.NodeConnections = append(s.NodeConnections, connection)
	}

	return s, nil
}
