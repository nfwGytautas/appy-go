package appy_http

// Base type for http endpoints
type Endpoint struct {
}

// Options for endpoints
type EndpointOptions struct {
}

// Interface for endpoint input
type EndpointIf interface {
	Options() EndpointOptions
	Parse(pr *ParamReader)
	Validate() *ValidationResult
	Handle() EndpointResult
}

// Type to return from 'Validate' call
type ValidationResult struct {
	valid   bool
	message string
}

// Type to return from 'Handle' call
type EndpointResult struct {
}

func Valid() *ValidationResult {
	return &ValidationResult{
		valid:   true,
		message: "",
	}
}

func Invalid(reason string) *ValidationResult {
	return &ValidationResult{
		valid:   false,
		message: reason,
	}
}

func Ok() EndpointResult {
	return EndpointResult{}
}

func Error() EndpointResult {
	return EndpointResult{}
}
