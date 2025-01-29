package resource

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

type CreateResourceInput struct {
	Annotations map[string]string   `json:"annotations"`
	Kind        models.ResourceKind `json:"kind"`
	OwnerRef    *models.OwnerRef    `json:"owner_ref"`
	Spec        any                 `json:"spec"`
}

type res struct {
	Ok bool `json:"ok"`
}

func Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
		defer cancel()

		input := new(CreateResourceInput)
		if err := c.Bind(input); err != nil {
			fmt.Println(err)
			return c.JSON(http.StatusBadRequest, common.MalformedInput())
		}

		resource := new(models.Resource)

		resource.Id = uuid.NewString()
		resource.Owner = input.OwnerRef
		resource.Kind = string(input.Kind)
		resource.Annotations = input.Annotations
		resource.Spec = input.Spec

		path := fmt.Sprintf("%s/%s/%s", models.RootKey, resource.Kind, resource.Id)
		if resource.Kind == models.ResourceKindStorageVolume {
			path = fmt.Sprintf("%s/%s/%s/%s", models.RootKey, resource.Kind, resource.Owner.Id, resource.Spec.(map[string]interface{})["name"].(string))
		}

		db := db.Extract(c)

		if _, err := db.Put(ctx, path, string(resource.Marshal())); err != nil {
			fmt.Println(err)
			return c.JSON(http.StatusInternalServerError, common.InternalServerError())
		}

		response := new(res)
		response.Ok = true

		return c.JSON(http.StatusCreated, response)
	}
}
