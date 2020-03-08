package room

import (
	"time"

	interfaces "github.com/taufikardiyan28/chat/interfaces"
	MessageModel "github.com/taufikardiyan28/chat/model/messages"
)

type (
	ChatRoom struct {
		ID        string //UUIDv4
		Name      string
		CreatedBy string
		CreatedAt time.Time
		Members   RoomMember
		Clients   map[string]interfaces.IClient
	}
	RoomMember struct {
		ID   string
		Role string
	}
)

var Rooms *map[string]*ChatRoom

func NewRoom(room_id string, name string) *ChatRoom {
	chatRoom := &ChatRoom{
		ID:        room_id,
		Name:      name,
		CreatedAt: time.Now(),
		Clients:   make(map[string]interfaces.IClient),
	}

	(*Rooms)[room_id] = chatRoom

	return chatRoom
}

func GetRoom(room_id string) (*ChatRoom, bool) {
	r, exists := (*Rooms)[room_id]

	return r, exists
}

func GetRoomByName(name string) *ChatRoom {
	for _, r := range *Rooms {
		if r.Name == name {
			return r
		}
	}
	return nil
}

func (c *ChatRoom) Join(client_id string, client interfaces.IClient) {
	c.Clients[client_id] = client
}

func (c *ChatRoom) Broadcast(sender_id string, msg MessageModel.MessagePayload) {
	for _, client := range c.Clients {
		if client.GetID() != sender_id {
			client.GetPrivateChannel() <- msg
		}
	}
}
