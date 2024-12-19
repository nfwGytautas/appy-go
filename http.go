package appy

import (
	"context"
	"runtime"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	appy_config "github.com/nfwGytautas/appy-go/config"
	appy_logger "github.com/nfwGytautas/appy-go/logger"
)

var server *Server

func InitializeHTTP(options appy_config.HttpConfig) error {
	e := gin.Default()
	e.Use(cors.Default())

	server = &Server{
		engine:  e,
		options: options,
	}

	return nil
}

func HTTP() *Server {
	return server
}

type Server struct {
	engine  *gin.Engine
	options appy_config.HttpConfig
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

func (s *Server) HandleError(ctx context.Context, c *gin.Context, err error) {
	_, file, line, _ := runtime.Caller(1)
	appy_logger.Logger().Error("Error while handling request: '%v:%v', error: '%v'", file, line, err)

	statusCode, body := s.options.ErrorMapper.Map(ctx, err)

	if body != nil {
		c.JSON(statusCode, body)
	} else {
		c.Status(statusCode)
	}
}
