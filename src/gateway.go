package src

import (
	"github.com/gorilla/websocket"
)

type Gateway struct {
	wsConnection *websocket.Conn
}
