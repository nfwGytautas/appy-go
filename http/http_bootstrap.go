package appy_http

import (
	"context"
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
	appy_driver "github.com/nfwGytautas/appy-go/driver"
	appy_logger "github.com/nfwGytautas/appy-go/logger"
	appy_tracker "github.com/nfwGytautas/appy-go/tracker"
	appy_utils "github.com/nfwGytautas/appy-go/utils"
)

type RequestContext struct {
	server *Server
	c      *gin.Context
	Ctx    context.Context

	Tx      *appy_driver.Tx
	Tracker appy_tracker.Tracker

	// Internal
	status int
	result any
	err    error

	postCommits []PostCommitJob
}

type HttpHandler func(r *RequestContext)
type WsHandler func(r *RequestContext) WsFn

type PostCommitJob func()

type WsFn func(c *gin.Context, context context.Context) error

type statusHijacker struct {
	gin.ResponseWriter
	statusCode int
	body       []byte
}

type BootstrapConfig struct {
	Tracked     bool
	Transaction bool
}

func (s *Server) Handle(handler HttpHandler) gin.HandlerFunc {
	return s.HandleConfig(handler, BootstrapConfig{
		Tracked:     true,
		Transaction: true,
	})
}

func (s *Server) HandleConfig(handler HttpHandler, config BootstrapConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		var tracker appy_tracker.Tracker
		var tx *appy_driver.Tx

		// Debug
		currentFunctionName := appy_utils.ReflectFunctionName(handler)
		appy_logger.Logger().Debug("Running: '%v'", currentFunctionName)

		ctx := c.Request.Context()

		if config.Tracked {
			ctx, tracker = appy_tracker.Begin(c.Request.Context(), currentFunctionName)
			defer tracker.Finish()

			tracker.SetRequest(c.Request)
		} else {
			tracker = appy_tracker.BeginDummy()
		}

		// DB Transaction setup
		if config.Transaction {
			tx, err = appy_driver.StartTransaction()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
				return
			}
		} else {
			tx = nil
		}

		// Handler code
		r := &RequestContext{
			c:       c,
			server:  s,
			Ctx:     ctx,
			Tx:      tx,
			Tracker: tracker,
			status:  0,
			result:  nil,
		}

		handler(r)

		if config.Transaction {
			if r.status >= 400 || r.err != nil {
				tx.Rollback()
				r.setGinStatus()
				return
			}

			err = tx.Commit()
			if err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
				return
			}
		}

		// Post commit jobs
		for _, job := range r.postCommits {
			go job()
		}

		r.setGinStatus()
	}
}

func (s *Server) Ws(handler WsHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Debug
		currentFunctionName := appy_utils.ReflectFunctionName(handler)
		appy_logger.Logger().Debug("Running: '%v'", currentFunctionName)

		// Tracker setup
		ctx, tracker := appy_tracker.Begin(c.Request.Context(), currentFunctionName)
		defer tracker.Finish()

		tracker.SetRequest(c.Request)

		// DB Transaction setup
		tx, err := appy_driver.StartTransaction()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
			return
		}

		statusHijacker := statusHijacker{
			ResponseWriter: c.Writer,
			statusCode:     http.StatusOK,
			body:           nil,
		}

		defer statusHijacker.FlushHijack()

		c.Writer = &statusHijacker

		// Handler code
		r := &RequestContext{
			c:       c,
			server:  s,
			Ctx:     ctx,
			Tx:      tx,
			Tracker: tracker,
			status:  0,
			result:  nil,
		}

		socketFn := handler(r)

		if statusHijacker.statusCode >= 400 || socketFn == nil || r.err != nil {
			tx.Rollback()
			r.setGinStatus()
			return
		}

		err = tx.Commit()
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
			return
		}

		// Socket code
		err = socketFn(c, ctx)
		if err != nil {
			r.setGinStatus()
			return
		}
	}
}

func (sh *statusHijacker) WriteHeader(code int) {
	sh.statusCode = code
	sh.body = nil
}

func (sh *statusHijacker) Write(data []byte) (int, error) {
	sh.body = data
	return len(data), nil
}

func (sh *statusHijacker) FlushHijack() {
	sh.ResponseWriter.WriteHeader(sh.statusCode)

	if sh.body != nil {
		sh.ResponseWriter.Write(sh.body)
		return
	}
}

func (r *RequestContext) Status(status int) {
	r.status = status
}

func (r *RequestContext) Result(status int, result any) {
	r.status = status
	r.result = result
}

func (r *RequestContext) ParamChain() *ParamChain {
	return &ParamChain{
		context:      r.c,
		currentError: nil,
	}
}

func (r *RequestContext) Error(err error) {
	_, file, line, _ := runtime.Caller(1)
	appy_logger.Logger().Error("Error while handling request: '%v:%v', error: '%v'", file, line, err)

	r.err = err
}

func (r *RequestContext) Redirect(status int, location string) {
	r.c.Redirect(status, location)
}

func (r *RequestContext) StoreMultipartFile(fileKey string, destination string) (string, error) {
	return storeMultipartFile(r.c, fileKey, destination)
}

func (r *RequestContext) PostForm(key string) string {
	return r.c.PostForm(key)
}

func (r *RequestContext) setGinStatus() {
	if r.err != nil {
		statusCode, body := r.server.options.ErrorMapper.Map(r.Ctx, r.err)

		if body != nil {
			r.c.JSON(statusCode, body)
		} else {
			r.c.Status(statusCode)
		}

		return
	}

	if r.result != nil {
		r.c.JSON(r.status, r.result)
		return
	}

	r.c.Status(r.status)
}

func (r *RequestContext) PostCommitJob(job PostCommitJob) {
	r.postCommits = append(r.postCommits, job)
}
