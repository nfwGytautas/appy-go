package appy_websocket

import (
	"bytes"
	"net/http"

	"github.com/gorilla/websocket"
	appy_logger "github.com/nfwGytautas/appy-go/logger"
)

// Character constants
var (
	cWebsocketNewline = []byte{'\n'}
	cWebsocketSpace   = []byte{' '}
)

// Timeout for client to send ready in seconds
const cWebsocketClientReadyTimeout = 10

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

// Upgrader for upgrading http requests to websocket connections
type websocketFactory struct {
	upgrader websocket.Upgrader
	id       uint64
}

var upgrader websocketFactory = websocketFactory{
	upgrader: websocket.Upgrader{
		// Check origin will check the cross region source
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	},
}

func WebsocketFactory() *websocketFactory {
	return &upgrader
}

// Create a new empty websocket object, without spinning it
func (wu *websocketFactory) Create(options WebsocketOptions) *Websocket {
	wu.id += 1

	return &Websocket{
		upgrader:  wu,
		id:        wu.id,
		chanClose: make(chan bool),
		chanSend:  make(chan []byte),
		options:   options,
	}
}

type Websocket struct {
	upgrader *websocketFactory
	ws       *websocket.Conn

	id uint64

	chanClose chan bool
	chanSend  chan []byte
	userClose bool

	options WebsocketOptions
}

// Start the websocket
func (ws *Websocket) Spin(writer http.ResponseWriter, request *http.Request) error {
	var err error
	ws.ws, err = ws.upgrader.upgrader.Upgrade(writer, request, nil)
	if err != nil {
		appy_logger.Logger().Error("Failed to spin websocket '%v'", err)
		return err
	}

	defer ws.Close()

	ws.userClose = false

	// Now run read/write process
	go ws.writerProcess()
	go ws.readerProcess()

	<-ws.chanClose
	appy_logger.Logger().Debug("Ending socket spin for '%v'", ws.id)

	err = ws.Close()
	if err != nil {
		appy_logger.Logger().Error("Failed to close websocket '%v'", err)
		return err
	}

	return nil
}

// Send a message to websocket
func (ws *Websocket) Send(message []byte) {
	appy_logger.Logger().Debug("Sending message '%v' to '%v'", string(message), ws.id)
	ws.chanSend <- message
}

// Close the websocket
func (ws *Websocket) Close() error {
	if ws.userClose {
		// Already closed by the user
		return nil
	}

	return ws.ws.Close()
}

func (ws *Websocket) writerProcess() {
	defer func() {
		ws.chanClose <- true
	}()

	// Infinite loop
	for {
		select {
		case message, ok := <-ws.chanSend:
			appy_logger.Logger().Debug("[%v] sending '%s'", ws.id, string(message))
			// c.conn.SetWriteDeadline(time.Now().Add(writeWait))

			if !ok {
				ws.ws.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := ws.ws.NextWriter(websocket.TextMessage)
			if err != nil {
				appy_logger.Logger().Error("%v", err)
				return
			}

			appy_logger.Logger().Debug("[%v] writing", ws.id)
			w.Write(message)
			w.Write(cWebsocketNewline)

			// Add queued chat messages to the current websocket message.
			n := len(ws.chanSend)
			appy_logger.Logger().Debug("[%v] queue backlog: '%v'", ws.id, n)
			for i := 0; i < n; i++ {
				w.Write(<-ws.chanSend)
				w.Write(cWebsocketNewline)
			}

			err = w.Close()
			if err != nil {
				appy_logger.Logger().Error("%v", err)
				return
			}

		case <-ws.chanClose:
			return
		}
	}
}

func (ws *Websocket) readerProcess() {
	defer func() {
		ws.chanClose <- true
	}()

	// Infinite loop
	for {
		_, message, err := ws.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
				websocket.CloseNormalClosure,
			) {
				ws.userClose = true
				appy_logger.Logger().Error("Error while reading socket: %v", err)
			}
			break
		}

		message = bytes.TrimSpace(bytes.Replace(message, cWebsocketNewline, cWebsocketSpace, -1))

		if ws.options.OnMessage != nil {
			ws.options.OnMessage(ws, message)
		}
	}

	if ws.options.OnClose != nil {
		ws.options.OnClose(ws)
	}
}
