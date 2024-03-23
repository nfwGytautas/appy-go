package appy_logger

// Logger is the interface that wraps the basic logging methods
type Logger interface {
	Debug(fmt string, args ...interface{})
	Info(fmt string, args ...interface{})
	Warn(fmt string, args ...interface{})
	Error(fmt string, args ...interface{})

	Initialize() error
	Flush()
}

var logger Logger

func Initialize() error {
	logger = &loggerImplementation{}
	return nil
}

func Get() Logger {
	return logger
}
