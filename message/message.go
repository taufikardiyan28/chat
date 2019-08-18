package message

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	From           map[string]interface{}
	MessagePayload struct {
		ID           primitive.ObjectID `json:"id" bson:"_id"`
		OwnerId      string             `json:"-" bson:"owner_id"`
		Sender       interface{}        `json:"sender" bson:"sender"`
		ReceiverId   string             `json:"receiver_id" bson:"receiver_id"`     //public, room_id, username
		ReceiverType string             `json:"receiver_type" bson:"receiver_type"` //private, room
		RequestId    string             `json:"request_id" bson:"request_id"`       //unique request id from client
		Time         time.Time          `json:"time" bson:"time"`
		MessageType  string             `json:"message_type" bson:"message_type"` //chat/info
		ContentType  string             `json:"content_type" bson:"content_type"` //text/image/audio/video
		Content      string             `json:"content" bson:"content"`
		Status       MessageStatus      `json:"status" bson:"status"`
	}

	MessageStatus struct {
		Received     bool      `json:"received" bson:"received"`
		ReceivedTime time.Time `json:"received_time" bson:"received_time"`
		Readed       bool      `json:"readed" bson:"readed"`
		ReadedTime   time.Time `json:"readed_time" bson:"readed_time"`
	}
)

func GenerateErrorResponse(receiverId, receiverType, text string) MessagePayload {
	sender_info := From{
		"id":       "system",
		"username": "System",
		"nickname": "System",
		"phone":    "System",
	}

	msgPayload := MessagePayload{
		Sender:       sender_info,
		ReceiverId:   receiverId,
		ReceiverType: receiverType,
		Time:         time.Now(),
		MessageType:  "info",
		ContentType:  "text",
		Content:      text,
	}

	return msgPayload
}
