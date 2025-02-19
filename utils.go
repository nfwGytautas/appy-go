package appy

import (
	"crypto/rand"
	"encoding/base64"
)

func isElementInArray[T comparable](arr []T, val T) bool {
	for _, element := range arr {
		if element == val {
			return true
		}
	}

	return false
}

func SafeRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}

type SelectorFn[T any, R any] func(*T) *R

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

type MapFn[T any, R any] func(*T, *R)

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

func JoinArrays[T any](arr []T, element T) []T {
	return append(arr, element)
}
