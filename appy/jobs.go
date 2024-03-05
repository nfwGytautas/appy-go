package appy

import "time"

type JobScheduler struct {
	provider JobSchedulerProvider
}

type JobSchedulerOptions struct {
	Provider JobSchedulerProvider

	// The duration to wait before executing another cycle of job pool checks
	PoolTick time.Duration
}

type JobSchedulerProvider interface {
	Initialize(*Appy, JobSchedulerOptions) error

	// Add a new job to the scheduler of the specified duration
	Add(JobOptions)

	// Start job execution
	Start()

	// Stop job execution
	Stop()
}

type JobOptions struct {
	// The job to execute
	Job Job

	// The duration to wait before executing the job
	Tick time.Duration

	// If the job should be pooled, if true the job will be executed in its own go routine, if not it will be scheduled to run on a pool routine
	Pooled bool

	// If the job should be persistent, if true the job will be executed periodically until the end of the program or false if it is a one off job
	Persistent bool
}

type Job func(Appy)

func (js *JobScheduler) Start() {
	js.provider.Start()
}

func (js *JobScheduler) Add(job JobOptions) {
	js.provider.Add(job)
}
