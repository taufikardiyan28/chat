package MessageModel

import (
	"database/sql/driver"
	"encoding/json"
)

type (
	MessageBody map[string]interface{}
	Chat        struct {
		ChatId          string      `json:"chat_id" bson:"chat_id" db:"chat_id"`
		DestinationId   string      `json:"destination_id" bson:"destionation_id" db:"destinationId"`
		DestinationType string      `json:"destination_type" bson:"destination_type" db:"destinationType"`
		SenderName      string      `json:"sender_name" bson:"sender_name" db:"senderName"`
		LastMessage     MessageBody `json:"last_message" bson:"last_message" db:"lastMessage"`
	}

	MessagePayload struct {
		ID              int64       `json:"-" db:"id" bson:"id"`
		MessageType     string      `json:"message_type" bson:"message_type"`
		OwnerId         string      `json:"-" bson:"owner_id" db:"ownerId"`           //public, room_id, username
		OwnerType       string      `json:"-" bson:"owner_type" db:"ownerType"`       //group, user, public
		ChatId          string      `json:"chat_id" bson:"chat_id" db:"chatId"`       //unique request id from client
		SenderId        string      `json:"sender_id" bson:"sender_id" db:"senderId"` //private, room
		DestinationId   string      `json:"destination_id" bson:"destination_id" db:"destinationId"`
		DestinationType string      `json:"destination_type" bson:"destination_type" db:"destinationType"` //group, user, public
		Msg             MessageBody `json:"msg" bson:"msg" db:"msg"`
		CreatedAt       string      `json:"created_at" bson:"created_at" db:"createdAt"`
	}
)

/*func (m MessageBody) Value() (driver.Value, error) {
	return json.Marshal(&m)
}*/

func (c *MessageBody) Value() (driver.Value, error) {
	if c != nil {
		b, err := json.Marshal(c)
		if err != nil {
			return nil, err
		}
		return string(b), nil
	}
	return nil, nil
}

func (c *MessageBody) Scan(src interface{}) error {
	var data []byte
	if b, ok := src.([]byte); ok {
		data = b
	} else if s, ok := src.(string); ok {
		data = []byte(s)
	}
	return json.Unmarshal(data, c)
}
