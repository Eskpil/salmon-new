package spec

type StoragePoolSpecSourceHost struct {
	Name string `json:"name"`
	Port string `json:"port"`
}

type StoragePoolSpecSourceAuth struct {
	Type     string `json:"type"`
	Username string `json:"username"`
	Secret   string `json:"secret"`
}

type StoragePoolSpecSource struct {
	Name  string                       `json:"name"`
	Hosts []*StoragePoolSpecSourceHost `json:"hosts,omitempty"`
	Auth  *StoragePoolSpecSourceAuth   `json:"auth"`
}

type StoragePoolSpec struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	Capacity   uint64 `json:"capacity"`
	Allocation uint64 `json:"allocation"`
	Available  uint64 `json:"available"`

	Source *StoragePoolSpecSource `json:"source"`
}
