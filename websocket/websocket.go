package appy_websockets

import (
	"bytes"
	"net/http"

	"github.com/gorilla/websocket"
	appy_logger "github.com/nfwGytautas/appy-go/logger"
)

type Websocket struct {
	factory *WebsocketFactory
	ws      *websocket.Conn

	id uint64

	chanClose chan bool
	chanSend  chan []byte

	options WebsocketOptions
}

// Start the websocket
func (ws *Websocket) Spin(writer http.ResponseWriter, request *http.Request) error {
	var err error
	ws.ws, err = ws.factory.upgrader.Upgrade(writer, request, nil)
	if err != nil {
		appy_logger.Get().Error("Failed to spin websocket '%v'", err)
		return err
	}

	defer ws.Close()

	// Now run read/write process
	go ws.writerProcess()
	go ws.readerProcess()

	<-ws.chanClose
	appy_logger.Get().Debug("Ending socket spin for '%v'", ws.id)

	err = ws.Close()
	if err != nil {
		appy_logger.Get().Error("Failed to close websocket '%v'", err)
	}

	return err
}

// Send a message to websocket
func (ws *Websocket) Send(message []byte) {
	appy_logger.Get().Debug("Sending message '%v' to '%v'", string(message), ws.id)
	ws.chanSend <- message
}

// Close the websocket
func (ws *Websocket) Close() error {
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
			appy_logger.Get().Debug("[%v] sending '%s'", ws.id, string(message))
			// c.conn.SetWriteDeadline(time.Now().Add(writeWait))

			if !ok {
				ws.ws.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := ws.ws.NextWriter(websocket.TextMessage)
			if err != nil {
				appy_logger.Get().Error("%v", err)
				return
			}

			appy_logger.Get().Debug("[%v] writing", ws.id)
			w.Write(message)
			w.Write(cNewline)

			// Add queued chat messages to the current websocket message.
			n := len(ws.chanSend)
			appy_logger.Get().Debug("[%v] queue backlog: '%v'", ws.id, n)
			for i := 0; i < n; i++ {
				w.Write(<-ws.chanSend)
				w.Write(cNewline)
			}

			err = w.Close()
			if err != nil {
				appy_logger.Get().Error("%v", err)
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
				appy_logger.Get().Error("Error while reading socket: %v", err)
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
