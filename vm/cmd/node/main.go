package main

import (
	"context"

	"github.com/eskpil/salmon/vm/internal/node"
)

func main() {
	state, err := node.New()
	if err != nil {
		panic(err)
	}

	if err := state.Watch(context.Background()); err != nil {
		panic(err)
	}
}
