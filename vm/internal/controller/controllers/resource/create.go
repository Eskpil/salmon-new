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

		// TODO: Owner is expected to be nil at some stage, for example
		// 		 when trying to create a vm and you do not care where,
		// 		 a node needs to be allocated for both volume creation
		// 		 and the vm.

		resource := new(models.Resource)

		id := uuid.NewString()
		if input.Kind == models.ResourceKindStorageVolume {
			id = fmt.Sprintf("%s/%s", input.OwnerRef.Id, input.Spec.(map[string]interface{})["name"].(string))
		}

		resource.Id = id
		resource.Owner = input.OwnerRef
		resource.Kind = string(input.Kind)
		resource.Annotations = input.Annotations
		resource.Spec = input.Spec
		resource.Status.Phase = models.PhaseRequested

		path := fmt.Sprintf("%s/%s/%s", models.RootKey, resource.Kind, resource.Id)

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
