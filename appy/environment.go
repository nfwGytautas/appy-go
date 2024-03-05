package appy

// Various settings for the environment in appy
type EnvironmentSettings struct {
	// If true then appy instructs its providers to run in debug mode
	DebugMode bool

	// If true appy won't return nil handlers for services when request in such cases the user needs to handle this himself,
	// however its advised to leave this false and use Has* methods of appy, defaulted to false
	FailOnInvalidService bool
}

// DefaultEnvironment returns a new EnvironmentSettings with default values
func DefaultEnvironment() EnvironmentSettings {
	return EnvironmentSettings{
		DebugMode:            true,
		FailOnInvalidService: false,
	}
}
