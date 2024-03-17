package appy

import (
	"os"
	"os/signal"
	"syscall"
)

// Appy is the main struct for the appy package, here you will find
// all the services and methods to run the appy application
type Appy struct {
	// Environment settings for appy
	Environment EnvironmentSettings

	// Logger for appy
	Logger Logger

	// Server for appy, if this is null no endpoints will be available
	http HttpServer

	// Scheduler for timed jobs to run in the background
	jobs JobScheduler

	// If there is a need for websocket support a factory needs to be specified
	sockets WebsocketFactory
}

// Options to pass when creating an appy
type AppyOptions struct {
	Environment EnvironmentSettings
	Logger      LoggerOptions
	HTTP        *HttpOptions
	Jobs        *JobSchedulerOptions
	Sockets     *WebsocketFactoryOptions
}

// New creates a new instance of the appy application
func New(options AppyOptions) (*Appy, error) {
	app := &Appy{}

	app.Environment = options.Environment

	// Logger
	if options.Logger.Provider == nil {
		return nil, ErrNoLogger
	}

	app.Logger = options.Logger.Provider

	// Http
	if options.HTTP != nil {
		app.http = options.HTTP.Provider

		err := app.http.Initialize(app, *options.HTTP)
		if err != nil {
			return nil, err
		}
	}

	// Jobs
	if options.Jobs != nil {
		app.jobs = options.Jobs.Provider

		err := app.jobs.Initialize(app, *options.Jobs)
		if err != nil {
			return nil, err
		}
	}

	// Sockets
	if options.Sockets != nil {
		app.sockets = options.Sockets.Provider

		err := app.sockets.Initialize(app, *options.Sockets)
		if err != nil {
			return nil, err
		}
	}

	return app, nil
}

// Run the appy application, this will return an error if something goes wrong, otherwise this is a blocking call
func (a *Appy) Run() error {
	a.Logger.Debug("Starting appy application")

	// Start jobs
	if a.jobs != nil {
		a.Logger.Debug("Starting job scheduler")
		go a.jobs.Start()
	}

	if a.http != nil {
		a.Logger.Debug("Starting http server")
		return a.http.Run()
	}

	// No http server check other conditions for running
	a.Logger.Info("No HTTP server available, running till signal (Ctrl+C)")
	a.waitForSignal()

	return nil
}

func (a *Appy) HasHttp() bool {
	return a.http != nil
}

func (a *Appy) Http() HttpServer {
	if a.http == nil && !a.Environment.FailOnInvalidService {
		a.Logger.Warn("No http server available, returning nil provider")
		return &nilHttpServer{}
	}

	return a.http
}

func (a *Appy) HasJobs() bool {
	return a.jobs != nil
}

func (a *Appy) Jobs() JobScheduler {
	if a.jobs == nil && !a.Environment.FailOnInvalidService {
		a.Logger.Warn("No job scheduler available, returning nil provider")
		return &nilJobScheduler{}
	}

	return a.jobs
}

func (a *Appy) HasSockets() bool {
	return a.sockets != nil
}

func (a *Appy) Sockets() WebsocketFactory {
	if a.sockets == nil && !a.Environment.FailOnInvalidService {
		a.Logger.Warn("No websocket factory available, returning nil provider")
		return &nilWebsocketFactory{}
	}

	return a.sockets
}

func (a *Appy) waitForSignal() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	<-sig
	a.Logger.Info("Received signal. Exiting gracefully...")
}
