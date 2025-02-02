package storagevolumesv1

type Spec struct {
	Name       string `json:"name"`
	Capacity   uint64 `json:"capacity"`
	Allocation uint64 `json:"allocation"`
	Key        string `json:"key"`
	Pool       string `json:"pool"`
}
