package appy

import (
	"github.com/gin-gonic/gin"
)

// MiddlewareProvider is an interface that should be implemented by all middleware providers.
type MiddlewareProvider interface {
	Provide() gin.HandlerFunc
}

// Controller is an interface that should be implemented by all controllers.
type Controller interface {
	SetupJobs(*App)
}
