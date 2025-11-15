package xdb

import "errors"

var ErrInvalidIPAddress = errors.New("invalid ip address")

func IsInvalidIPAddress(err error) bool {
	return errors.Is(err, ErrInvalidIPAddress)
}
