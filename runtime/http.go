package appy_runtime

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	appy_config "github.com/nfwGytautas/appy-go/config"
	appy_http "github.com/nfwGytautas/appy-go/http"
	appy_logger "github.com/nfwGytautas/appy-go/logger"
)

type EndpointProxyFactoryFn func() appy_http.EndpointIf

type HttpEngine struct {
	config *appy_config.HttpConfig
	engine *gin.Engine
}

func createHttpEngine(config *appy_config.HttpConfig) *HttpEngine {
	e := gin.Default()
	e.Use(cors.Default())

	return &HttpEngine{
		config: config,
		engine: e,
	}
}

func (he *HttpEngine) run() error {
	if he.config.SSL != nil {
		appy_logger.Logger().Debug("Running HTTPS")
		return he.engine.RunTLS(he.config.Address, he.config.SSL.CertFile, he.config.SSL.KeyFile)
	}

	appy_logger.Logger().Debug("Running HTTP")
	return he.engine.Run(he.config.Address)
}

func (he *HttpEngine) RegisterEndpoint(method string, path string, endpoint EndpointProxyFactoryFn) {
	appy_logger.Logger().Debug("Registering endpoint %s %s", method, path)
}
