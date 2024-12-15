package appy_http

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ParamReader is a utility struct for getting query parameters
type ParamReader struct {
	context *gin.Context

	err error
}

// Global validator instance for body validation
var validate = validator.New()

func newParamReader(c *gin.Context) *ParamReader {
	return &ParamReader{
		context: c,
		err:     nil,
	}
}
