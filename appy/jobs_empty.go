package appy

type nilJobScheduler struct{}

func (n *nilJobScheduler) Initialize(*Appy, JobSchedulerOptions) error {
	return nil
}

func (n *nilJobScheduler) Add(options JobOptions) {
}

func (n *nilJobScheduler) Start() {
}

func (n *nilJobScheduler) Stop() {
}
