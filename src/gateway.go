package src

import (
	"github.com/gorilla/websocket"
)

// Gateway is what is used to connect to discord
type Gateway struct {
	wsConnection *websocket.Conn
}
