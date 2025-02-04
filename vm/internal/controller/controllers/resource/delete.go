package resource

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/eskpil/salmon/vm/internal/controller/controllers/common"
	"github.com/eskpil/salmon/vm/internal/controller/db"
	"github.com/eskpil/salmon/vm/internal/controller/models"
	"github.com/labstack/echo/v4"
)

type DeleteRequest struct {
	Kind string `json:"kind"`
	Id   string `json:"id"`
}

func Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
		defer cancel()

		req := new(DeleteRequest)
		if err := c.Bind(req); err != nil {
			fmt.Println(err)
			return c.JSON(http.StatusBadRequest, common.MalformedInput())
		}

		db := db.Extract(c)

		path := fmt.Sprintf("%s/%s/%s", models.RootKey, req.Kind, req.Id)

		fmt.Println("deleting", path)

		_, err := db.Delete(ctx, path)
		if err != nil {
			fmt.Println("failed to delete resource", err)
			return c.JSON(http.StatusInternalServerError, common.InternalServerError())
		}

		return c.JSON(http.StatusOK, res{Ok: true})
	}
}
