package appy_middleware

import (
	"crypto/subtle"

	appy_http "github.com/nfwGytautas/appy/http"
)

func ApiKeyMiddleware(apiKey string, failStatusCode int) appy_http.HttpMiddleware {
	return func(c *appy_http.HttpContext) appy_http.HttpResult {
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
