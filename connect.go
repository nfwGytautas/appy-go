package appy

import (
	"os"
	"os/signal"
	"syscall"

	appy_config "github.com/nfwGytautas/appy-go/config"
	appy_jobs "github.com/nfwGytautas/appy-go/jobs"
	appy_logger "github.com/nfwGytautas/appy-go/logger"
	appy_tracker "github.com/nfwGytautas/appy-go/tracker"
)

// Options to pass when creating an appy
type AppyOptions struct {
	Environment EnvironmentSettings
	HTTP        *appy_config.HttpConfig
	Jobs        *appy_jobs.JobSchedulerOptions
	Tracker     *appy_tracker.TrackerOptions
}

func Initialize(options AppyOptions) {
	var err error

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

	if options.HTTP != nil {
		err = InitializeHTTP(*options.HTTP)
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

	defer appy_logger.Logger().Flush()

	if !appy_tracker.IsInitialized() {
		appy_tracker.Flush()
	}

	if HTTP() != nil {
		HTTP().Run()
	} else {
		appy_logger.Logger().Info("No HTTP server available, running till signal (Ctrl+C)")
		waitForSignal()
	}
}

func waitForSignal() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	<-sig
	appy_logger.Logger().Info("Received signal. Exiting gracefully...")
}
