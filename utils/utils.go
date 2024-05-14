package appy_utils

func SetFieldIfNotNull[T any](field *T, value *T) {
	if value != nil {
		*field = *value
	}
}
