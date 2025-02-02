package networksv1

type SpecBridge struct {
	Name string `json:"name"`
}

type SpecForward struct {
	// TODO: Add enums
	Mode string `json:"mode"`
}

type Spec struct {
	Name    string      `json:"name"`
	Bridge  SpecBridge  `json:"bridge"`
	Forward SpecForward `json:"forward"`
}
