package appy_driver_http

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nfwGytautas/appy"
)

type ginHttpEndpointGroup struct {
	provider *ginHttpServer
	group    *gin.RouterGroup
	parent   *ginHttpEndpointGroup

	pre  []appy.HttpMiddleware
	post []appy.HttpMiddleware
}

func (g *ginHttpEndpointGroup) Subgroup(path string) appy.HttpEndpointGroup {
	return &ginHttpEndpointGroup{
		parent:   g,
		provider: g.provider,
		group:    g.group.Group(path),
	}
}

func (g *ginHttpEndpointGroup) Pre(middleware ...appy.HttpMiddleware) {
	g.pre = append(g.pre, middleware...)
}

func (g *ginHttpEndpointGroup) Post(middleware ...appy.HttpMiddleware) {
	g.post = append(g.post, middleware...)
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
	g.provider.app.Logger.Debug("Handling request: %v", appy.ReflectFunctionName(handler))

	ctx := g.appyCtx(c)

	res := g.runPreHandlerMiddleware(&ctx)
	if res.IsFailed() {
		g.handleResult(c, res)
		return
	}

	handlerRes := handler(&ctx)
	if handlerRes.IsFailed() {
		g.handleResult(c, handlerRes)
		return
	}

	res = g.runPostHandlerMiddleware(&ctx)
	if res.IsFailed() {
		g.handleResult(c, res)
		return
	}

	g.handleResult(c, handlerRes)
}

func (g *ginHttpEndpointGroup) handleResult(c *gin.Context, res appy.HttpResult) {
	// Unexpected error
	if res.HasError() {
		// Try and map the error from the error map
		g.provider.options.ErrorMapper.Map(&res)

		c.JSON(
			res.StatusCode,
			gin.H{
				"body":  res.Body,
				"error": strings.Split(res.Error.Error(), "\n"),
				"debug": res.Tracker,
			},
		)

		g.provider.app.Logger.Debug("Error in handler: '%v', at: '%v'", res.Error.Error(), res.Tracker.At)

		return
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
}

func (g *ginHttpEndpointGroup) runPreHandlerMiddleware(ctx *appy.HttpContext) appy.HttpResult {
	if g.parent != nil {
		res := g.parent.runPreHandlerMiddleware(ctx)
		if res.HasError() {
			return res
		}
	}

	for _, pre := range g.pre {
		res := pre(ctx)
		if res.HasError() {
			return res
		}
	}

	return ctx.Nil()
}

func (g *ginHttpEndpointGroup) runPostHandlerMiddleware(ctx *appy.HttpContext) appy.HttpResult {
	for _, post := range g.post {
		res := post(ctx)
		if res.HasError() {
			return res
		}
	}

	if g.parent != nil {
		res := g.parent.runPostHandlerMiddleware(ctx)
		if res.HasError() {
			return res
		}
	}

	return ctx.Nil()
}
