package appy_driver_http

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/nfwGytautas/appy"
)

// GIN powered http provider for appy
type ginHttpServer struct {
	engine  *gin.Engine
	app     *appy.Appy
	options appy.HttpOptions

	rootGroup *gin.RouterGroup
}

type ginQueryParser struct {
	ctx *gin.Context
}

type ginPathParser struct {
	ctx *gin.Context
}

// Create a new appy HttpProvider
func Provider() appy.HttpServer {
	return &ginHttpServer{}
}

func (g *ginHttpServer) Initialize(app *appy.Appy, options appy.HttpOptions) error {
	if !app.Environment.DebugMode {
		gin.SetMode(gin.ReleaseMode)
	}

	g.engine = gin.Default()
	g.engine.Use(cors.Default())
	g.app = app
	g.options = options

	g.rootGroup = g.engine.Group("/")

	if g.options.ErrorMapper == nil {
		g.options.ErrorMapper = &defaultErrorMapper{}
	}

	return nil
}

func (g *ginHttpServer) Run() error {
	if g.options.SSL != nil {
		g.app.Logger.Debug("Running HTTPS")
		return g.engine.RunTLS(g.options.Address, g.options.SSL.CertFile, g.options.SSL.KeyFile)
	}

	g.app.Logger.Debug("Running HTTP")
	return g.engine.Run(g.options.Address)
}

func (g *ginHttpServer) RootGroup() appy.HttpEndpointGroup {
	return &ginHttpEndpointGroup{
		provider: g,
		group:    g.rootGroup,
	}
}
