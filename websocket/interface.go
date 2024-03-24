package appy_websockets

import (
	"net/http"

	"github.com/gorilla/websocket"
)

// Callbacks for websocket close
type OnCloseCallback func(*Websocket)

// Callback for websocket messages
type OnMessageCallback func(*Websocket, []byte)

// Options to pass when creating a websocket
type WebsocketOptions struct {
	OnClose   OnCloseCallback
	OnMessage OnMessageCallback

	UserData any
}

var factory WebsocketFactory

func Initialize() error {
	factory = WebsocketFactory{
		upgrader: websocket.Upgrader{
			// Check origin will check the cross region source
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
	return nil
}

func Get() *WebsocketFactory {
	return &factory
}
