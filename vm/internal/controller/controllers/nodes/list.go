package nodes

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/eskpil/salmon/vm/internal/controller/models"
	"github.com/eskpil/salmon/vm/internal/controller/state"
	"github.com/eskpil/salmon/vm/nodeapi"
	"github.com/labstack/echo/v4"
)

type ListResponse struct {
	List []*models.Node `json:"list"`
}

func listNodes(ctx context.Context, s *state.State, res *ListResponse) {
	for _, conn := range s.NodeConnections {
		node := new(models.Node)

		node.Name = conn.Config.Name
		node.Url = conn.Config.Url

		pingRes, err := conn.Client.Ping(ctx, new(nodeapi.PingRequest))
		if err != nil {
			node.Active = false
			fmt.Println(err)
		} else {
			node.Active = true

			node.ActiveMachines = &pingRes.ActiveMachines
			node.TotalMachines = &pingRes.TotalMachines
		}

		res.List = append(res.List, node)
	}
}

func List(s *state.State) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
		defer cancel()

		res := new(ListResponse)

		// TODO: Implement graceful errors
		listNodes(ctx, s, res)
		return c.JSON(http.StatusOK, res)
	}
}
