package appy_utils

type SelectorFn[T any, R any] func(*T) *R
type MapFn[T any, R any] func(*T, *R)

func SetFieldIfNotNull[T any](field *T, value *T) {
	if value != nil {
		*field = *value
	}
}

func MapEntries[T any, R any](entriesToMap []T, mappables []R, mapFn MapFn[T, R]) {
	for i := range entriesToMap {
		for j := range mappables {
			mapFn(&entriesToMap[i], &mappables[j])
		}
	}
}

func CompareCloseEnough(a, b int64, window int64) bool {
	if a > b {
		return a-b <= window
	}

	return b-a <= window
}
