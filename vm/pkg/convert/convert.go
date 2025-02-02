package convert

import (
	"encoding/json"

	"google.golang.org/protobuf/types/known/structpb"
)

func Convert[T any](in *structpb.Struct) (*T, error) {
	out := new(T)

	encoded, err := in.MarshalJSON()
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(encoded, out); err != nil {
		return nil, err
	}

	return out, nil
}
