package status

import (
	"fmt"
	"strconv"
	"strings"
)

type ErrorCode = int

const (
	ErrNoResults ErrorCode = iota
)

func NewError(code ErrorCode, message string) error {
	return fmt.Errorf("%d: %s", code, message)
}

func Is(err error, code ErrorCode) bool {
	parts := strings.Split(err.Error(), ":")
	if 1 >= len(parts) {
		return false
	}

	if possibleCode, err := strconv.ParseInt(parts[0], 10, 8); err != nil {
		return false
	} else {
		if code != ErrorCode(possibleCode) {
			return false
		}
	}

	return true
}
