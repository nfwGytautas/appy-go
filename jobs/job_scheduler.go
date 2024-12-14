package appy_jobs

import (
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

type JobScheduler struct {
	sync.RWMutex

	stop chan bool
	jobs []jobEntry

	cronScheduler *cron.Cron // Lazy loading, only initialized when a cron job is added
}

type jobEntry struct {
	job  Job
	tick time.Duration
}

func (n *JobScheduler) Add(options JobOptions) {
	n.Lock()
	defer n.Unlock()

	if options.Cron != "" {
		if n.cronScheduler == nil {
			n.cronScheduler = cron.New(cron.WithLocation(time.UTC))
		}

		n.cronScheduler.AddFunc(options.Cron, options.Job)

		return
	}

	job := jobEntry{
		job:  options.Job,
		tick: options.Tick,
	}

	n.jobs = append(n.jobs, job)
}

func (n *JobScheduler) Start() {
	if n.cronScheduler != nil {
		n.cronScheduler.Start()
	}

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
