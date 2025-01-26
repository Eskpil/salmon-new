package network

import "encoding/xml"

type Schema struct {
	XMLName             xml.Name `xml:"network"`
	Name                string   `xml:"name"`
	TrustGuestRxFilters string   `xml:"trustGuestRxFilters,attr"`
	Ipv6                string   `xml:"ipv6,attr"`
	Uuid                string   `xml:"uuid"`
	Description         string   `xml:"description"`
	Bridge              Bridge   `xml:"bridge"`
	Mtu                 Mtu      `xml:"mtu"`
	Domain              Domain   `xml:"domain"`
	Forward             Forward  `xml:"forward"`
}

type Mtu struct {
	Size int `xml:"size,attr"`
}

type Bridge struct {
	Name            string `xml:"name,attr"`
	Stp             string `xml:"stp,attr"`
	Delay           string `xml:"delay,attr"`
	MacTableManager string `xml:"macTableManager,attr"`
	Zone            string `xml:"zone,attr"`
}

type Domain struct {
	Name      string `xml:"name,attr"`
	LocalOnly string `xml:"localOnly,attr"`
	Register  string `xml:"register,attr"`
}

// https://libvirt.org/formatnetwork.html#nat-based-network
// TODO: Support bridge, private, route etc
type Forward struct {
	Mode string      `xml:"mode,attr"`
	Dev  string      `xml:"dev,attr"`
	Nat  *ForwardNat `xml:"nat,omitempty" json:",omitempty"`
}

type ForwardNat struct {
	Address ForwardNatAddress `xml:"address"`
	Port    ForwardNatPort    `xml:"port"`
}

type ForwardNatAddress struct {
	Start string `xml:"start"`
	End   string `xml:"end"`
}

type ForwardNatPort struct {
	Start string `xml:"start"`
	End   string `xml:"end"`
}
