package main

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/eskpil/salmon/vm/internal/controller/config"
	"github.com/eskpil/salmon/vm/internal/controller/controllers/nodes"
	"github.com/eskpil/salmon/vm/internal/controller/state"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gopkg.in/yaml.v3"

	"go.etcd.io/etcd/server/v3/embed"
)

func runDb(dir string) {
	cfg := embed.NewConfig()
	cfg.Dir = dir
	e, err := embed.StartEtcd(cfg)
	if err != nil {
		log.Fatal(err)
	}

	defer e.Close()
	select {
	case <-e.Server.ReadyNotify():
		log.Printf("Database is running")
	case <-time.After(60 * time.Second):
		e.Server.Stop() // trigger a shutdown
		log.Printf("Server took too long to start!")
	}

	log.Fatal(<-e.Err())
}

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
	wg := new(sync.WaitGroup)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		runDb("salmon_vm.etcd")
		wg.Done()
	}(wg)

	config := readConfig()

	s, err := state.New(config)
	if err != nil {
		panic(err)
	}

	server := echo.New()

	server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	server.GET("/v1/nodes", nodes.List(s))

	if err := server.Start("0.0.0.0:8080"); err != nil {
		panic(err)
	}

	wg.Wait()
}
