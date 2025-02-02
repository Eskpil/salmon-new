package queries

import (
	"encoding/json"
	"encoding/xml"

	"github.com/digitalocean/go-libvirt"
	"github.com/eskpil/salmon/vm/nodeapi"
	storagevolumesv1 "github.com/eskpil/salmon/vm/pkg/rockferry/v1/storagevolumes"
	"github.com/eskpil/salmon/vm/pkg/virtwrap/storagepool"
	"github.com/eskpil/salmon/vm/pkg/virtwrap/storagevol"
	"github.com/google/uuid"
)

type StoragePool struct {
	Name             string
	Uuid             uuid.UUID
	AllocatedVolumes uint64
	Capacity         uint64
	Available        uint64
	State            uint8
}

type StorageVolume struct {
	Name       string
	Pool       string
	Key        string
	Type       uint8
	Allocation uint64
	Capacity   uint64
}

func listAllStoragePools(v *libvirt.Libvirt) ([]libvirt.StoragePool, error) {
	pools, _, err := v.ConnectListAllStoragePools(100, 1|2)
	if err != nil {
		return nil, err
	}
	return pools, nil
}

func completeStoragePool(v *libvirt.Libvirt, pool libvirt.StoragePool) (*nodeapi.StoragePool, error) {
	state, capacity, allocated, avaliable, err := v.StoragePoolGetInfo(pool)
	if err != nil {
		return nil, err
	}

	xmlSchema, err := v.StoragePoolGetXMLDesc(pool, 0)
	if err != nil {
		return nil, err
	}

	var schema storagepool.Schema
	if err := xml.Unmarshal([]byte(xmlSchema), &schema); err != nil {
		return nil, err
	}

	schemaJson, err := json.Marshal(schema)
	if err != nil {
		return nil, err
	}

	mapped := new(nodeapi.StoragePool)

	mapped.Kind = string(schema.Type)
	mapped.Uuid = schema.Uuid
	mapped.Schema = schemaJson
	mapped.State = uint64(state)
	mapped.Capacity = capacity
	mapped.Allocated = allocated
	mapped.Available = avaliable
	mapped.Name = pool.Name

	return mapped, nil
}

func completeVolume(v *libvirt.Libvirt, unmappedVolume libvirt.StorageVol, poolUuid string) (*nodeapi.StorageVolume, error) {
	mapped := new(nodeapi.StorageVolume)
	xmlSchema, err := v.StorageVolGetXMLDesc(unmappedVolume, 0)
	if err != nil {
		return nil, err
	}

	var schema storagevol.Schema
	if err := xml.Unmarshal([]byte(xmlSchema), &schema); err != nil {
		return nil, err
	}

	schemaJson, err := json.Marshal(schema)
	if err != nil {
		return nil, err
	}

	mapped.Schema = schemaJson
	mapped.Pool = poolUuid
	mapped.Name = unmappedVolume.Name
	mapped.Key = unmappedVolume.Key
	mapped.Capacity = uint64(schema.Capacity.Value)
	mapped.Allocation = uint64(schema.Allocation.Value)

	return mapped, nil
}

func completeVolumes(v *libvirt.Libvirt, pool libvirt.StoragePool, uuid string) ([]*nodeapi.StorageVolume, error) {
	unmappedVolumes, _, err := v.StoragePoolListAllVolumes(pool, 100, 0)
	if err != nil {
		return nil, err
	}

	volumes := make([]*nodeapi.StorageVolume, len(unmappedVolumes))

	for i, unmappedVolume := range unmappedVolumes {
		volume, err := completeVolume(v, unmappedVolume, uuid)
		if err != nil {
			continue
		}

		volumes[i] = volume
	}

	return volumes, nil
}

func (c *Client) QueryStoragePools() ([]*nodeapi.StoragePool, error) {
	return c.pools, nil
}

func (c *Client) QueryStorageVolumes() ([]*nodeapi.StorageVolume, error) {
	return c.volumes, nil
}

func (c *Client) preloadStorageVolumes(pool *nodeapi.StoragePool, virtPool libvirt.StoragePool, uuid string) error {
	v, err := completeVolumes(c.v, virtPool, uuid)
	if err != nil {
		return err
	}

	pool.AllocatedVolumes = uint64(len(v))

	c.volumes = append(c.volumes, v...)

	return nil
}

func (c *Client) preloadStorage() error {
	unmappedPools, err := listAllStoragePools(c.v)
	if err != nil {
		return err
	}

	// TODO: this is inefficient, use map
	c.pools = nil
	c.volumes = nil

	for _, unmappedPool := range unmappedPools {
		pool, err := completeStoragePool(c.v, unmappedPool)
		if err != nil {
			continue
		}

		c.pools = append(c.pools, pool)
	}

	for i, pool := range c.pools {
		unmappedPool := unmappedPools[i]
		if err := c.preloadStorageVolumes(pool, unmappedPool, pool.GetUuid()); err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) CreateVolume(poolUuid string, name string, format string, capacity uint64, allocation uint64) error {
	t, err := c.findPoolByUuid(poolUuid)
	if err != nil {
		return err
	}

	pool, err := c.v.StoragePoolLookupByName(t.Name)
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

	return c.preloadStorage()
}

func (c *Client) QueryVolumeSpec(poolUuid string, name string) (*storagevolumesv1.Spec, error) {
	t, err := c.findPoolByUuid(poolUuid)
	if err != nil {
		return nil, err
	}

	pool, err := c.v.StoragePoolLookupByName(t.Name)
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

	spec := new(storagevolumesv1.Spec)

	spec.Key = xmlSchema.Key
	spec.Name = xmlSchema.Name
	spec.Allocation = uint64(xmlSchema.Allocation.Value)
	spec.Capacity = uint64(xmlSchema.Capacity.Value)

	return spec, nil
}
