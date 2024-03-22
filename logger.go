package appy

// Logger is the interface that wraps the basic logging methods
type LoggerImplementation interface {
	Debug(fmt string, args ...interface{})
	Info(fmt string, args ...interface{})
	Warn(fmt string, args ...interface{})
	Error(fmt string, args ...interface{})

	Initialize() error
	Flush()
}

// Options when creating a new logger
type LoggerOptions struct {
	Implementations []LoggerImplementation
}

// Logger is a struct for containing all the implementations
type Logger struct {
	implementations []LoggerImplementation
}

func newLoger() Logger {
	return Logger{
		implementations: make([]LoggerImplementation, 0),
	}
}

func (l *Logger) setImplementations(impls []LoggerImplementation) error {
	l.implementations = impls

	for _, impl := range l.implementations {
		err := impl.Initialize()
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *Logger) Debug(fmt string, args ...interface{}) {
	for _, impl := range l.implementations {
		impl.Debug(fmt, args...)
	}
}

func (l *Logger) Info(fmt string, args ...interface{}) {
	for _, impl := range l.implementations {
		impl.Info(fmt, args...)
	}
}

func (l *Logger) Warn(fmt string, args ...interface{}) {
	for _, impl := range l.implementations {
		impl.Warn(fmt, args...)
	}
}

func (l *Logger) Error(fmt string, args ...interface{}) {
	for _, impl := range l.implementations {
		impl.Error(fmt, args...)
	}
}

func (l *Logger) Initialize() error {
	for _, impl := range l.implementations {
		err := impl.Initialize()
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *Logger) Flush() {
	for _, impl := range l.implementations {
		impl.Flush()
	}
}

func (l *Logger) Attach(impl LoggerImplementation) {
	l.implementations = append(l.implementations, impl)
}
