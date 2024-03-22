package appy_default_drivers

import (
	"github.com/nfwGytautas/appy"
	appy_driver_http "github.com/nfwGytautas/appy/http"
	appy_driver_jobs "github.com/nfwGytautas/appy/jobs"
	appy_driver_logger "github.com/nfwGytautas/appy/logger"
	appy_driver_tracker "github.com/nfwGytautas/appy/tracker"
	appy_driver_websocket "github.com/nfwGytautas/appy/websocket"
)

func DefaultLogger() appy.LoggerImplementation {
	return appy_driver_logger.ConsoleProvider()
}

func DefaultHttpServer() appy.HttpServer {
	return appy_driver_http.Provider()
}

func DefaultJobScheduler() appy.JobScheduler {
	return appy_driver_jobs.NewScheduler()
}

func DefaultWebsocketFactory() appy.WebsocketFactory {
	return appy_driver_websocket.Factory()
}

func DefaultTracker(dsn string, sampleRate float32) appy.Tracker {
	return appy_driver_tracker.NewTracker(dsn, sampleRate)
}
