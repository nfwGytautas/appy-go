package main

import (
	"github.com/nfwGytautas/appy"
	appy_default_drivers "github.com/nfwGytautas/appy/defaults"
)

func main() {
	options := appy.AppyOptions{
		Environment: appy.DefaultEnvironment(),
		Logger: appy.LoggerOptions{
			Provider: appy_default_drivers.DefaultLogger(),
			Name:     "Appy",
		},
		HTTP: &appy.HttpOptions{
			Provider: appy_default_drivers.DefaultHttpServer(),
			Address:  "127.0.0.1:8080",
			SSL:      nil, // HTTP
		},
		Sockets: &appy.WebsocketFactoryOptions{
			Provider: appy_default_drivers.DefaultWebsocketFactory(),
		},
	}

	// Create
	app, err := appy.New(options)
	if err != nil {
		panic(err)
	}

	// Add an endpoint handler
	app.Http().RootGroup().GET("/connect", func(c *appy.HttpContext) appy.HttpResult {
		socket := c.App.Sockets().Create(appy.WebsocketOptions{
			OnMessage: func(socket appy.Websocket, data []byte) {
				c.App.Logger.Debug("Received message: '%v'", string(data))
				socket.Send([]byte("Server: " + string(data)))
			},
		})

		// Don't need to return anything the socket will handle the result
		socket.Spin(c)
		return c.Nil()
	})

	// Run
	err = app.Run()
	if err != nil {
		panic(err)
	}
}
