package appy_http

import (
	"context"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	appy_logger "github.com/nfwGytautas/appy-go/logger"
)

// SSL settings
type SSLSettings struct {
	CertFile string
	KeyFile  string
}

// HTTP server config
type HttpConfig struct {
	Address string
	SSL     *SSLSettings

	ErrorMapper HttpErrorMapper
}

type HttpErrorMapper interface {
	Map(context.Context, error) (int, any)
}

type Server struct {
	engine  *gin.Engine
	options *HttpConfig
}

func InitializeHTTP(options *HttpConfig) (*Server, error) {
	e := gin.Default()
	e.Use(cors.Default())

	return &Server{
		engine:  e,
		options: options,
	}, nil
}

func (s *Server) Run() error {
	if s.options.SSL != nil {
		appy_logger.Logger().Debug("Running HTTPS")
		return s.engine.RunTLS(s.options.Address, s.options.SSL.CertFile, s.options.SSL.KeyFile)
	}

	appy_logger.Logger().Debug("Running HTTP")
	return s.engine.Run(s.options.Address)
}

func (s *Server) Root() *gin.RouterGroup {
	return s.engine.Group("/")
}
