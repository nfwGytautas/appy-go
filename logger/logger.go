package appy_logger

import (
	"log"
)

type Logger struct {
}

func (c *Logger) Debug(fmt string, args ...interface{}) {
	log.Printf("[DEBUG] "+cyanSequence+fmt+resetSequence, args...)
}

func (c *Logger) Info(fmt string, args ...interface{}) {
	log.Printf("[ INFO] "+greenSequence+fmt+resetSequence, args...)
}

func (c *Logger) Warn(fmt string, args ...interface{}) {
	log.Printf("[ WARN] "+yellowSequence+fmt+resetSequence, args...)
}

func (c *Logger) Error(fmt string, args ...interface{}) {
	log.Printf("[ERROR] "+redSequence+fmt+resetSequence, args...)
}

func (c *Logger) Initialize() error {
	return nil
}

func (c *Logger) Flush() {
	// Nothing to do here
}
