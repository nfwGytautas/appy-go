package driver_websocket

import (
	"bytes"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/nfwGytautas/appy"
)

type websocketFactory struct {
	upgrader websocket.Upgrader
	app      *appy.Appy
	id       uint64
}

type socket struct {
	factory *websocketFactory
	ws      *websocket.Conn

	id uint64

	chanClose chan bool
	chanSend  chan []byte

	options appy.WebsocketOptions
}

func Factory() appy.WebsocketFactory {
	return &websocketFactory{
		upgrader: websocket.Upgrader{
			// Check origin will check the cross region source
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (w *websocketFactory) Initialize(app *appy.Appy, options appy.WebsocketFactoryOptions) error {
	w.app = app
	return nil
}

func (w *websocketFactory) Create(options appy.WebsocketOptions) appy.Websocket {
	w.id += 1

	return &socket{
		factory:   w,
		id:        w.id,
		chanClose: make(chan bool),
		chanSend:  make(chan []byte),
		options:   options,
	}
}

func (ws *socket) Spin(c appy.HttpContext) error {
	var err error

	ws.ws, err = ws.factory.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		ws.factory.app.Logger.Error("Failed to spin websocket '%v'", err)
		return err
	}

	defer ws.Close()

	// Now run read/write process
	go ws.writerProcess()
	go ws.readerProcess()

	<-ws.chanClose
	ws.factory.app.Logger.Debug("Ending socket spin for '%v'", ws.id)

	err = ws.Close()
	if err != nil {
		ws.factory.app.Logger.Error("Failed to close websocket '%v'", err)
	}

	return err
}

func (ws *socket) Send(message []byte) {
	ws.chanSend <- message
}

func (ws *socket) Close() error {
	return ws.ws.Close()
}

func (ws *socket) writerProcess() {
	defer func() {
		ws.chanClose <- true
	}()

	// Infinite loop
	for {
		select {
		case message, ok := <-ws.chanSend:
			// c.conn.SetWriteDeadline(time.Now().Add(writeWait))

			if !ok {
				ws.ws.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := ws.ws.NextWriter(websocket.TextMessage)
			if err != nil {
				ws.factory.app.Logger.Error("%v", err)
				return
			}

			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(ws.chanSend)
			for i := 0; i < n; i++ {
				w.Write(cNewline)
				w.Write(<-ws.chanSend)
			}

			err = w.Close()
			if err != nil {
				ws.factory.app.Logger.Error("%v", err)
				return
			}

		case <-ws.chanClose:
			return
		}
	}
}

func (ws *socket) readerProcess() {
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
				ws.factory.app.Logger.Error("Error while reading socket: %v", err)
			}
			break
		}

		message = bytes.TrimSpace(bytes.Replace(message, cNewline, cSpace, -1))

		if ws.options.OnMessage != nil {
			ws.options.OnMessage(ws, message)
		}
	}

	if ws.options.OnClose != nil {
		ws.options.OnClose(ws)
	}
}
