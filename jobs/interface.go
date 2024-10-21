package appy_jobs

import (
	"time"
)

// A single job, just a simple function
type Job func()

// Options for creating a scheduler
type JobSchedulerOptions struct {
}

// Options to describe a single job
type JobOptions struct {
	// The job to execute
	Job Job

	// The duration to wait before executing the job
	Tick time.Duration

	// If specified the job will be executed using cron scheduler
	Cron string
}

var scheduler *JobScheduler

func Initialize(options JobSchedulerOptions) error {
	scheduler = &JobScheduler{
		stop: make(chan bool),
	}

	return nil
}

func Get() *JobScheduler {
	return scheduler
}
