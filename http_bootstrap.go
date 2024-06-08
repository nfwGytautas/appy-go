package appy

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	appy_driver "github.com/nfwGytautas/appy-go/driver"
	appy_http "github.com/nfwGytautas/appy-go/http"
	appy_logger "github.com/nfwGytautas/appy-go/logger"
	appy_tracker "github.com/nfwGytautas/appy-go/tracker"
	appy_utils "github.com/nfwGytautas/appy-go/utils"
)

type HttpHandler func(c *gin.Context, context context.Context, transaction *appy_driver.Tx, tracker appy_tracker.Tracker)

type statusHijacker struct {
	gin.ResponseWriter
	statusCode int
	body       []byte
}

func AppyHttpBootstrap(handler HttpHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Debug
		currentFunctionName := appy_utils.ReflectFunctionName(handler)
		appy_logger.Get().Debug("Running: '%v'", currentFunctionName)

		// Tracker setup
		ctx, tracker := appy_tracker.Begin(c.Request.Context(), currentFunctionName)
		defer tracker.Finish()

		tracker.SetRequest(c.Request)

		// DB Transaction setup
		tx, err := appy_driver.StartTransaction()
		if err != nil {
			appy_http.Get().HandleError(ctx, c, err)
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
			appy_http.Get().HandleError(ctx, c, err)
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
