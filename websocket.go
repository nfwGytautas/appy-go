package appy

// A factory for websocket objects
type WebsocketFactory interface {
	// Initialize the factory
	Initialize(*Appy, WebsocketFactoryOptions) error

	// Create a new empty websocket object, without spinning it
	Create(WebsocketOptions) Websocket
}

// Options for a websocket factory
type WebsocketFactoryOptions struct {
	Provider WebsocketFactory
}

// Websocket object
type Websocket interface {
	// Start the websocket
	Spin(HttpContext) error

	// Send a message to websocket
	Send([]byte)

	// Close the websocket
	Close() error
}

// Options to pass when creating a websocket
type WebsocketOptions struct {
	OnClose   OnCloseCallback
	OnMessage OnMessageCallback

	UserData any
}

type OnCloseCallback func(Websocket)
type OnMessageCallback func(Websocket, []byte)
