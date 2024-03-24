package appy_utils

func ErrorIsAnyOf(err error, errs ...error) bool {
	for _, e := range errs {
		if err == e {
			return true
		}
	}
	return false
}
