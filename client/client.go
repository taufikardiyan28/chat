package client

import (
	"fmt"
	"strings"
	"time"

	"github.com/taufikardiyan28/chat/interfaces"

	"github.com/gorilla/websocket"
	h "github.com/taufikardiyan28/chat/helper"
	MessageModel "github.com/taufikardiyan28/chat/model/messages"
	UserModel "github.com/taufikardiyan28/chat/model/user"
	//room "github.com/taufikardiyan28/chat/room"
)

type (
	Connection struct {
		*websocket.Conn
		UserModel.User
		OnlineUsers     *map[string]*Connection
		messageChannel  chan interface{}
		lastSeenChannel chan UserModel.User
		DB              interfaces.Database
	}
)

func (c *Connection) GetID() string {
	return c.ID
}

func (c *Connection) GetmessageChannel() chan interface{} {
	return c.messageChannel
}

// start listen for client incoming messages
func (c *Connection) Start() {
	defer func() {
		/*if r := recover(); r != nil {
			log.Println("Error", fmt.Sprintf("%v", r))
		}*/
		c.Close()
		delete(*c.OnlineUsers, c.ID)
	}()

	c.messageChannel = make(chan interface{})

	go c.handleClientMessage()

	// get all pending message
	go c.onGetPendingMessage()

	/*msg := MessageModel.MessagePayload{
		DestinationId:   "085246497498",
		DestinationType: "user",
		Msg:             map[string]interface{}{"status": "tes"},
	}
	c.sendToPrivate(msg)

	msg = MessageModel.MessagePayload{
		DestinationId:   "085246497497",
		DestinationType: "user",
		Msg:             map[string]interface{}{"status": "halloooooo"},
	}

	c.sendToPrivate(msg)
	*/

	for {
		msgPayload := MessageModel.MessagePayload{}

		err := c.ReadJSON(&msgPayload)
		if err != nil {
			if strings.Contains(err.Error(), "closed network connection") || strings.Contains(err.Error(), "websocket: close") {
				return
			}
			fmt.Println("Invalid Message ", err)
			continue
		}

		//Send to destination if destination type is not a command
		if msgPayload.MessageType == "chat" {
			if msgPayload.MessageType == "chat" && msgPayload.Msg["status"] != "pending" {
				//msgPayload.OwnerId = c.ID
				//msgPayload.OwnerType = "user"
				// insert to destination id message
				//go c.handleInsertMessage(msgPayload)
			}
			c.Send(msgPayload)
		} else {
			c.onCMD(msgPayload)
		}
	}
}

/******
 ##CHAT GOROUTINE
******/
// handle for received message
func (c *Connection) handleClientMessage() {
	for {
		msg := <-c.messageChannel

		err := c.WriteJSON(msg)
		if err != nil {
			if strings.Contains(err.Error(), "closed network connection") || strings.Contains(err.Error(), "websocket: close") {
				c.Close()
				delete(*c.OnlineUsers, c.ID)
				return
			}
			fmt.Println(err)
			continue
		}

		//fmt.Printf("message receivd chatId: %s, ChatId Length: %d, %s: %f\n", msg.ChatId, len(msg.ChatId), msg.SenderId, msg.Msg["count"].(float64))
	}
}

func (c *Connection) handleInsertMessage(msg MessageModel.MessagePayload) {
	fmt.Println("OwnerID", c.ID)
	msg.OwnerId = c.ID
	err := c.DB.InsertMessage(msg)
	if err != nil {
		fmt.Println("ERROR INSERT MESSAGE ", err)
	}
}

/******##END CHAT GOROUTINE******/

/******
 ##CHAT EVENTS
******/
func (c *Connection) onMessageDelivered(msg MessageModel.MessagePayload) {
	res, err := c.DB.GetMessage(msg)
	if err != nil {
		fmt.Println("ERROR GET MESSAGE ", err)
		return
	}
	for _, elMsg := range res {
		//if el.OwnerId == c.ID {
		elMsg.Msg["delivered_time"] = time.Now().Unix()
		elMsg.Msg["status"] = "delivered"
		err := c.DB.UpdateMessage(elMsg) // need error handling
		if err != nil {
			fmt.Println("ERROR UPDATE DELIVERED ", err)
			continue
		}

		elMsg.MessageType = "info-delivered"
		dstClient, exists := (*c.OnlineUsers)[elMsg.SenderId]
		resp := []MessageModel.MessagePayload{elMsg}
		if exists {
			dstClient.GetmessageChannel() <- resp
		}
		//}
	}
}

