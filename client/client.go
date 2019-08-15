package client

import (
	"fmt"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	message "github.com/taufikardiyan28/chat/message"
)

type (
	Connection struct {
		*websocket.Conn
		UserInfo
		Public         *map[string]*Connection
		PrivateChannel chan message.ResponseMessage
	}

	UserInfo struct {
		Key      string `json:"key"`
		UserId   int64  `json:"-"`
		UserName string `json:"username"`
		NickName string `json:"nickname"`
		Phone    string `json:"phone"`
		Email    string `json:"-"`
	}
)

// start listen for client incoming messages
func (c *Connection) Listen() {
	defer func() {
		/*if r := recover(); r != nil {
			log.Println("Error", fmt.Sprintf("%v", r))
		}*/
		c.Close()
		delete(*c.Public, c.Key)
	}()

	c.PrivateChannel = make(chan message.ResponseMessage)

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
		msg := <-c.PrivateChannel
		err := c.WriteJSON(msg)
		if err != nil {
			if strings.Contains(err.Error(), "websocket: close") {
				c.Close()
				delete(*c.Public, c.Key)
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
	case "private":
		c.sendToPrivate(msg.Dst, msg)
		break
	default:
		c.sendToRoom(msg.Dst, msg)
	}
}

func (c *Connection) sendToPrivate(to string, msg message.MessagePayload) {
	c_info := UserInfo{
		Key:      c.Key,
		UserName: c.UserName,
		NickName: c.NickName,
		Phone:    c.Phone,
	}

	// set to current time
	msg.Message.Time = time.Now()

	resp := message.ResponseMessage{
		Status:         0,
		From:           c_info,
		MessagePayload: msg,
	}

	dstClient, exists := (*c.Public)[to]
	if !exists {
		c.PrivateChannel <- c.generateErrorResponse(fmt.Sprintf("User \"%s\" not found", to))
	} else {
		dstClient.PrivateChannel <- resp
	}

}

func (c *Connection) sendToRoom(group_id string, msg message.MessagePayload) {
	c.WriteJSON(msg)
}

func (c *Connection) onCMD(cmd string, cmdValue string) {

}

func (c *Connection) generateErrorResponse(text string) message.ResponseMessage {
	sender_info := UserInfo{
		Key:      "system",
		UserName: "System",
		NickName: "System",
		Phone:    "System",
	}

	msg := message.Message{
		Time:        time.Now(),
		ContentType: "text",
		Content:     text,
	}

	msgPayload := message.MessagePayload{
		Dst:     c.Key,
		DstType: "private",
		Message: msg,
	}
	resp := message.ResponseMessage{
		Status:         1,
		From:           sender_info,
		MessagePayload: msgPayload,
		Error:          text,
	}

	return resp
}
