package appy

import "net/http"

// The HttpProvider interface is used to define the methods that are required for a http server, these are implemented by drivers
type HttpServer interface {
	Initialize(*Appy, HttpOptions) error
	Run() error

	// Get the root group for the http server i.e. /
	RootGroup() HttpEndpointGroup
}

// Utility struct for mapping errors to http responses
type HttpErrorMapper interface {
	Map(error) HttpResult
}

// Options when creating a new http server
type HttpOptions struct {
	Provider HttpServer

	// Error mapper
	ErrorMapper HttpErrorMapper

	// The address to bind the server to
	Address string

	// SSL settings for HTTPS, runs on HTTP if nil
	SSL *SSLSettings
}

// SSLSettings is used to define the settings for a https server
type SSLSettings struct {
	CertFile string
	KeyFile  string
}

// HttpEndpointGroup is used to group http methods together
type HttpEndpointGroup interface {
	// Create a subgroup of the current group i.e. /api -> /api/v1
	Subgroup(string) HttpEndpointGroup

	Use(HttpMiddleware)

	GET(string, HttpHandler)
	POST(string, HttpHandler)
	PUT(string, HttpHandler)
	DELETE(string, HttpHandler)
}

// HttpHandler is a function that handles an http request
type HttpHandler func(c HttpContext) HttpResult

// A http handler context
type HttpContext struct {
	App   *Appy
	Query QueryParameterParser
	Path  PathParameterParser

	Writer  http.ResponseWriter
	Request *http.Request
}

// A result of a http request
type HttpResult struct {
	StatusCode int
	Body       interface{}
	Error      error

	failed bool
}

// QueryParameterParser is used to parse query parameters from a http request
type QueryParameterParser interface{}

// PathParameterParser is used to parse path parameters from a http request
type PathParameterParser interface{}

// This is a middleware function that can be used to add functionality that runs before the main handler,
// if the returned result is not Nil it will be passed down the chain
type HttpMiddleware func(c HttpContext) HttpResult

func (c *HttpContext) Nil() HttpResult {
	return HttpResult{
		failed: false,
	}
}

func (c *HttpContext) Ok(statusCode int, body interface{}) HttpResult {
	return HttpResult{
		StatusCode: statusCode,
		Body:       body,
		failed:     false,
	}
}

func (c *HttpContext) Fail(statusCode int, body interface{}) HttpResult {
	return HttpResult{
		StatusCode: statusCode,
		Body:       body,
		failed:     true,
	}
}

func (c *HttpContext) Error(err error) HttpResult {
	return HttpResult{
		Error:  err,
		failed: true,
	}
}

func (hr HttpResult) IsFailed() bool {
	return hr.failed || hr.Error != nil
}

func (hr HttpResult) HasError() bool {
	return hr.Error != nil
}
