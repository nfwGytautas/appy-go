package appy_hooks

type Slot func()

// Signal is a simple signal/slot implementation.
type Signal struct {
	callbacks []Slot
}

// NewSignal initializes a new Signal.
func NewSignal() *Signal {
	return &Signal{}
}

// Connect connects a slot to the signal.
func (s *Signal) Connect(slot Slot) {
	s.callbacks = append(s.callbacks, slot)
}

// Emit emits the signal.
func (s *Signal) Emit() {
	for _, callback := range s.callbacks {
		callback()
	}
}
