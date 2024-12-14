package appy_logger

import "log"

// Color constants
const cLoggerCyanSequence = "\u001B[36m"
const cLoggerGreenSequence = "\u001B[32m"
const cLoggerYellowSequence = "\u001B[33m"
const cLoggerRedSequence = "\u001B[31m"
const cLoggerResetSequence = "\u001B[0m"

// Logger used to log messages
type logger struct {
}

var loggerInstance logger = logger{}

// Logger returns the logger instance
func Logger() *logger {
	return &loggerInstance
}

func (c *logger) Debug(fmt string, args ...interface{}) {
	log.Printf("[DEBUG] "+cLoggerCyanSequence+fmt+cLoggerResetSequence, args...)
}

func (c *logger) Info(fmt string, args ...interface{}) {
	log.Printf("[ INFO] "+cLoggerGreenSequence+fmt+cLoggerResetSequence, args...)
}

func (c *logger) Warn(fmt string, args ...interface{}) {
	log.Printf("[ WARN] "+cLoggerYellowSequence+fmt+cLoggerResetSequence, args...)
}

func (c *logger) Error(fmt string, args ...interface{}) {
	log.Printf("[ERROR] "+cLoggerRedSequence+fmt+cLoggerResetSequence, args...)
}

func (c *logger) Initialize() error {
	return nil
}

func (c *logger) Flush() {
	// Nothing to do here
}
