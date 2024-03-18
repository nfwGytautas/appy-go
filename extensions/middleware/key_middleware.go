package appy_middleware

import (
	"crypto/subtle"

	"github.com/nfwGytautas/appy"
)

func ApiKeyMiddleware(apiKey string, failStatusCode int) appy.HttpMiddleware {
	return func(c *appy.HttpContext) appy.HttpResult {
		// Check header
		token, err := c.Header.ExpectSingleString("Authorization")
		if err != nil {
			return c.Error(err)
		}

		if subtle.ConstantTimeCompare([]byte(token), []byte(apiKey)) == 0 {
			return c.Fail(failStatusCode, "Invalid token")
		}

		return c.Nil()
	}
}
