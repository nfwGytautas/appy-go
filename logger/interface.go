package appy_logger

var logger Logger

func Initialize() error {
	logger = Logger{}
	return nil
}

func Get() *Logger {
	return &logger
}
