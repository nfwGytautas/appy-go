package appy

import "github.com/gin-gonic/gin"

// MiddlewareProvider is an interface that should be implemented by all middleware providers.
type MiddlewareProvider interface {
	Provide() gin.HandlerFunc
}
