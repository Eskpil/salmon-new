package storagepoolsv1

type SpecSourceHost struct {
	Name string `json:"name"`
	Port string `json:"port"`
}

type SpecSourceAuth struct {
	Type     string `json:"type"`
	Username string `json:"username"`
	Secret   string `json:"secret"`
}

type SpecSource struct {
	Name  string            `json:"name"`
	Hosts []*SpecSourceHost `json:"hosts,omitempty"`
	Auth  *SpecSourceAuth   `json:"auth"`
}

type Spec struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	Capacity   uint64 `json:"capacity"`
	Allocation uint64 `json:"allocation"`
	Available  uint64 `json:"available"`

	Source *SpecSource `json:"source"`
}
