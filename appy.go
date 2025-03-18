package appy

import (
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	appy_firebase "github.com/nfwGytautas/appy-go/firebase"
	appy_http "github.com/nfwGytautas/appy-go/http"
	appy_jobs "github.com/nfwGytautas/appy-go/jobs"
	appy_logger "github.com/nfwGytautas/appy-go/logger"
	appy_mail "github.com/nfwGytautas/appy-go/mail"
	appy_tracker "github.com/nfwGytautas/appy-go/tracker"
	appy_utils "github.com/nfwGytautas/appy-go/utils"
)

type App struct {
	scheduler  *appy_jobs.JobScheduler
	httpServer *appy_http.Server

	opts AppOpts
}

type ControllersSetupFn func(app *App) error
type EndpointsSetupFn func(app *App) error
type DatabaseSetupFn func(app *App) error
type JobsSetupFn func(app *App) error

type ConfigureFn func() *ModuleOpts

type AppOpts struct {
	Version string

	Controllers ControllersSetupFn
	Endpoints   EndpointsSetupFn
	Database    DatabaseSetupFn
	Jobs        JobsSetupFn

	Configure ConfigureFn
}

type ModuleOpts struct {
	Scheduler *appy_jobs.JobSchedulerOptions
	Tracker   *appy_tracker.TrackerOptions
	Mail      *appy_mail.MailerLiteOptions
	Firebase  *appy_firebase.FirebaseServicesOptions
}

// Create new appy app
func Go(opts AppOpts) {
	app := &App{
		opts: opts,
	}

	err := app.Initialize()
	if err != nil {
		log.Fatal(err)
	}

	app.Run()
}

func (a *App) Initialize() error {
	if a.opts.Configure == nil {
		return errors.New("configure hook not set")
	}

	appy_logger.Logger().Info("Initializing appy-go")

	// Set some environment variables
	os.Setenv("VERSION", a.opts.Version)

	// Load environment variables if any
	if appy_utils.FsFileExists(".env") {
		appy_logger.Logger().Info("Loading environment variables from .env file")
		err := Environment().LoadFromFile(".env")
		if err != nil {
			return err
		}
	}

	moduleSettings := a.opts.Configure()

	// Initialize modules
	if moduleSettings.Scheduler != nil {
		err := appy_jobs.Initialize(*moduleSettings.Scheduler)
		if err != nil {
			return err
		}
	}

	if moduleSettings.Tracker != nil {
		err := appy_tracker.Initialize(*moduleSettings.Tracker)
		if err != nil {
			return err
		}
	}

	if moduleSettings.Mail != nil {
		err := appy_mail.Mailerlite().Configure(*moduleSettings.Mail)
		if err != nil {
			return err
		}
	}

	if moduleSettings.Firebase != nil {
		err := appy_firebase.Firebase().Configure(*moduleSettings.Firebase)
		if err != nil {
			return err
		}
	}

	// Register endpoints
	if a.opts.Endpoints != nil {
		appy_logger.Logger().Debug("Running endpoint hook")
		err := a.opts.Endpoints(a)
		if err != nil {
			return err
		}
	}

	if a.opts.Database != nil {
		appy_logger.Logger().Debug("Running database hook")
		err := a.opts.Database(a)
		if err != nil {
			return err
		}
	}

	if a.opts.Controllers != nil {
		appy_logger.Logger().Debug("Running controllers hook")
		err := a.opts.Controllers(a)
		if err != nil {
			return err
		}
	}

	if a.opts.Jobs != nil {
		appy_logger.Logger().Debug("Running jobs hook")
		err := a.opts.Jobs(a)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) Run() {
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

func (a *App) SetHttpServer(server *appy_http.Server) {
	a.httpServer = server
}

func (a *App) handleCLI() {

}

func waitForSignal() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	<-sig
	appy_logger.Logger().Info("Received signal. Exiting gracefully...")
}
