package appy_websockets

import "github.com/gorilla/websocket"

type WebsocketFactory struct {
	upgrader websocket.Upgrader
	id       uint64
}

// Create a new empty websocket object, without spinning it
func (wf *WebsocketFactory) Create(options WebsocketOptions) *Websocket {
	wf.id += 1

	return &Websocket{
		factory:   wf,
		id:        wf.id,
		chanClose: make(chan bool),
		chanSend:  make(chan []byte),
		options:   options,
	}
}
