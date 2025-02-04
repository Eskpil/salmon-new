package queries

import (
	"encoding/xml"
	"fmt"

	"github.com/digitalocean/go-libvirt"
	"github.com/eskpil/salmon/vm/pkg/rockferry"
	"github.com/eskpil/salmon/vm/pkg/rockferry/resource"
	"github.com/eskpil/salmon/vm/pkg/rockferry/spec"
	"github.com/eskpil/salmon/vm/pkg/virtwrap/storagepool"
	"github.com/eskpil/salmon/vm/pkg/virtwrap/storagevol"
	"github.com/google/uuid"
)

func (c *Client) CreateVolume(poolName string, name string, format string, capacity uint64, allocation uint64) error {
	pool, err := c.v.StoragePoolLookupByName(poolName)
	if err != nil {
		return err
	}

	volume := new(storagevol.Schema)

	volume.Name = name
	volume.XMLName.Space = "volume"

	volume.Allocation.Unit = "bytes"
	volume.Allocation.Value = int(allocation)

	// TODO: Just for testing
	volume.Capacity.Unit = "bytes"
	volume.Capacity.Value = int(capacity)

	volume.Target.Format.Type = format

	volume.Annotations = new(storagevol.Annotations)
	volume.Annotations.Id = uuid.NewString()

	volumeXML, err := xml.Marshal(volume)
	if err != nil {
		return err
	}

	_, err = c.v.StorageVolCreateXML(pool, string(volumeXML), 0)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) QueryVolumeSpec(poolName string, name string) (*spec.StorageVolumeSpec, error) {
	pool, err := c.v.StoragePoolLookupByName(poolName)
	if err != nil {
		return nil, err
	}

	vol, err := c.v.StorageVolLookupByName(pool, name)
	if err != nil {
		return nil, err
	}

	xmlDesc, err := c.v.StorageVolGetXMLDesc(vol, 0)
	if err != nil {
		return nil, err
	}

	xmlSchema := new(storagevol.Schema)
	if err := xml.Unmarshal([]byte(xmlDesc), xmlSchema); err != nil {
		return nil, err
	}

	spec := new(spec.StorageVolumeSpec)

	spec.Key = xmlSchema.Key
	spec.Name = xmlSchema.Name
	spec.Allocation = uint64(xmlSchema.Allocation.Value)
	spec.Capacity = uint64(xmlSchema.Capacity.Value)

	return spec, nil
}

func (c *Client) QueryStorageVolumes() ([]*rockferry.StorageVolume, error) {
	pools, _, err := c.v.ConnectListAllStoragePools(100, 0)
	if err != nil {
		return nil, err
	}

	volumes := []*rockferry.StorageVolume{}

	for _, pool := range pools {
		names, err := c.v.StoragePoolListVolumes(pool, 100)
		if err != nil {
			return nil, err
		}

		poolId := uuid.UUID([16]byte(pool.UUID))

		for _, name := range names {
			volume := new(rockferry.StorageVolume)

			volume.Id = fmt.Sprintf("%s/%s", poolId.String(), name)
			volume.Owner = new(resource.OwnerRef)
			volume.Annotations = map[string]string{}
			volume.Annotations["origin"] = "sync"

			volume.Owner.Kind = resource.ResourceKindStoragePool
			volume.Owner.Id = poolId.String()
			volume.Kind = resource.ResourceKindStorageVolume
			volume.Status.Phase = resource.PhaseCreated

			volumeSpec, err := c.QueryVolumeSpec(pool.Name, name)
			if err != nil {
				return nil, err
			}

			volume.Spec = *volumeSpec

			volumes = append(volumes, volume)
		}

	}
	return volumes, nil
}

func (c *Client) QueryStoragePools() ([]*rockferry.StoragePool, error) {
	unmapped, _, err := c.v.ConnectListAllStoragePools(100, 0)
	if err != nil {
		return nil, err
	}

	out := []*rockferry.StoragePool{}

	for _, u := range unmapped {
		xmlDesc, err := c.v.StoragePoolGetXMLDesc(u, 0)
		if err != nil {
			return nil, err
		}

		xmlSchema := new(storagepool.Schema)
		if err := xml.Unmarshal([]byte(xmlDesc), xmlSchema); err != nil {
			return nil, err
		}

		_, capacity, allocation, avaliable, err := c.v.StoragePoolGetInfo(u)
		if err != nil {
			return nil, err
		}

		storagePoolSpec := new(spec.StoragePoolSpec)
		storagePoolSpec.Name = xmlSchema.Name

		storagePoolSpec.Type = string(xmlSchema.Type)
		storagePoolSpec.Allocation = allocation
		storagePoolSpec.Capacity = capacity
		storagePoolSpec.Available = avaliable

		storagePoolSpec.Source = new(spec.StoragePoolSpecSource)

		storagePoolSpec.Source.Name = xmlSchema.Source.Name

		host := new(spec.StoragePoolSpecSourceHost)
		host.Name = xmlSchema.Source.Host.Name
		host.Port = xmlSchema.Source.Host.Port
		storagePoolSpec.Source.Hosts = append(storagePoolSpec.Source.Hosts, host)

		auth := new(spec.StoragePoolSpecSourceAuth)

		if xmlSchema.Source.Auth.Type != "" {
			auth.Type = xmlSchema.Source.Auth.Type
			auth.Username = xmlSchema.Source.Auth.Username
			auth.Secret = xmlSchema.Source.Auth.Secrets[0].Uuid

			storagePoolSpec.Source.Auth = auth
		}

		res := new(rockferry.StoragePool)

		res.Id = xmlSchema.Uuid
		res.Owner = new(resource.OwnerRef)
		res.Owner.Kind = resource.ResourceKindNode

		res.Spec = *storagePoolSpec

		res.Kind = resource.ResourceKindStoragePool
		res.Status.Phase = resource.PhaseCreated

		out = append(out, res)
	}

	return out, nil
}

func (c *Client) DeleteStorageVolume(key string) error {
	vol, err := c.v.StorageVolLookupByKey(key)
	if err != nil {
		return err
	}

	return c.v.StorageVolDelete(vol, libvirt.StorageVolDeleteZeroed)
}
