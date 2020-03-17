package client

import (
	"fmt"
	"sync"
	"time"

	"github.com/taufikardiyan28/chat/helper"
	"github.com/taufikardiyan28/chat/interfaces"

	"github.com/gorilla/websocket"
	h "github.com/taufikardiyan28/chat/helper"
	MessageModel "github.com/taufikardiyan28/chat/model/messages"
	UserModel "github.com/taufikardiyan28/chat/model/user"
	//room "github.com/taufikardiyan28/chat/room"
)

type (
	Connection struct {
		UUID   string
		Config *helper.Configuration
		*websocket.Conn
		UserModel.User
		OnlineUsers     *sync.Map //*map[string]*Connection
		messageChannel  chan interface{}
		lastSeenChannel chan UserModel.User
		Pool            interfaces.IDatabase
		UserRepo        interfaces.IUserRepo
		MessageRepo     interfaces.IMessageRepo
	}
)

func (c *Connection) GetID() string {
	return c.ID
}

func (c *Connection) GetmessageChannel() (chan interface{}, bool) {
	return c.messageChannel, c.IsChanClosed(c.messageChannel)
}

// start listen for client incoming messages
func (c *Connection) Start() {
	defer func() {
		c.Close()
		//c.OnlineUsers.Delete(c.ID)
		c.RemoveConnection()
		if !c.IsChanClosed(c.messageChannel) {
			close(c.messageChannel)
		}
		fmt.Println("Closing", c.ID)
	}()

	c.messageChannel = make(chan interface{})

	go c.handleClientMessage()

	// get all pending message
	go c.onGetPendingMessage()

	//update online status
	go c.OnUserOnline()

	//go c.Ping()
	//c.SetReadDeadline(time.Now().Add(time.Second * 5))
	for {
		msgPayload := MessageModel.MessagePayload{}

		err := c.ReadJSON(&msgPayload)
		if err != nil {
			c.OnUserOffline()
			fmt.Println("Invalid Message ", err)
			break
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

func (c *Connection) RemoveConnection() {
	//fmt.Println("removing")
	iClients, exists := (*c.OnlineUsers).Load(c.ID)
	if exists {
		connections := iClients.([]*Connection)
		conIdx := 0
		for i, con := range connections {
			if c.UUID == con.UUID {
				conIdx = i
			}
		}
		connections = append(connections[:conIdx], connections[conIdx+1:]...)
		if len(connections) == 0 {
			(*c.OnlineUsers).Delete(c.ID)
		} else {
			(*c.OnlineUsers).Store(c.ID, connections)
		}
	}
}

/******
 ##CHAT GOROUTINE
******/
// handle for received message
func (c *Connection) Ping() {
	msg := []MessageModel.MessagePayload{MessageModel.MessagePayload{
		MessageType: "ping",
	}}
	for {
		time.Sleep(10 * time.Second)
		c.SendMessage(msg)
	}
}

func (c *Connection) IsChanClosed(ch <-chan interface{}) bool {
	select {
	case <-ch:
		return true
	default:
	}

	return false
}

func (c *Connection) handleClientMessage() {
	defer func() {
		c.Close()
		//c.OnlineUsers.Delete(c.ID)
		c.RemoveConnection()
		if !c.IsChanClosed(c.messageChannel) {
			close(c.messageChannel)
		}
		//fmt.Println("Closing from Write", c.ID)
	}()
	for {
		msg, ok := <-c.messageChannel
		if !ok {
			break
		}

		err := c.WriteJSON(msg)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func (c *Connection) handleInsertMessage(msg MessageModel.MessagePayload) {
	//msg.OwnerId = c.ID
	if msg.SenderId == msg.OwnerId {
		msg.InterlocutorsId = msg.DestinationId
	} else {
		msg.InterlocutorsId = msg.SenderId
	}

	err := c.MessageRepo.InsertMessage(msg)
	if err != nil {
		fmt.Println("ERROR INSERT MESSAGE ", err)
	}
}

/******##END CHAT GOROUTINE******/

/******
 ##CHAT EVENTS
******/
func (c *Connection) onMessageDelivered(msg MessageModel.MessagePayload) {
	res, err := c.MessageRepo.GetMessage(msg)
	if err != nil {
		fmt.Println("ERROR GET MESSAGE ", err)
		return
	}
	for _, elMsg := range res {
		//if el.OwnerId == c.ID {
		elMsg.Msg["delivered_time"] = time.Now().Unix()
		elMsg.Msg["status"] = "delivered"
		err := c.MessageRepo.UpdateMessage(elMsg) // need error handling
		if err != nil {
			fmt.Println("ERROR UPDATE DELIVERED ", err)
			continue
		}

		elMsg.MessageType = "info-delivered"
		iDstClient, exists := (*c.OnlineUsers).Load(elMsg.SenderId) //(*c.OnlineUsers)[elMsg.SenderId]
		resp := []MessageModel.MessagePayload{elMsg}
		if exists {
			dstClient := iDstClient.([]*Connection)
			c.BroadcastToClientId(dstClient, resp)
		}
		//}
	}
}

func (c *Connection) onMessageReaded(msg MessageModel.MessagePayload) {
	res, err := c.MessageRepo.GetMessage(msg)
	if err != nil {
		fmt.Println("ERROR GET MESSAGE ", err)
		return
	}
	for _, elMsg := range res {
		//if el.OwnerId == el.SenderId {
		elMsg.Msg["readed_time"] = time.Now().Unix()
		elMsg.Msg["status"] = "readed"
		err := c.MessageRepo.UpdateMessage(elMsg) // need error handling
		if err != nil {
			fmt.Println("ERROR UPDATE READED ", err)
			continue
		}
		//}
		elMsg.MessageType = "info-readed"
		iDstClient, exists := (*c.OnlineUsers).Load(elMsg.SenderId) //(*c.OnlineUsers)[elMsg.SenderId]
		resp := []MessageModel.MessagePayload{elMsg}
		if exists {
			dstClient := iDstClient.([]*Connection)
			c.BroadcastToClientId(dstClient, resp)
		}
	}
}

func (c *Connection) onGetPendingMessage() {
	res, err := c.MessageRepo.GetPendingMessage(c.ID)

	if err == nil {
		for i, _ := range res {
			// send push notification
			//go c.SendPushNotification(pendingMsg)
			res[i].MessageType = "chat"
		}

		if len(res) > 0 {
			c.SendMessage(res)
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
		c.SendMessage(resp)
		return
	}

	res, err := c.MessageRepo.GetChatHistory(c.ID, msg.DestinationId, int(limit), int(offset))
	if err != nil {
		msg := h.GenerateErrorResponse(msg.OwnerId, "user", "Error get chat history")
		resp := []MessageModel.MessagePayload{msg}
		c.SendMessage(resp)
		return
	}

	for i, _ := range res {
		res[i].MessageType = "chat-history"
	}
	c.SendMessage(res)
}

func (c *Connection) onGetChatList(msg MessageModel.MessagePayload) {
	limit, limitValid := msg.Msg["limit"].(float64)
	offset, offsetValid := msg.Msg["offset"].(float64)

	if !limitValid || !offsetValid {
		msg := h.GenerateErrorResponse(msg.OwnerId, "user", "invalid limit or offset value")
		resp := []MessageModel.MessagePayload{msg}
		c.SendMessage(resp)
		return
	}

	res, err := c.MessageRepo.GetChatList(c.ID, int(limit), int(offset))
	if err != nil {
		msg := h.GenerateErrorResponse(msg.OwnerId, "user", "Error get chat history")
		fmt.Println("ERROR GET CHAT LIST", err)
		resp := []MessageModel.MessagePayload{msg}
		c.SendMessage(resp)
		return
	}

	for i, _ := range res {
		res[i].MessageType = "chat-list"
	}

	c.SendMessage(res)
}

func (c *Connection) OnUserOnline() {
	var vals []interface{}
	vals = append(vals, time.Now())
	vals = append(vals, "online")
	var cols []string
	cols = append(cols, "lastSeen")
	cols = append(cols, "status")
	c.UserRepo.UpdateUser(c.ID, cols, vals...)
}

func (c *Connection) OnUserOffline() {
	var vals []interface{}
	vals = append(vals, time.Now())
	vals = append(vals, "offline")
	var cols []string
	cols = append(cols, "lastSeen")
	cols = append(cols, "status")
	c.UserRepo.UpdateUser(c.ID, cols, vals...)
}

func (c *Connection) onGetUserStatus(msg MessageModel.MessagePayload) {
	res, err := c.UserRepo.GetUserInfo(msg.DestinationId)
	if err != nil {
		fmt.Println("ERROR GET USER STATUS ", err)
		return
	}
	msg.Msg = make(map[string]interface{})
	msg.Msg["id"] = msg.DestinationId
	msg.Msg["name"] = res.Name
	msg.Msg["last_seen"] = res.LastSeen
	msg.Msg["status"] = res.Status

	resp := []MessageModel.MessagePayload{msg}

	c.SendMessage(resp)
}

/******##END CHAT EVENTS******/

/******
 ##START SOCKET COMMAND
******/

// send message to all connected socket with same id
func (c *Connection) BroadcastToClientId(clients []*Connection, msg interface{}) {
	for _, cl := range clients {
		cl.SendMessage(msg)
	}
}

func (c *Connection) SendMessage(msg interface{}) {
	ch, isClosed := c.GetmessageChannel()
	if !isClosed {
		ch <- msg
	}
}

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
		msg.OwnerId = c.ID
		go c.handleInsertMessage(msg)

		c.sendToPrivate(msg)
	}
}

func (c *Connection) sendToPrivate(msg MessageModel.MessagePayload) {
	iDstClient, exists := (*c.OnlineUsers).Load(msg.DestinationId) //(*c.OnlineUsers)[msg.DestinationId]
	if !exists {
		// send push notif
		go c.SendPushNotification(msg)

		msg.OwnerId = msg.DestinationId
		go c.handleInsertMessage(msg)
	} else {
		dstClient := iDstClient.([]*Connection)
		resp := []MessageModel.MessagePayload{msg}
		for _, cl := range dstClient {
			msg.OwnerId = cl.ID
			cl.SendMessage(resp)
		}

		go c.handleInsertMessage(msg)
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
	case "chat-list":
		go c.onGetChatList(msg)
	case "user-status":
		go c.onGetUserStatus(msg)
		break
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
