package appy_utils

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

func SelectFromArray[T any, R any](entries []T, selectorFn SelectorFn[T, R]) []R {
	result := make([]R, 0, len(entries))
	for i := range entries {
		res := selectorFn(&entries[i])
		if res != nil {
			result = append(result, *res)
		}
	}

	return result
}

func JoinArrays[T any](arr []T, element T) []T {
	return append(arr, element)
}
