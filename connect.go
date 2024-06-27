package appy

import (
	"os"
	"os/signal"
	"syscall"

	appy_http "github.com/nfwGytautas/appy-go/http"
	appy_jobs "github.com/nfwGytautas/appy-go/jobs"
	appy_logger "github.com/nfwGytautas/appy-go/logger"
	appy_tracker "github.com/nfwGytautas/appy-go/tracker"
	appy_websockets "github.com/nfwGytautas/appy-go/websocket"
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

	if options.Jobs != nil {
		err = appy_jobs.Initialize(*options.Jobs)
		if err != nil {
			panic(err)
		}
	}

	if options.Tracker != nil {
		err = appy_tracker.Initialize(*options.Tracker)
		if err != nil {
			panic(err)
		}
	}

	err = appy_websockets.Initialize()
	if err != nil {
		panic(err)
	}

	if options.HTTP != nil {
		err = appy_http.Initialize(*options.HTTP)
		if err != nil {
			panic(err)
		}
	}
}

func Takeover() {
	start()
}

func start() {
	if appy_jobs.Get() != nil {
		go appy_jobs.Get().Start()
	}

	defer appy_logger.Get().Flush()

	if !appy_tracker.IsInitialized() {
		appy_tracker.Flush()
	}

	if appy_http.Get() != nil {
		appy_http.Get().Run()
	} else {
		appy_logger.Get().Info("No HTTP server available, running till signal (Ctrl+C)")
		waitForSignal()
	}
}

func waitForSignal() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	<-sig
	appy_logger.Get().Info("Received signal. Exiting gracefully...")
}
