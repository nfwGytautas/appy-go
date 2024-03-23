package appy_logger

import (
	"log"
)

type loggerImplementation struct {
}

func (c *loggerImplementation) Debug(fmt string, args ...interface{}) {
	log.Printf("[DEBUG] "+cyanSequence+fmt+resetSequence, args...)
}

func (c *loggerImplementation) Info(fmt string, args ...interface{}) {
	log.Printf("[ INFO] "+greenSequence+fmt+resetSequence, args...)
}

func (c *loggerImplementation) Warn(fmt string, args ...interface{}) {
	log.Printf("[ WARN] "+yellowSequence+fmt+resetSequence, args...)
}

func (c *loggerImplementation) Error(fmt string, args ...interface{}) {
	log.Printf("[ERROR] "+redSequence+fmt+resetSequence, args...)
}

func (c *loggerImplementation) Initialize() error {
	return nil
}

func (c *loggerImplementation) Flush() {
	// Nothing to do here
}
