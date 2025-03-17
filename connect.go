package appy

import (
	appy_jobs "github.com/nfwGytautas/appy-go/jobs"
	appy_tracker "github.com/nfwGytautas/appy-go/tracker"
)

// Options to pass when creating an appy
type AppyOptions struct {
	Environment EnvironmentSettings
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
}
