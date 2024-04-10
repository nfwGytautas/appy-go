package appy_utils

import (
	"errors"
)

func ErrorIsAnyOf(err error, errs ...error) bool {
	for _, e := range errs {
		if err == e {
			return true
		}
	}
	return false
}

func ErrorIsTypeOf[T any](err error) bool {
	var t T
	return errors.As(err, &t)
}
