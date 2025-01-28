package db

import (
	"time"

	"github.com/labstack/echo/v4"
	clientv3 "go.etcd.io/etcd/client/v3"
)

const MiddlewareKey = "database"

func Middleware() echo.MiddlewareFunc {
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

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(MiddlewareKey, cli)
			return next(c)
		}
	}
}

func Extract(c echo.Context) *clientv3.Client {
	return c.Get(MiddlewareKey).(*clientv3.Client)
}
