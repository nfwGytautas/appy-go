package appy_utils

func TrueOnAny[T comparable](value T, values ...T) bool {
	for _, v := range values {
		if value == v {
			return true
		}
	}
	return false
}

func TrueOnAnyNonNil[T comparable](value *T, values ...T) bool {
	if value == nil {
		return false
	}
	return TrueOnAny(*value, values...)
}