func (c *Connection) onMessageReaded(msg MessageModel.MessagePayload) {
	res, err := c.DB.GetMessage(msg)
	if err != nil {
		fmt.Println("ERROR GET MESSAGE ", err)
		return
	}
	for _, elMsg := range res {
		//if el.OwnerId == el.SenderId {
		elMsg.Msg["readed_time"] = time.Now().Unix()
		elMsg.Msg["status"] = "readed"
		err := c.DB.UpdateMessage(elMsg) // need error handling
		if err != nil {
			fmt.Println("ERROR UPDATE READED ", err)
			continue
		}
		//}
		elMsg.MessageType = "info-readed"
		dstClient, exists := (*c.OnlineUsers)[elMsg.SenderId]
		resp := []MessageModel.MessagePayload{elMsg}
		if exists {
			dstClient.GetmessageChannel() <- resp
		}
	}
}

func (c *Connection) onGetPendingMessage() {
	res, err := c.DB.GetPendingMessage(c.ID)

	if err == nil {
		for i, elMsg := range res {
			fmt.Println(elMsg)
			res[i].MessageType = "chat"
		}

		if len(res) > 0 {
			c.GetmessageChannel() <- res
		}
	} else {
		fmt.Println("err get pending message", err)
	}
}

func (c *Connection) onGetHistory(msg MessageModel.MessagePayload) {
	limit, limitValid := msg.Msg["limit"].(float64)
	offset, offsetValid := msg.Msg["offset"].(float64)

	if !limitValid || !offsetValid {
		msg := h.GenerateErrorResponse(msg.OwnerId, "user", "invalid limit or offset value")
		resp := []MessageModel.MessagePayload{msg}
		c.GetmessageChannel() <- resp
		return
	}

	res, err := c.DB.GetChatHistory(c.ID, msg.DestinationId, int(limit), int(offset))
	if err != nil {
		msg := h.GenerateErrorResponse(msg.OwnerId, "user", "Error get chat history")
		resp := []MessageModel.MessagePayload{msg}
		c.GetmessageChannel() <- resp
		return
	}

	for i, _ := range res {
		res[i].MessageType = "chat-history"
	}
	c.GetmessageChannel() <- res
}

/******##END CHAT EVENTS******/

/******
 ##START SOCKET COMMAND
******/
func (c *Connection) Send(msg MessageModel.MessagePayload) {
	switch msg.DestinationType {
	case "room":
		// set message ownerId to groupId
		msg.OwnerId = msg.DestinationId
		msg.OwnerType = "group"
		c.sendToRoom(msg.DestinationId, msg)
		break
	default:
		// set message ownerId to senderId
		msg.OwnerId = c.ID
		msg.OwnerType = "user"
		msg.SenderId = c.ID
		msg.Msg["status"] = "pending"
		msg.Msg["server_time"] = time.Now().Unix()
		msg.Msg["sender_name"] = c.Name
		// insert to owner message
		go c.handleInsertMessage(msg)

		c.sendToPrivate(msg)
	}
}

func (c *Connection) sendToPrivate(msg MessageModel.MessagePayload) {
	dstClient, exists := (*c.OnlineUsers)[msg.DestinationId]
	if !exists {
		// send push notif
		//c.GetmessageChannel() <- h.GenerateErrorResponse(c.ID, "private", fmt.Sprintf("User %s is offline", msg.DestinationId))
	} else {
		go dstClient.handleInsertMessage(msg)

		// make private to array
		resp := []MessageModel.MessagePayload{msg}
		dstClient.GetmessageChannel() <- resp
	}
}

func (c *Connection) sendToRoom(room_id string, msg MessageModel.MessagePayload) {
	c.WriteJSON(msg)
	/*r, exists := room.GetRoom(room_id)
	if !exists {
		c.messageChannel <- message.GenerateErrorResponse(c.ID, "private", "Room ID not found")
	} else {
		respMsg := c.createResponse(msg)
		r.Broadcast(c.ID, respMsg)
	}*/
}

func (c *Connection) onCMD(msg MessageModel.MessagePayload) {
	switch msg.MessageType {
	case "info-delivered":
		go c.onMessageDelivered(msg)
		break
	case "info-readed":
		go c.onMessageReaded(msg)
		break
	case "chat-history":
		go c.onGetHistory(msg)
	}
}

func (c *Connection) createRoom(roomName string) {
	/*
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

		msgPayload := message.MessagePayload{
			Sender:       sender_info,
			ReceiverId:   new_room_id,
			ReceiverType: "room",
			Time:         time.Now(),
			MessageType:  "info",
			ContentType:  "text",
			Content:      fmt.Sprintf("Room Created id : %s", new_room_id),
		}

		c.messageChannel <- msgPayload*/
}

func (c *Connection) joinRoom(roomId string) {
	/*
		r, exists := room.GetRoom(roomId)
		if !exists {
			c.messageChannel <- message.GenerateErrorResponse(c.ID, "private", "Room ID not found")
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

			msgPayload := message.MessagePayload{
				Sender:       sender_info,
				ReceiverId:   r.ID,
				ReceiverType: "room",
				Time:         time.Now(),
				MessageType:  "info",
				ContentType:  "text",
				Content:      fmt.Sprintf("User %s Joined the room", c.UserName),
			}

			r.Broadcast("", msgPayload)
		}
	*/
}

/******##END CHAT COMMAND******/
