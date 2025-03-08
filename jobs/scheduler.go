package appy_jobs

import (
	"sync"
	"time"

	appy_logger "github.com/nfwGytautas/appy-go/logger"
	appy_utils "github.com/nfwGytautas/appy-go/utils"
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
			prevFinished := true
			jobName := appy_utils.ReflectFunctionName(job.job)
			for {
				select {
				case <-n.stop:
					return
				case <-time.After(job.tick):
					if !prevFinished {
						// Skip
						appy_logger.Logger().Warn("Job overrun %s, skipping", jobName)
						continue
					}

					prevFinished = false
					start := time.Now()
					job.job()
					elapsed := time.Since(start)
					appy_logger.Logger().Info("Job %s took: '%v'", jobName, elapsed)
					prevFinished = true
				}
			}
		}(job)
	}
}
