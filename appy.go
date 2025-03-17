package appy

import (
	"os"
	"os/signal"
	"syscall"

	appy_http "github.com/nfwGytautas/appy-go/http"
	appy_jobs "github.com/nfwGytautas/appy-go/jobs"
	appy_logger "github.com/nfwGytautas/appy-go/logger"
	appy_tracker "github.com/nfwGytautas/appy-go/tracker"
	appy_utils "github.com/nfwGytautas/appy-go/utils"
)

type Appy struct {
	scheduler  *appy_jobs.JobScheduler
	httpServer *appy_http.Server

	opts AppOpts
}

type RegisterEndpointsFn func(app *Appy) error

type AppOpts struct {
	Version     string
	Controllers []Controller

	RegisterEndpoints RegisterEndpointsFn
}

// Create new appy app
func NewApp(opts AppOpts) *Appy {
	return &Appy{
		opts: opts,
	}
}

func (a *Appy) Initialize() error {
	appy_logger.Logger().Info("Initializing appy-go")

	// Set some environment variables
	os.Setenv("VERSION", a.opts.Version)

	err := appy_jobs.Initialize(appy_jobs.JobSchedulerOptions{})
	if err != nil {
		return err
	}

	// TODO: Remove global access
	a.scheduler = appy_jobs.Get()

	// Load environment variables if any
	if appy_utils.FsFileExists(".env") {
		appy_logger.Logger().Info("Loading environment variables from .env file")
		err := Environment().LoadFromFile(".env")
		if err != nil {
			return err
		}
	}

	// Controllers
	for _, controller := range a.opts.Controllers {
		controller.SetupJobs(a.scheduler)
	}

	// Register endpoints
	if a.opts.RegisterEndpoints != nil {
		err := a.opts.RegisterEndpoints(a)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *Appy) Run() {
	appy_logger.Logger().Info("Starting app, version: '%s'", a.opts.Version)

	// CLI Override
	if len(os.Args) > 1 {
		a.handleCLI()
		return
	}

	if appy_jobs.Get() != nil {
		go appy_jobs.Get().Start()
	}

	defer appy_logger.Logger().Flush()

	if !appy_tracker.IsInitialized() {
		appy_tracker.Flush()
	}

	if a.httpServer != nil {
		a.httpServer.Run()
	} else {
		appy_logger.Logger().Info("No HTTP server available, running till signal (Ctrl+C)")
		waitForSignal()
	}
}

func (a *Appy) SetHttpServer(server *appy_http.Server) {
	a.httpServer = server
}

func (a *Appy) handleCLI() {

}

func waitForSignal() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	<-sig
	appy_logger.Logger().Info("Received signal. Exiting gracefully...")
}
