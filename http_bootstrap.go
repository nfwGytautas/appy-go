package appy

import (
	"context"

	"github.com/gin-gonic/gin"
	appy_driver "github.com/nfwGytautas/appy-go/driver"
	appy_http "github.com/nfwGytautas/appy-go/http"
	appy_logger "github.com/nfwGytautas/appy-go/logger"
	appy_tracker "github.com/nfwGytautas/appy-go/tracker"
	appy_utils "github.com/nfwGytautas/appy-go/utils"
)

type HttpHandler func(c *gin.Context, context context.Context, transaction *appy_driver.Tx, tracker appy_tracker.Tracker)

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
		defer tx.CommitOrRollback()

		// Handler code
		handler(c, ctx, tx, tracker)
	}
}
