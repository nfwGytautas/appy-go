package appy_jobs

import (
	"sync"
	"time"
)

type JobScheduler struct {
	sync.RWMutex

	stop chan bool
	jobs []jobEntry
}

type jobEntry struct {
	job  Job
	tick time.Duration
}

func (n *JobScheduler) Add(options JobOptions) {
	n.Lock()
	defer n.Unlock()

	job := jobEntry{
		job:  options.Job,
		tick: options.Tick,
	}

	n.jobs = append(n.jobs, job)
}

func (n *JobScheduler) Start() {
	for _, job := range n.jobs {
		go func(job jobEntry) {
			for {
				select {
				case <-n.stop:
					return
				case <-time.After(job.tick):
					job.job()
				}
			}
		}(job)
	}
}
