package storagepool

import "encoding/xml"

type StoragePoolType string

const (
	StoragePoolTypeDir         = "dir"
	StoragePoolTypeFs          = "fs"
	StoragePoolTypeNetfs       = "netfs"
	StoragePoolTypeDisk        = "disk"
	StoragePoolTypeIscsi       = "iscsi"
	StoragePoolTypeLogical     = "logical"
	StoragePoolTypeScsi        = "scsi"
	StoragePoolTypeMpath       = "mpath"
	StoragePoolTypeRbd         = "rbd"
	StoragePoolTypeSheepdog    = "sheepdog"
	StoragePoolTypeGluster     = "gluster"
	StoragePoolTypeZfs         = "zfs"
	StoragePoolTypeVstorage    = "vstorage"
	StoragePoolTypeIscsiDirect = "iscsi-direct"
)

type Schema struct {
	XMLName  xml.Name        `xml:"pool"`
	Type     StoragePoolType `xml:"type,attr"`
	Uuid     string          `xml:"uuid"`
	Name     string          `xml:"name"`
	Features Features        `xml:"features,omitempty"`
	Source   Source          `xml:"source,omitempty"`
	Target   Target          `xml:"target,omitempty"`
}

type Features struct {
	State string `xml:"state,attr"`
}

type Source struct {
	Name     string         `xml:"name,omitempty"`
	Auth     SourceAuth     `xml:"auth,omitempty"`
	Host     SourceHost     `xml:"host,omitempty"`
	Device   SourceDevice   `xml:"device,omitempty"`
	Dir      SourceDir      `xml:"dir,omitempty"`
	Vendor   SourceVendor   `xml:"vendor,omitempty"`
	Product  SourceProduct  `xml:"product,omitempty"`
	Format   SourceFormat   `xml:"format,omitempty"`
	Protocol SourceProtocol `xml:"protocol,omitempty"`
}

type SourceHost struct {
	Name string `xml:"name,attr"`
	Port string `xml:"port,attr,omitempty"`
}

type SourceDevice struct {
	Path          string `xml:"path,attr"`
	PartSeperator string `xml:"path_seperator,attr"`
}

type SourceDir struct {
	Path string `xml:"path,attr"`
}

type SourceAuth struct {
	Type     string             `xml:"type,attr"`
	Username string             `xml:"username,attr"`
	Secrets  []SourceAuthSecret `xml:"secret"`
}

type SourceAuthSecret struct {
	Usage string `xml:"usage,attr"`
	Uuid  string `xml:"uuid,attr"`
}

type SourceVendor struct {
	Name string `xml:"name,attr"`
}

type SourceProduct struct {
	Name string `xml:"name,attr"`
}

type SourceAdapter struct {
	Name            string                  `xml:"name,attr,omitempty"`
	Type            string                  `xml:"type,attr"`
	Wwnn            string                  `xml:"wwnn,attr,omitempty"`
	Wwpn            string                  `xml:"wwpn,attr,omitempty"`
	Parent          string                  `xml:"parent,attr,omitempty"`
	ParentWwnn      string                  `xml:"parent_wwnn,attr,omitempty"`
	ParentWwpn      string                  `xml:"parent_wwpn,attr,omitempty"`
	ParentFabricWwn string                  `xml:"parent_fabric_wwn,omitempty"`
	Managed         string                  `xml:"managed,omitempty"`
	ParentAddr      SourceAdapterParentAddr `xml:"parentaddr,omitempty"`
}

type SourceAdapterParentAddr struct {
	UniqueId string                         `xml:"unique_id,attr,omitempty"`
	Address  SourceAdapterParentAddrAddress `xml:"address"`
}

type SourceAdapterParentAddrAddress struct {
	Domain   int `xml:"domain,attr"`
	Slot     int `xml:"slot,attr"`
	Bus      int `xml:"bus,attr"`
	Function int `xml:"function,attr"`
}

type SourceFormat struct {
	Type string `xml:"type,attr"`
}

type SourceProtocol struct {
	Ver string `xml:"ver,attr"`
}

type Target struct {
	Path        string            `xml:"path"`
	Permissions TargetPermissions `xml:"permissions"`
}

type TargetPermissions struct {
	Owner int    `xml:"owner"`
	Group int    `xml:"group"`
	Mode  int    `xml:"mode"`
	Label string `xml:"label"`
}
