package driver_logger

import (
	"log"

	"github.com/nfwGytautas/appy"
)

const cyanSequence = "\u001B[36m"
const greenSequence = "\u001B[32m"
const yellowSequence = "\u001B[33m"
const redSequence = "\u001B[31m"
const resetSequence = "\u001B[0m"

type ConsoleLogger struct {
}

// A very simple console logging provider, mainly to be used for debugging or a quick setup
func ConsoleProvider() appy.Logger {
	return &ConsoleLogger{}
}

func (c *ConsoleLogger) Debug(fmt string, args ...interface{}) {
	log.Printf("[DEBUG] "+cyanSequence+fmt+resetSequence, args...)
}

func (c *ConsoleLogger) Info(fmt string, args ...interface{}) {
	log.Printf("[ INFO] "+greenSequence+fmt+resetSequence, args...)
}

func (c *ConsoleLogger) Warn(fmt string, args ...interface{}) {
	log.Printf("[ WARN] "+yellowSequence+fmt+resetSequence, args...)
}

func (c *ConsoleLogger) Error(fmt string, args ...interface{}) {
	log.Printf("[ERROR] "+redSequence+fmt+resetSequence, args...)
}
