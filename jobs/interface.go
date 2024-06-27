package appy_jobs

import (
	"time"

	appy_errors "github.com/nfwGytautas/appy-go/errors"
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

var scheduler *JobScheduler

func Initialize(options JobSchedulerOptions) error {
	scheduler = &JobScheduler{
		stop:     make(chan bool),
		poolTick: options.PoolTick,
	}

	if options.PoolTick < 1 {
		return appy_errors.ErrInvalidPoolTick
	}

	return nil
}

func Get() *JobScheduler {
	return scheduler
}
