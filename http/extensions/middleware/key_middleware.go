package appy_middleware

import (
	"crypto/subtle"

	"github.com/gin-gonic/gin"
	appy_http "github.com/nfwGytautas/appy-go/http"
)

func ApiKeyMiddleware(apiKey string, failStatusCode int) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check header
		token := c.GetHeader("Authorization")
		if token == "" {
			c.Abort()
			appy_http.Get().HandleError(c.Request.Context(), c, ErrAuthorizationHeaderMissing)
			return
		}

		if subtle.ConstantTimeCompare([]byte(token), []byte(apiKey)) == 0 {
			c.Abort()
			appy_http.Get().HandleError(c.Request.Context(), c, ErrApiKeysDontMatch)
			return
		}

		c.Next()
	}
}
