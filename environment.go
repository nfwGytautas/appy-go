package appy

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Various settings for the environment in appy
type EnvironmentSettings struct {
	// If true then appy instructs its providers to run in debug mode
	DebugMode bool
}

var environmentSettings EnvironmentSettings = EnvironmentSettings{
	DebugMode: false,
}

// Get the current environment settings for appy
func Environment() *EnvironmentSettings {
	return &environmentSettings
}

// LoadFromFile loads the environment settings from a file
func (es *EnvironmentSettings) LoadFromFile(file string) error {
	err := godotenv.Load(file)
	if err != nil {
		return err
	}

	return nil
}

func (es *EnvironmentSettings) GetValue(key string) (string, error) {
	val := os.Getenv(key)

	if val == "" {
		return "", errors.New(fmt.Sprintf("environment variable '%s' not set", key))
	}

	return val, nil
}
