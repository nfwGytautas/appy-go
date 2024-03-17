package appy_driver_http

import (
	"github.com/gin-gonic/gin"
	"github.com/nfwGytautas/appy"
)

type ginHttpEndpointGroup struct {
	provider *ginHttpServer
	group    *gin.RouterGroup
}

func (g *ginHttpEndpointGroup) Subgroup(path string) appy.HttpEndpointGroup {
	return &ginHttpEndpointGroup{
		provider: g.provider,
		group:    g.group.Group(path),
	}
}

func (g *ginHttpEndpointGroup) Use(middleware appy.HttpMiddleware) {
	g.group.Use(func(c *gin.Context) {
		res := middleware(g.appyCtx(c))
		if g.handleResult(c, res) {
			return
		}

		c.Next()
	})
}

func (g *ginHttpEndpointGroup) GET(path string, handler appy.HttpHandler) {
	g.group.GET(path, func(c *gin.Context) {
		g.handle(c, handler)
	})
}

func (g *ginHttpEndpointGroup) POST(path string, handler appy.HttpHandler) {
	g.group.POST(path, func(c *gin.Context) {
		g.handle(c, handler)
	})
}

func (g *ginHttpEndpointGroup) PUT(path string, handler appy.HttpHandler) {
	g.group.PUT(path, func(c *gin.Context) {
		g.handle(c, handler)
	})
}

func (g *ginHttpEndpointGroup) DELETE(path string, handler appy.HttpHandler) {
	g.group.DELETE(path, func(c *gin.Context) {
		g.handle(c, handler)
	})
}

func (g *ginHttpEndpointGroup) appyCtx(c *gin.Context) appy.HttpContext {
	return appy.HttpContext{
		App: g.provider.app,
		Query: &ginQueryParser{
			ctx: c,
		},
		Path: &ginPathParser{
			ctx: c,
		},
		Writer:  c.Writer,
		Request: c.Request,
	}
}

func (g *ginHttpEndpointGroup) handle(c *gin.Context, handler appy.HttpHandler) {
	res := handler(g.appyCtx(c))
	g.handleResult(c, res)
}

func (g *ginHttpEndpointGroup) handleResult(c *gin.Context, res appy.HttpResult) bool {
	failed := false

	// Unexpected error
	if res.HasError() {
		// Try and map the error from the error map
		res = g.provider.options.ErrorMapper.Map(res.Error)
		failed = true
	}

	// Write the response
	if res.Body != nil {
		c.JSON(
			res.StatusCode,
			res.Body,
		)
	} else {
		c.Status(res.StatusCode)
	}

	return failed
}
