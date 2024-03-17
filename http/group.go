package appy_driver_http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nfwGytautas/appy"
)

type ginHttpEndpointGroup struct {
	provider *ginHttpServer
	group    *gin.RouterGroup

	pre  []appy.HttpMiddleware
	post []appy.HttpMiddleware
}

func (g *ginHttpEndpointGroup) Subgroup(path string) appy.HttpEndpointGroup {
	return &ginHttpEndpointGroup{
		provider: g.provider,
		group:    g.group.Group(path),
	}
}

func (g *ginHttpEndpointGroup) Pre(middleware appy.HttpMiddleware) {
	g.pre = append(g.pre, middleware)
}

func (g *ginHttpEndpointGroup) Post(middleware appy.HttpMiddleware) {
	g.post = append(g.post, middleware)
}

func (g *ginHttpEndpointGroup) StaticFile(path, file string) {
	g.group.StaticFile(path, file)
}

func (g *ginHttpEndpointGroup) StaticDir(path string, dir http.FileSystem) {
	g.group.StaticFS(path, dir)
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

func (g *ginHttpEndpointGroup) PATCH(path string, handler appy.HttpHandler) {
	g.group.PATCH(path, func(c *gin.Context) {
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
		Header: &ginHeaderParser{
			ctx: c,
		},
		Query: &ginQueryParser{
			ctx: c,
		},
		Path: &ginPathParser{
			ctx: c,
		},
		Body: &ginBodyParser{
			ctx: c,
		},
		Writer:  c.Writer,
		Request: c.Request,
	}
}

func (g *ginHttpEndpointGroup) handle(c *gin.Context, handler appy.HttpHandler) {
	for _, pre := range g.pre {
		res := pre(g.appyCtx(c))
		if res.HasError() {
			g.handleResult(c, res)
			return
		}
	}

	res := handler(g.appyCtx(c))

	for _, post := range g.post {
		postRes := post(g.appyCtx(c))
		if postRes.HasError() {
			g.handleResult(c, res)
			return
		}
	}

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
