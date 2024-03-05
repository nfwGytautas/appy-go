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
	http *HttpServer

	// Scheduler for timed jobs to run in the background
	jobs *JobScheduler
}

// Options to pass when creating an appy
type AppyOptions struct {
	Environment EnvironmentSettings
	Logger      LoggerOptions
	HTTP        *HttpOptions
	Jobs        *JobSchedulerOptions
}

// New creates a new instance of the appy application
func New(options AppyOptions) (*Appy, error) {
	app := &Appy{}

	app.Environment = options.Environment

	// Logger
	app.Logger = Logger{
		Name:     options.Logger.Name,
		provider: options.Logger.Provider,
	}

	// Http
	if options.HTTP != nil {
		app.http = &HttpServer{
			provider: options.HTTP.Provider,
		}

		err := app.http.provider.Initialize(app, *options.HTTP)
		if err != nil {
			return nil, err
		}
	}

	// Jobs
	if options.Jobs != nil {
		app.jobs = &JobScheduler{
			provider: options.Jobs.Provider,
		}

		err := app.jobs.provider.Initialize(app, *options.Jobs)
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
		return a.http.provider.Run()
	}

	// No http server check other conditions for running
	a.Logger.Info("No HTTP server available, running till signal (Ctrl+C)")
	a.waitForSignal()

	return nil
}

func (a *Appy) HasHttp() bool {
	return a.http != nil
}

func (a *Appy) Http() *HttpServer {
	if a.http == nil && !a.Environment.FailOnInvalidService {
		a.Logger.Warn("No http server available, returning nil provider")
		return &HttpServer{
			provider: &nilHttpProvider{},
		}
	}

	return a.http
}

func (a *Appy) HasJobs() bool {
	return a.jobs != nil
}

func (a *Appy) Jobs() *JobScheduler {
	if a.jobs == nil && !a.Environment.FailOnInvalidService {
		a.Logger.Warn("No job scheduler available, returning nil provider")
		return &JobScheduler{
			provider: &nilJobSchedulerProvider{},
		}
	}

	return a.jobs
}

func (a *Appy) waitForSignal() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	<-sig
	a.Logger.Info("Received signal. Exiting gracefully...")
}
