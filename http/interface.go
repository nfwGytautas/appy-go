package appy_http

import (
	"context"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	appy_tracker "github.com/nfwGytautas/appy/tracker"
)

// Options when creating a new http server
type HttpOptions struct {
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

// HttpHandler is a function that handles an http request
type HttpHandler func(*HttpContext) HttpResult

// The HttpProvider interface is used to define the methods that are required for a http server, these are implemented by drivers
type HttpServer interface {
	Run() error

	// Get the root group for the http server i.e. /
	RootGroup() HttpEndpointGroup
}

// Utility struct for mapping errors to http responses
type HttpErrorMapper interface {
	Map(*HttpResult)
}

// HttpEndpointGroup is used to group http methods together
type HttpEndpointGroup interface {
	// Create a subgroup of the current group i.e. /api -> /api/v1
	Subgroup(string) HttpEndpointGroup

	// Attach pre-handle middleware
	Pre(...HttpMiddleware)

	// Attach post-handle middleware
	Post(...HttpMiddleware)

	StaticFile(string, string)
	StaticDir(string, http.FileSystem)

	GET(string, HttpHandler)
	POST(string, HttpHandler)
	PATCH(string, HttpHandler)
	PUT(string, HttpHandler)
	DELETE(string, HttpHandler)
}

// HeaderParser is used to parse headers from a http request
type HeaderParser interface {
	// Get a header by name
	ExpectSingleString(string) (string, error)
}

// QueryParameterParser is used to parse query parameters from a http request
type QueryParameterParser interface {
	// A utility function to get a page from a query parameter, same as calling GetInt("page")
	Page() int

	// Get a query paramter by name
	GetString(string) string

	// Get a query paramter by name and convert to int
	GetInt(string) int

	ExpectPage() (int, error)
	ExpectString(string) (string, error)
	ExpectInt(string) (int, error)
}

// PathParameterParser is used to parse path parameters from a http request
type PathParameterParser interface {
	// Get a path paramter by name and convert to int
	GetInt(string) int

	ExpectInt(string) (int, error)
}

// BodyParser is used to parse the body of a http request
type BodyParser interface {
	// Parse the body into a struct
	ParseSingle(any) error

	// Parse an array of structs from the body
	ParseArray(any) error
}

// A http handler context
type HttpContext struct {
	Header HeaderParser
	Query  QueryParameterParser
	Path   PathParameterParser
	Body   BodyParser

	Context context.Context

	Writer  http.ResponseWriter
	Request *http.Request

	Tracker     appy_tracker.TrackerScope
	Transaction appy_tracker.TrackerTransaction

	// Temporary storage to pass from middleware to handler
	tempStorage map[string]any
}

// A result of a http request
type HttpResult struct {
	StatusCode int
	Body       interface{}
	Error      error
	Tracker    HttpResultTrackerInfo

	failed bool
}

type HttpResultTrackerInfo struct {
	At string
}

// This is a middleware function that can be used to add functionality that runs before the main handler,
// if the returned result is not Nil it will be passed down the chain
type HttpMiddleware func(*HttpContext) error

var server HttpServer

func Initialize(options HttpOptions) error {
	srv := &ginHttpServer{}
	// if !app.Environment.DebugMode {
	// 	gin.SetMode(gin.ReleaseMode)
	// }

	srv.engine = gin.Default()
	srv.engine.Use(cors.Default())
	srv.options = options

	srv.rootGroup = srv.engine.Group("/")

	if srv.options.ErrorMapper == nil {
		srv.options.ErrorMapper = &defaultErrorMapper{}
	}

	server = srv

	return nil
}

func Get() HttpServer {
	return server
}
