package appy

import (
	"context"
	"runtime"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	appy_logger "github.com/nfwGytautas/appy-go/logger"
)

// Utility struct for mapping errors to http responses
type HttpErrorMapper interface {
	Map(context.Context, error) (int, any)
}

var server *Server

func InitializeHTTP(options HttpOptions) error {
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

// SSLSettings is used to define the settings for a https server
type SSLSettings struct {
	CertFile string
	KeyFile  string
}

// Options when creating a new http server
type HttpOptions struct {
	Mapper HttpErrorMapper

	// The address to bind the server to
	Address string

	// SSL settings for HTTPS, runs on HTTP if nil
	SSL *SSLSettings
}

type Server struct {
	engine  *gin.Engine
	options HttpOptions
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

	statusCode, body := s.options.Mapper.Map(ctx, err)

	if body != nil {
		c.JSON(statusCode, body)
	} else {
		c.Status(statusCode)
	}
}
