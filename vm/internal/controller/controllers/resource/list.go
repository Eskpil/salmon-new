package resource

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/eskpil/salmon/vm/internal/controller/controllers/common"
	"github.com/eskpil/salmon/vm/internal/controller/db"
	"github.com/eskpil/salmon/vm/internal/controller/models"
	"github.com/labstack/echo/v4"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type ListFilter struct {
	Kind      string `query:"kind"`
	Id        string `query:"id"`
	OwnerKind string `query:"owner_kind"`
	OwnerId   string `query:"owner_id"`
}

func List() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
		defer cancel()

		filter := new(ListFilter)
		if err := c.Bind(filter); err != nil {
			return c.JSON(http.StatusBadRequest, common.MalformedInput())
		}

		db := db.Extract(c)

		// TODO: Validate kind

		var opts []clientv3.OpOption

		if filter.Id == "" {
			opts = append(opts, clientv3.WithPrefix())
		}
		path := fmt.Sprintf("%s/%s/%s", models.RootKey, filter.Kind, filter.Id)

		// TODO: Avoid this hack
		if filter.Kind == models.ResourceKindStorageVolume {
			path = fmt.Sprintf("%s/%s/%s", models.RootKey, filter.Kind, filter.OwnerId)
		}

		res, err := db.Get(ctx, path, opts...)
		if err != nil {
			fmt.Println("failed to fetch resources", err)
			return c.JSON(http.StatusInternalServerError, common.InternalServerError())
		}

		list := new(common.ListResponse[models.Resource])

		for _, kv := range res.Kvs {
			resource := new(models.Resource)
			if err := json.Unmarshal(kv.Value, resource); err != nil {
				fmt.Println("unable to unmarshal resource", err)
				return c.JSON(http.StatusInternalServerError, common.InternalServerError())
			}

			if filter.OwnerId != "" && filter.OwnerKind != "" && resource.Owner != nil {
				if filter.OwnerId != resource.Owner.Id && filter.OwnerKind != resource.Owner.Kind {
					continue
				}
			}

			list.List = append(list.List, resource)

		}

		return c.JSON(http.StatusOK, list)
	}
}
