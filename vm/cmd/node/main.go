package main

import (
	"context"
	"os"

	"github.com/eskpil/salmon/vm/internal/node"
	"github.com/eskpil/salmon/vm/internal/node/config"
	"gopkg.in/yaml.v2"
)

func readConfig() *config.Config {
	contents, err := os.ReadFile("./config.yml")
	if err != nil {
		panic(err)
	}

	config := new(config.Config)

	if err := yaml.Unmarshal(contents, config); err != nil {
		panic(err)
	}

	return config
}

func main() {
	state, err := node.New(readConfig())
	if err != nil {
		panic(err)
	}

	if err := state.Watch(context.Background()); err != nil {
		panic(err)
	}
}
