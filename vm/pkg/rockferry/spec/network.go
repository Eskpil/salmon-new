package spec

type NetworkSpecBridge struct {
	Name string `json:"name"`
}

type NetworkSpecForward struct {
	// TODO: Add enums
	Mode string `json:"mode"`
}

type NetworkSpec struct {
	Name    string             `json:"name"`
	Bridge  NetworkSpecBridge  `json:"bridge"`
	Forward NetworkSpecForward `json:"forward"`
}
