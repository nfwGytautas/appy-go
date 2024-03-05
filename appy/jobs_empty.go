package appy

type nilJobSchedulerProvider struct{}

func (n *nilJobSchedulerProvider) Initialize(*Appy, JobSchedulerOptions) error {
	return nil
}

func (n *nilJobSchedulerProvider) Add(options JobOptions) {
}

func (n *nilJobSchedulerProvider) Start() {
}

func (n *nilJobSchedulerProvider) Stop() {
}
