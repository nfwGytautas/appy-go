package appy_http

import (
	"github.com/gin-gonic/gin"
	appy_logger "github.com/nfwGytautas/appy/logger"
)

// GIN powered http provider for appy
type ginHttpServer struct {
	engine  *gin.Engine
	options HttpOptions

	rootGroup *gin.RouterGroup
}

func (g *ginHttpServer) Run() error {
	if g.options.SSL != nil {
		appy_logger.Get().Debug("Running HTTPS")
		return g.engine.RunTLS(g.options.Address, g.options.SSL.CertFile, g.options.SSL.KeyFile)
	}

	appy_logger.Get().Debug("Running HTTP")
	return g.engine.Run(g.options.Address)
}

func (g *ginHttpServer) RootGroup() HttpEndpointGroup {
	return &ginHttpEndpointGroup{
		provider: g,
		group:    g.rootGroup,
	}
}
