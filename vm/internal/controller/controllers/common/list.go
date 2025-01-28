package common

type ListResponse[T any] struct {
	List []*T `json:"list"`
}
