package storagepool

type StoragePoolSpec struct {
	//	string name = 1;
	//
	// string uuid = 2;
	// bytes schema = 3;
	//
	// uint64 allocated_volumes = 4;
	// uint64 capacity = 5;
	// uint64 available = 6;
	// uint64 state = 7;
	//
	// string kind = 8;
	// uint64 allocated = 9;

	Name    string `json:"name"`
	Volumes uint64 `json:"volumes"`

	Capacity  uint64 `json:"capacity"`
	Available uint64 `json:"available"`
	Allocated uint64 `json:"allocated"`

	State uint64 `json:"state"`
	Kind  string `json:"kind"`
}
