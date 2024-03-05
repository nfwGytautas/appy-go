package driver_gin

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/nfwGytautas/appy"
)

// GIN powered http provider for appy
type ginHttpProvider struct {
	engine *gin.Engine
	app    *appy.Appy

	address   string
	rootGroup *gin.RouterGroup
	ssl       *appy.SSLSettings
	mapper    appy.HttpErrorMapper
}

type ginQueryParser struct {
	ctx *gin.Context
}

type ginPathParser struct {
	ctx *gin.Context
}

// Create a new appy HttpProvider
func Provider() appy.HttpProvider {
	return &ginHttpProvider{}
}

func (g *ginHttpProvider) Initialize(app *appy.Appy, options appy.HttpOptions) error {
	if !app.Environment.DebugMode {
		gin.SetMode(gin.ReleaseMode)
	}

	g.engine = gin.Default()
	g.engine.Use(cors.Default())
	g.address = options.Address
	g.ssl = options.SSL
	g.app = app
	g.mapper = options.ErrorMapper

	g.rootGroup = g.engine.Group("/")

	if g.mapper == nil {
		g.mapper = &defaultErrorMapper{}
	}

	return nil
}

func (g *ginHttpProvider) Run() error {
	if g.ssl != nil {
		g.app.Logger.Debug("Running HTTPS")
		return g.engine.RunTLS(g.address, g.ssl.CertFile, g.ssl.KeyFile)
	}

	g.app.Logger.Debug("Running HTTP")
	return g.engine.Run(g.address)
}

func (g *ginHttpProvider) RootGroup() appy.HttpEndpointGroup {
	return &ginHttpEndpointGroup{
		provider: g,
		group:    g.rootGroup,
	}
}
