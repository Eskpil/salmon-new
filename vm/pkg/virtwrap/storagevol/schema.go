package storagevol

import (
	"encoding/xml"
)

type Size struct {
	Unit  string `xml:"unit,attr,omitempty"`
	Value int    `xml:",chardata"`
}

type Annotations struct {
	Id string `xml:"id"`
}

type Schema struct {
	XMLName      xml.Name     `xml:"volume"`
	Name         string       `xml:"name"`
	Annotations  *Annotations `xml:"annotations,omitempty"`
	Key          string       `xml:"key,omitempty"`
	Target       Target       `xml:"target,omitempty"`
	BackingStore BackingStore `xml:"backingStore,omitempty"`
	Allocation   Size         `xml:"allocation,omitempty"`
	Capacity     Size         `xml:"capacity,omitempty"`
}

type Target struct {
	Path        string      `xml:"path,omitempty"`
	Format      Format      `xml:"format,omitempty"`
	Permissions Permissions `xml:"permissions,omitempty"`
	Timestamps  Timestamps  `xml:"timestamps,omitempty"`

	//	TODO: Figure out why omitempty not working
	//	Encryption  Encryption  `xml:"encryption,omitempty"`
	Compat      string      `xml:"compat,omitempty"`
	Nocow       bool        `xml:"nocow,omitempty"`
	ClusterSize ClusterSize `xml:"clusterSize,omitempty"`
	Features    Features    `xml:"features,omitempty"`
}

type Format struct {
	Type string `xml:"type,attr"`
}

type Permissions struct {
	Owner int    `xml:"owner"`
	Group int    `xml:"group"`
	Mode  int    `xml:"mode"`
	Label string `xml:"label"`
}

type Timestamps struct {
	ATime float64 `xml:"atime"`
	MTime float64 `xml:"mtime"`
	CTime float64 `xml:"ctime"`
}

// TODO: Implement encryption
type Encryption struct {
	Format string `xml:"format,attr"`
}

type ClusterSize struct {
	Unit  string `xml:"unit,attr"`
	Value int    `xml:",chardata"`
}

type Features struct {
	LazyRefcounts bool `xml:"lazy_refcounts,omitempty"`
	Extended12    bool `xml:"extended_12,omitempty"`
}

type BackingStore struct {
	Path        string      `xml:"path"`
	Format      Format      `xml:"format"`
	Permissions Permissions `xml:"permissions"`
}
