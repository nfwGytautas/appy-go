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
