package main

import (
	"github.com/nfwGytautas/appy"
	appy_driver_http "github.com/nfwGytautas/appy/http"
	appy_driver_logger "github.com/nfwGytautas/appy/logger"
	appy_driver_websocket "github.com/nfwGytautas/appy/websocket"
)

func main() {
	options := appy.AppyOptions{
		Environment: appy.DefaultEnvironment(),
		Logger: appy.LoggerOptions{
			Provider: appy_driver_logger.ConsoleProvider(),
			Name:     "Appy",
		},
		HTTP: &appy.HttpOptions{
			Provider: appy_driver_http.Provider(),
			Address:  "127.0.0.1:8080",
			SSL:      nil, // HTTP
		},
		Sockets: &appy.WebsocketFactoryOptions{
			Provider: appy_driver_websocket.Factory(),
		},
	}

	// Create
	app, err := appy.New(options)
	if err != nil {
		panic(err)
	}

	// Add an endpoint handler
	app.Http().RootGroup().GET("/connect", func(c appy.HttpContext) appy.HttpResult {
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
