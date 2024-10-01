package appy

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	appy_driver "github.com/nfwGytautas/appy-go/driver"
	appy_logger "github.com/nfwGytautas/appy-go/logger"
	appy_tracker "github.com/nfwGytautas/appy-go/tracker"
	appy_utils "github.com/nfwGytautas/appy-go/utils"
)

type HttpHandler func(c *gin.Context, context context.Context, transaction *appy_driver.Tx, tracker appy_tracker.Tracker)
type WsHandler func(c *gin.Context, context context.Context, transaction *appy_driver.Tx, tracker appy_tracker.Tracker) WsFn

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

func AppyHttpBootstrap(handler HttpHandler) gin.HandlerFunc {
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
			HTTP().HandleError(ctx, c, err)
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
		handler(c, ctx, tx, tracker)

		if statusHijacker.statusCode >= 400 {
			tx.Rollback()
			return
		}

		err = tx.Commit()
		if err != nil {
			tx.Rollback()
			HTTP().HandleError(ctx, c, err)
			return
		}
	}
}

func AppyHttpBootstrapConfig(handler HttpHandler, config BootstrapConfig) gin.HandlerFunc {
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
				HTTP().HandleError(ctx, c, err)
				return
			}
		} else {
			tx = nil
		}

		statusHijacker := statusHijacker{
			ResponseWriter: c.Writer,
			statusCode:     http.StatusOK,
			body:           nil,
		}

		defer statusHijacker.FlushHijack()

		c.Writer = &statusHijacker

		// Handler code
		handler(c, ctx, tx, tracker)

		if statusHijacker.statusCode >= 400 {
			tx.Rollback()
			return
		}

		if config.Transaction {
			err = tx.Commit()
			if err != nil {
				tx.Rollback()
				HTTP().HandleError(ctx, c, err)
				return
			}
		}
	}
}

func AppyWsBootstrap(handler WsHandler) gin.HandlerFunc {
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
			HTTP().HandleError(ctx, c, err)
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
		socketFn := handler(c, ctx, tx, tracker)

		if statusHijacker.statusCode >= 400 || socketFn == nil {
			tx.Rollback()
			return
		}

		err = tx.Commit()
		if err != nil {
			tx.Rollback()
			HTTP().HandleError(ctx, c, err)
			return
		}

		// Socket code
		err = socketFn(c, ctx)
		if err != nil {
			HTTP().HandleError(ctx, c, err)
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
