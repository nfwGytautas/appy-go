package appy_http

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	appy_logger "github.com/nfwGytautas/appy/logger"
	appy_tracker "github.com/nfwGytautas/appy/tracker"
	appy_utils "github.com/nfwGytautas/appy/utils"
)

type ginHttpEndpointGroup struct {
	provider *ginHttpServer
	group    *gin.RouterGroup
	parent   *ginHttpEndpointGroup

	pre []HttpMiddleware
}

func (g *ginHttpEndpointGroup) Subgroup(path string) HttpEndpointGroup {
	return &ginHttpEndpointGroup{
		parent:   g,
		provider: g.provider,
		group:    g.group.Group(path),
	}
}

func (g *ginHttpEndpointGroup) Pre(middleware ...HttpMiddleware) {
	g.pre = append(g.pre, middleware...)
}

func (g *ginHttpEndpointGroup) StaticFile(path, file string) {
	g.group.StaticFile(path, file)
}

func (g *ginHttpEndpointGroup) StaticDir(path string, dir http.FileSystem) {
	g.group.StaticFS(path, dir)
}

func (g *ginHttpEndpointGroup) GET(path string, handler HttpHandler) {
	g.group.GET(path, func(c *gin.Context) {
		g.handle(c, handler)
	})
}

func (g *ginHttpEndpointGroup) POST(path string, handler HttpHandler) {
	g.group.POST(path, func(c *gin.Context) {
		g.handle(c, handler)
	})
}

func (g *ginHttpEndpointGroup) PATCH(path string, handler HttpHandler) {
	g.group.PATCH(path, func(c *gin.Context) {
		g.handle(c, handler)
	})
}

func (g *ginHttpEndpointGroup) PUT(path string, handler HttpHandler) {
	g.group.PUT(path, func(c *gin.Context) {
		g.handle(c, handler)
	})
}

func (g *ginHttpEndpointGroup) DELETE(path string, handler HttpHandler) {
	g.group.DELETE(path, func(c *gin.Context) {
		g.handle(c, handler)
	})
}

func (g *ginHttpEndpointGroup) appyCtx(c *gin.Context) HttpContext {
	return HttpContext{
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
		Context: context.TODO(),
	}
}

func (g *ginHttpEndpointGroup) handle(c *gin.Context, handler HttpHandler) {
	handlerName := appy_utils.ReflectFunctionName(handler)
	appy_logger.Get().Debug("Handling request: %v", handlerName)

	ctx := g.appyCtx(c)

	ctx.Tracker = appy_tracker.Get().OpenScope(handlerName)
	ctx.Tracker.SetRequest(c.Request)
	ctx.Transaction = appy_tracker.Get().OpenTransaction(ctx.Context, handlerName)

	defer func() {
		ctx.Transaction.Finish()
		appy_tracker.Get().Flush()
	}()

	err := g.runPreHandlerMiddleware(&ctx)
	if err != nil {
		g.handleResult(c, ctx.Tracker, ctx.Error(err))
		return
	}

	handlerRes := handler(&ctx)
	if handlerRes.IsFailed() {
		g.handleResult(c, ctx.Tracker, handlerRes)
		return
	}

	g.handleResult(c, ctx.Tracker, handlerRes)
}

func (g *ginHttpEndpointGroup) handleResult(c *gin.Context, tracker *appy_tracker.Scope, res HttpResult) {
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

		if res.StatusCode >= 500 {
			tracker.CaptureError(res.Error)
			appy_logger.Get().Debug("Error in handler: '%v', at: '%v'", res.Error.Error(), res.Tracker.At)
		}
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

func (g *ginHttpEndpointGroup) runPreHandlerMiddleware(ctx *HttpContext) error {
	if g.parent != nil {
		err := g.parent.runPreHandlerMiddleware(ctx)
		if err != nil {
			return err
		}
	}

	for _, pre := range g.pre {
		name := appy_utils.ReflectFunctionName(pre)
		ctx.Tracker.AddBreadcrumb("Pre middleware", name)

		err := pre(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
