package appy

// Logger is the interface that wraps the basic logging methods
type Logger interface {
	Debug(fmt string, args ...interface{})
	Info(fmt string, args ...interface{})
	Warn(fmt string, args ...interface{})
	Error(fmt string, args ...interface{})
}

// Options when creating a new logger
type LoggerOptions struct {
	Provider Logger
	Name     string
}
