package nodes

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/eskpil/salmon/vm/internal/controller/controllers/common"
	"github.com/eskpil/salmon/vm/internal/controller/db"
	"github.com/eskpil/salmon/vm/internal/controller/models"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type CreateNodeInput struct {
	Url string `json:"url"`
}

func Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
		defer cancel()

		input := new(CreateNodeInput)
		if err := c.Bind(input); err != nil {
			return c.JSON(http.StatusBadRequest, common.MalformedInput())
		}

		db := db.Extract(c)

		id := uuid.New()

		resource := new(models.Resource)

		resource.Id = id.String()
		resource.Kind = models.ResourceKindNode

		resource.Annotations = make(map[string]string)
		resource.Annotations["node.url"] = input.Url

		path := fmt.Sprintf("%s/%s/%s", models.RootKey, models.ResourceKindNode, id)
		res, err := db.Put(ctx, path, string(resource.Marshal()))
		if err != nil {
			panic(err)
		}

		_ = res

		return c.String(http.StatusCreated, "made")
	}
}
