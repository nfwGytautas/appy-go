package appy_utils

func SetFieldIfNotNull[T any](field *T, value *T) {
	if value != nil {
		*field = *value
	}
}

func InArray[T comparable](value T, array []T) bool {
	for _, v := range array {
		if v == value {
			return true
		}
	}

	return false
}

func FilterArray[T comparable](array []T, filter func(T) bool) []T {
	var result []T

	for _, v := range array {
		if filter(v) {
			result = append(result, v)
		}
	}

	return result
}
