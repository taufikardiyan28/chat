package message

import (
	"time"
)

type (
	From           map[string]interface{}
	MessagePayload struct {
		Dst     string  `json:"dst"`      //public, room_id, username
		DstType string  `json:"dst_type"` //private, room
		Message Message `json:"message"`
	}
	Message struct {
		ID          string    `json:"id"`         // use UUIDv4
		RequestId   string    `json:"request_id"` //unique request id from client
		Time        time.Time `json:"time"`
		ContentType string    `json:"content_type"` //text/image/audio/video
		Content     string    `json:"content"`
	}

	ResponseMessage struct {
		Status int         `json:"status"`
		From   interface{} `json:"from"`
		MessagePayload
		Error string `json:"error"`
	}
)

func GenerateErrorResponse(dstId string, text string) ResponseMessage {
	sender_info := From{
		"id":       "system",
		"username": "System",
		"nickname": "System",
		"phone":    "System",
	}

	msg := Message{
		Time:        time.Now(),
		ContentType: "text",
		Content:     text,
	}

	msgPayload := MessagePayload{
		Dst:     dstId,
		DstType: "private",
		Message: msg,
	}
	resp := ResponseMessage{
		Status:         1,
		From:           sender_info,
		MessagePayload: msgPayload,
		Error:          text,
	}

	return resp
}
