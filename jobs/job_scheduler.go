package appy_jobs

import (
	"sync"
	"time"

	appy_logger "github.com/nfwGytautas/appy/logger"
	appy_utils "github.com/nfwGytautas/appy/utils"
)

type JobScheduler struct {
	sync.RWMutex

	stop    chan bool
	started bool

	oneOffJobs    []jobEntry
	nonPooledJobs []jobEntry
	pooledJobs    []jobEntry

	poolTick time.Duration
}

type jobEntry struct {
	job          Job
	currentDelta int
	tick         time.Duration
}

func (n *JobScheduler) Add(options JobOptions) {
	n.Lock()
	defer n.Unlock()

	if options.Tick < n.poolTick {
		appy_logger.Get().Warn("Job tick is smaller than pool tick, job will be throttled, job: '%v'", appy_utils.ReflectFunctionName(options.Job))
	}

	job := jobEntry{
		job:          options.Job,
		currentDelta: 0,
		tick:         options.Tick,
	}

	switch options.Type {
	case JobTypeInstant:
		go func() {
			time.Sleep(options.Tick)
			options.Job()
		}()
	case JobTypeOneOff:
		n.oneOffJobs = append(n.oneOffJobs, job)
	case JobTypePersistent:
		if n.started {
			appy_logger.Get().Error("Cannot add a persistent job to a started scheduler")
			return
		}
		n.nonPooledJobs = append(n.nonPooledJobs, job)
	case JobTypePooled:
		if n.started {
			appy_logger.Get().Error("Cannot add a persistent job to a started scheduler")
			return
		}
		n.pooledJobs = append(n.pooledJobs, job)
	}
}

func (n *JobScheduler) Start() {
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

func (n *JobScheduler) Stop() {
	n.stop <- true
}

func (n *JobScheduler) handleTick() {
	n.Lock()
	defer n.Unlock()

	for i := 0; i < len(n.oneOffJobs); i++ {
		job := n.oneOffJobs[i]
		go job.job()
	}

	for i := 0; i < len(n.nonPooledJobs); i++ {
		job := &n.nonPooledJobs[i]
		job.currentDelta += int(n.poolTick)
		if job.currentDelta >= int(job.tick) {
			go job.job()
			job.currentDelta = 0
		}
	}

	for i := 0; i < len(n.pooledJobs); i++ {
		job := &n.pooledJobs[i]
		job.currentDelta += int(n.poolTick)
		if job.currentDelta >= int(job.tick) {
			job.job()
			job.currentDelta = 0
		}
	}

	// Remove one off jobs
	n.oneOffJobs = nil
}
