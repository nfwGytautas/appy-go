package main

import (
	"time"

	"github.com/nfwGytautas/appy"
	appy_driver_jobs "github.com/nfwGytautas/appy/jobs"
	appy_driver_logger "github.com/nfwGytautas/appy/logger"
)

func main() {
	options := appy.AppyOptions{
		Environment: appy.DefaultEnvironment(),
		Logger: appy.LoggerOptions{
			Name:     "Appy",
			Provider: appy_driver_logger.ConsoleProvider(),
		},
		Jobs: &appy.JobSchedulerOptions{
			Provider: appy_driver_jobs.NewScheduler(),
			PoolTick: 1 * time.Second,
		},
	}

	// Create
	app, err := appy.New(options)
	if err != nil {
		panic(err)
	}

	// A a persistent job that is executed every second,
	// NOTE: that if the tick is smaller than the pool tick, the job will not be executed every specified tick
	app.Jobs().Add(appy.JobOptions{
		Job: func(app appy.Appy) {
			app.Logger.Debug("I am a persistent job")

			// You can queue one-off jobs from within a persistent job, make sure you do it on a different thread tho
			go app.Jobs().Add(appy.JobOptions{
				Job: func(app appy.Appy) {
					app.Logger.Debug("I am a one-off job queued from a persistent job")
				},
				Persistent: false,
			})
		},
		Tick:       5 * time.Second,
		Pooled:     true,
		Persistent: true,
	})

	app.Jobs().Add(appy.JobOptions{
		Job: func(app appy.Appy) {
			app.Logger.Debug("I am a throttled job")
		},
		Tick:       1 * time.Millisecond,
		Pooled:     true,
		Persistent: true,
	})

	app.Jobs().Add(appy.JobOptions{
		Job: func(app appy.Appy) {
			app.Logger.Debug("I am a one-off job")
		},
		Persistent: false,
	})

	// Run
	err = app.Run()
	if err != nil {
		panic(err)
	}
}
