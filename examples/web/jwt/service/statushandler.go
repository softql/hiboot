package service

import (
	"fmt"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hiboot/pkg/starter/websocket"
)

// StatusHandler is the websocket handler
type StatusHandler struct {
	at.ContextAware
	connection *websocket.Connection
}

func newStatusHandler(connection *websocket.Connection) *StatusHandler {
	h := &StatusHandler{connection: connection}
	return h
}

func init() {
	app.Register(newStatusHandler)
}

// OnMessage is the websocket message handler
func (h *StatusHandler) OnMessage(data []byte) {
	message := string(data)
	log.Debugf("client: %v", message)

	h.connection.EmitMessage([]byte(fmt.Sprintf("Status: Up")))

}

// OnDisconnect is the websocket disconnection handler
func (h *StatusHandler) OnDisconnect() {
	log.Debugf("Connection with ID: %v has been disconnected!", h.connection.ID())
}
