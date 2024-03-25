package appy

import (
	"os"
	"os/signal"
	"syscall"

	appy_http "github.com/nfwGytautas/appy/http"
	appy_jobs "github.com/nfwGytautas/appy/jobs"
	appy_logger "github.com/nfwGytautas/appy/logger"
	appy_tracker "github.com/nfwGytautas/appy/tracker"
	appy_websockets "github.com/nfwGytautas/appy/websocket"
)

// Options to pass when creating an appy
type AppyOptions struct {
	Environment EnvironmentSettings
	HTTP        *appy_http.HttpOptions
	Jobs        *appy_jobs.JobSchedulerOptions
	Tracker     *appy_tracker.TrackerOptions
}

func Initialize(options AppyOptions) {
	err := appy_logger.Initialize()
	if err != nil {
		panic(err)
	}

	err = appy_jobs.Initialize(*options.Jobs)
	if err != nil {
		panic(err)
	}

	err = appy_tracker.Initialize(*options.Tracker)
	if err != nil {
		panic(err)
	}

	err = appy_websockets.Initialize()
	if err != nil {
		panic(err)
	}

	err = appy_http.Initialize(*options.HTTP)
	if err != nil {
		panic(err)
	}
}

func Takeover() {
	start()
}

func start() {
	go appy_jobs.Get().Start()

	defer appy_logger.Get().Flush()
	appy_tracker.Flush()

	appy_http.Get().Run()

	// appy_logger.Get().Info("No HTTP server available, running till signal (Ctrl+C)")
	// waitForSignal()
}

func waitForSignal() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	<-sig
	appy_logger.Get().Info("Received signal. Exiting gracefully...")
}
