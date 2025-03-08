package appy

import (
	"github.com/gin-gonic/gin"
	appy_jobs "github.com/nfwGytautas/appy-go/jobs"
)

// MiddlewareProvider is an interface that should be implemented by all middleware providers.
type MiddlewareProvider interface {
	Provide() gin.HandlerFunc
}

// Controller is an interface that should be implemented by all controllers.
type Controller interface {
	SetupJobs(*appy_jobs.JobScheduler)
}
