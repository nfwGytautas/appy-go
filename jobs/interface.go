package appy_jobs

import (
	"time"

	appy_errors "github.com/nfwGytautas/appy/errors"
)

// Possible job types
const (
	JobTypePooled     = 1
	JobTypePersistent = 2
	JobTypeOneOff     = 3
	JobTypeInstant    = 4
)

// Alias for job types
type JobType int

// A single job, just a simple function
type Job func()

// Job scheduler used to add and execute jobs in the application
type JobScheduler interface {
	// Add a new job to the scheduler of the specified duration
	Add(JobOptions)

	// Start job execution
	Start()

	// Stop job execution
	Stop()
}

// Options for creating a scheduler
type JobSchedulerOptions struct {
	// The duration to wait before executing another cycle of job pool checks
	PoolTick time.Duration
}

// Options to describe a single job
type JobOptions struct {
	// The job to execute
	Job Job

	// The duration to wait before executing the job
	Tick time.Duration

	// Type of the job
	Type JobType
}

var scheduler JobScheduler

func Initialize(options JobSchedulerOptions) error {
	scheduler = &jobScheduler{
		stop:     make(chan bool),
		poolTick: options.PoolTick,
	}

	if options.PoolTick < 1 {
		return appy_errors.ErrInvalidPoolTick
	}

	return nil
}

func Get() JobScheduler {
	return scheduler
}
