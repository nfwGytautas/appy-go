package appy_driver_jobs

import (
	"errors"
	"sync"
	"time"

	"github.com/nfwGytautas/appy"
)

type jobScheduler struct {
	sync.RWMutex

	app     *appy.Appy
	stop    chan bool
	started bool

	oneOffJobs    []job
	nonPooledJobs []job
	pooledJobs    []job

	poolTick time.Duration
}

type job struct {
	job          appy.Job
	currentDelta int
	tick         time.Duration
}

func NewScheduler() appy.JobScheduler {
	return &jobScheduler{
		started: false,
		stop:    make(chan bool),
	}
}

func (n *jobScheduler) Initialize(app *appy.Appy, options appy.JobSchedulerOptions) error {
	if options.PoolTick < 1 {
		return errors.New("pool tick must be greater than 0")
	}

	n.app = app
	n.poolTick = options.PoolTick
	return nil
}

func (n *jobScheduler) Add(options appy.JobOptions) {
	if options.Persistent && n.started {
		n.app.Logger.Error("Cannot add a persistent job to a started scheduler")
		return
	}

	if options.Tick < n.poolTick {
		n.app.Logger.Warn("Job tick is smaller than pool tick, job will be throttled, job: '%v'", appy.ReflectFunctionName(options.Job))
	}

	n.Lock()
	defer n.Unlock()

	job := job{
		job:          options.Job,
		currentDelta: 0,
		tick:         options.Tick,
	}

	if options.Instant {
		go func() {
			time.Sleep(job.tick)
			job.job(*n.app)
		}()
		return
	}

	if !options.Persistent {
		n.oneOffJobs = append(n.oneOffJobs, job)
		return
	}

	if !options.Pooled {
		n.nonPooledJobs = append(n.nonPooledJobs, job)
	} else {
		n.pooledJobs = append(n.pooledJobs, job)
	}
}

func (n *jobScheduler) Start() {
	ticker := time.NewTicker(n.poolTick)

	for {
		select {
		case <-ticker.C:
			n.handleTick()
		case <-n.stop:
			return
		}
	}
}

func (n *jobScheduler) Stop() {
	n.stop <- true
}

func (n *jobScheduler) handleTick() {
	n.Lock()
	defer n.Unlock()

	for i := 0; i < len(n.oneOffJobs); i++ {
		job := n.oneOffJobs[i]
		go job.job(*n.app)
	}

	for i := 0; i < len(n.nonPooledJobs); i++ {
		job := &n.nonPooledJobs[i]
		job.currentDelta += int(n.poolTick)
		if job.currentDelta >= int(job.tick) {
			go job.job(*n.app)
			job.currentDelta = 0
		}
	}

	for i := 0; i < len(n.pooledJobs); i++ {
		job := &n.pooledJobs[i]
		job.currentDelta += int(n.poolTick)
		if job.currentDelta >= int(job.tick) {
			job.job(*n.app)
			job.currentDelta = 0
		}
	}

	// Remove one off jobs
	n.oneOffJobs = nil
}
