package appy

import (
	"os"

	appy_jobs "github.com/nfwGytautas/appy-go/jobs"
	appy_logger "github.com/nfwGytautas/appy-go/logger"
)

type Appy struct {
	controllers []Controller
	version     string
	scheduler   *appy_jobs.JobScheduler
}

// Create new appy app
func NewApp(version string, controllers []Controller) *Appy {
	return &Appy{
		version:     version,
		controllers: controllers,
	}
}

func (a *Appy) Initialize() error {
	appy_logger.Logger().Info("Initializing appy-go")

	err := appy_jobs.Initialize(appy_jobs.JobSchedulerOptions{})
	if err != nil {
		return err
	}

	// TODO: Remove global access
	a.scheduler = appy_jobs.Get()

	for _, controller := range a.controllers {
		controller.SetupJobs(a.scheduler)
	}

	return nil
}

func (a *Appy) Run() {
	appy_logger.Logger().Info("Starting app, version: '%s'", a.version)

	// CLI Override
	if len(os.Args) > 1 {
		a.handleCLI()
		return
	}

	Takeover()
}

func (a *Appy) handleCLI() {

}
