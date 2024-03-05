package appy

// Logger is a simple logging interface
type Logger struct {
	Name     string
	provider LoggerProvider
}

// Options when creating a new logger
type LoggerOptions struct {
	Provider LoggerProvider
	Name     string
}

// LoggerProvider is the interface that wraps the basic logging methods
type LoggerProvider interface {
	Debug(fmt string, args ...interface{})
	Info(fmt string, args ...interface{})
	Warn(fmt string, args ...interface{})
	Error(fmt string, args ...interface{})
}

func (l *Logger) Debug(fmt string, args ...interface{}) {
	l.provider.Debug(fmt, args...)
}

func (l *Logger) Info(fmt string, args ...interface{}) {
	l.provider.Info(fmt, args...)
}

func (l *Logger) Warn(fmt string, args ...interface{}) {
	l.provider.Warn(fmt, args...)
}

func (l *Logger) Error(fmt string, args ...interface{}) {
	l.provider.Error(fmt, args...)
}
