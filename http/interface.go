package appy_http

import (
	"context"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Utility struct for mapping errors to http responses
type HttpErrorMapper interface {
	Map(context.Context, error) (int, any)
}

var server *Server

func Initialize(options HttpOptions) error {
	e := gin.Default()
	e.Use(cors.Default())

	server = &Server{
		engine:  e,
		options: options,
	}

	return nil
}

func Get() *Server {
	return server
}
