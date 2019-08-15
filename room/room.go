package room

import (
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

type (
	ChatRoom struct {
		ID      string //UUIDv4
		Name    string
		Clients map[string]*websocket.Conn
	}
)

var Rooms *map[string]*ChatRoom

func NewRoom(name string) *ChatRoom {
	chatRoom := &ChatRoom{
		ID:      uuid.NewV4().String(),
		Name:    name,
		Clients: make(map[string]*websocket.Conn),
	}

	(*Rooms)[name] = chatRoom

	return chatRoom
}
