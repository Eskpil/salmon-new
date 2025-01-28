package common

import "net/http"

type Error struct {
	Code    int    `json:"code"`
	Message string `message:"message"`
}

func MalformedInput() Error {
	return Error{
		Code:    http.StatusBadRequest,
		Message: "malformed input request",
	}
}

func InternalServerError() Error {
	return Error{
		Code:    http.StatusInternalServerError,
		Message: "Something went wrong...",
	}
}
