package client

import (
	"fmt"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	interfaces "github.com/taufikardiyan28/chat/interfaces"
	message "github.com/taufikardiyan28/chat/message"
	room "github.com/taufikardiyan28/chat/room"
)

type (
	Connection struct {
		*websocket.Conn
		UserInfo
		Public         *map[string]*Connection
		privateChannel chan message.ResponseMessage
	}

	UserInfo struct {
		ID       string `json:"id"`
		UserId   int64  `json:"-"`
		UserName string `json:"username"`
		NickName string `json:"nickname"`
		Phone    string `json:"phone"`
		Email    string `json:"-"`
	}
)

func (c *Connection) GetID() string {
	return c.ID
}

func (c *Connection) GetPrivateChannel() chan message.ResponseMessage {
	return c.privateChannel
}

// start listen for client incoming messages
func (c *Connection) Listen() {
	defer func() {
		/*if r := recover(); r != nil {
			log.Println("Error", fmt.Sprintf("%v", r))
		}*/
		c.Close()
		delete(*c.Public, c.ID)
	}()

	c.privateChannel = make(chan message.ResponseMessage)

	go c.handleClientMessage()

	for {
		msgPayload := message.MessagePayload{}

		err := c.ReadJSON(&msgPayload)
		if err != nil {
			if strings.Contains(err.Error(), "websocket: close") {
				return
			}
			fmt.Println(err)
			continue
		}
		//fmt.Println(msgPayload)

		//Send to destination if destination type is not a command
		if msgPayload.DstType != "cmd" {
			c.Send(msgPayload)
		} else {
			c.onCMD(msgPayload.Dst, msgPayload.Message.Content)
		}
	}
}

func (c *Connection) handleClientMessage() {
	for {
		msg := <-c.privateChannel
		err := c.WriteJSON(msg)
		if err != nil {
			if strings.Contains(err.Error(), "websocket: close") {
				c.Close()
				delete(*c.Public, c.ID)
				return
			}
			fmt.Println(err)
			continue
		}
	}
}

func (c *Connection) Send(msg message.MessagePayload) {
	msg.Message.ID = uuid.NewV4().String()

	switch msg.DstType {
	case "room":
		c.sendToRoom(msg.Dst, msg)
		break
	default:
		c.sendToPrivate(msg.Dst, msg)
	}
}

func (c *Connection) createResponse(msg message.MessagePayload) message.ResponseMessage {
	c_info := UserInfo{
		ID:       c.ID,
		UserName: c.UserName,
		NickName: c.NickName,
		Phone:    c.Phone,
	}

	// set to current time
	msg.Message.Time = time.Now()

	respMsg := message.ResponseMessage{
		Status:         0,
		From:           c_info,
		MessagePayload: msg,
	}
	return respMsg
}

func (c *Connection) sendToPrivate(to string, msg message.MessagePayload) {
	dstClient, exists := (*c.Public)[to]
	if !exists {
		c.privateChannel <- message.GenerateErrorResponse(c.ID, fmt.Sprintf("User \"%s\" not found", to))
	} else {
		resp := c.createResponse(msg)
		dstClient.privateChannel <- resp
	}
}

func (c *Connection) sendToRoom(room_id string, msg message.MessagePayload) {
	c.WriteJSON(msg)
	r, exists := room.GetRoom(room_id)
	if !exists {
		c.privateChannel <- message.GenerateErrorResponse(c.ID, "Room ID not found")
	} else {
		respMsg := c.createResponse(msg)
		r.Broadcast(c.ID, respMsg)
	}
}

func (c *Connection) onCMD(cmd string, cmdValue string) {
	switch cmd {
	case "create-room":
		c.createRoom(cmdValue)
		break
	case "join-room":
		c.joinRoom(cmdValue)
		break
	}
}

func (c *Connection) createRoom(roomName string) {
	new_room_id := uuid.NewV4().String()
	r := room.NewRoom(new_room_id, roomName)
	var IClient interfaces.Client
	IClient = c
	r.Join(c.ID, IClient)
	sender_info := UserInfo{
		ID:       new_room_id,
		UserName: roomName,
		NickName: roomName,
		Phone:    roomName,
	}

	msg := message.Message{
		Time:        time.Now(),
		ContentType: "text",
		Content:     fmt.Sprintf("Room Created id : %s", new_room_id),
	}

	msgPayload := message.MessagePayload{
		Dst:     new_room_id,
		DstType: "room",
		Message: msg,
	}
	respMsg := message.ResponseMessage{
		Status:         0,
		From:           sender_info,
		MessagePayload: msgPayload,
		Error:          "",
	}
	c.privateChannel <- respMsg
}

func (c *Connection) joinRoom(roomId string) {
	r, exists := room.GetRoom(roomId)
	if !exists {
		c.privateChannel <- c.generateErrorResponse("Room ID not found")
	} else {
		var IClient interfaces.Client
		IClient = c
		r.Join(c.ID, IClient)

		sender_info := UserInfo{
			ID:       c.ID,
			UserName: c.UserName,
			NickName: c.NickName,
			Phone:    c.Phone,
		}

		msg := message.Message{
			Time:        time.Now(),
			ContentType: "text",
			Content:     fmt.Sprintf("User %s Joined the room", c.UserName),
		}

		msgPayload := message.MessagePayload{
			Dst:     r.ID,
			DstType: "room",
			Message: msg,
		}
		respMsg := message.ResponseMessage{
			Status:         0,
			From:           sender_info,
			MessagePayload: msgPayload,
			Error:          "",
		}

		r.Broadcast("", respMsg)
	}
}

func (c *Connection) createRoom(roomName string) {
	new_room_id := uuid.NewV4().String()
	r := room.NewRoom(new_room_id, roomName)
	var IClient interfaces.Client
	IClient = c
	r.Join(c.ID, IClient)
	sender_info := UserInfo{
		ID:       new_room_id,
		UserName: roomName,
		NickName: roomName,
		Phone:    roomName,
	}

	msg := message.Message{
		Time:        time.Now(),
		ContentType: "text",
		Content:     fmt.Sprintf("Room Created id : %s", new_room_id),
	}

	msgPayload := message.MessagePayload{
		Dst:     new_room_id,
		DstType: "room",
		Message: msg,
	}
	respMsg := message.ResponseMessage{
		Status:         0,
		From:           sender_info,
		MessagePayload: msgPayload,
		Error:          "",
	}
	c.privateChannel <- respMsg
}

func (c *Connection) joinRoom(roomId string) {
	r, exists := room.GetRoom(roomId)
	if !exists {
		c.privateChannel <- message.GenerateErrorResponse(c.ID, "Room ID not found")
	} else {
		var IClient interfaces.Client
		IClient = c
		r.Join(c.ID, IClient)

		sender_info := UserInfo{
			ID:       c.ID,
			UserName: c.UserName,
			NickName: c.NickName,
			Phone:    c.Phone,
		}

		msg := message.Message{
			Time:        time.Now(),
			ContentType: "text",
			Content:     fmt.Sprintf("User %s Joined the room", c.UserName),
		}

		msgPayload := message.MessagePayload{
			Dst:     r.ID,
			DstType: "room",
			Message: msg,
		}
		respMsg := message.ResponseMessage{
			Status:         0,
			From:           sender_info,
			MessagePayload: msgPayload,
			Error:          "",
		}

		r.Broadcast("", respMsg)
	}
}
