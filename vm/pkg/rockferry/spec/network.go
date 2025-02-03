package spec

type NetworkSpecBridge struct {
	Name  string `json:"name"`
	Stp   string `json:"stp"`
	Delay string `json:"delay"`
}

type NetworkSpecForwardNat struct {
	AddressStart string `json:"address_start"`
	AddressEnd   string `json:"address_end"`

	PortStart string `json:"port_start"`
	PortEnd   string `json:"port_end"`
}

type NetworkSpecForward struct {
	// TODO: Add enums
	Mode string                 `json:"mode"`
	Dev  string                 `json:"dev"`
	Nat  *NetworkSpecForwardNat `json:"nat,omitempty"`
}

type NetworkSpec struct {
	Name    string             `json:"name"`
	Bridge  NetworkSpecBridge  `json:"bridge"`
	Forward NetworkSpecForward `json:"forward"`
	Ipv6    bool               `json:"ipv6"`
	Mtu     uint64             `json:"mtu"`
}
