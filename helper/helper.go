package helper

import (
	"github.com/google/uuid"
	MessageModel "github.com/taufikardiyan28/chat/model/messages"
)

type (
	Configuration struct {
		Database struct {
			Host     string `yaml:"host"`
			Port     int    `yaml:"port"`
			User     string `yaml:"user"`
			Password string `yaml:"password"`
			DbName   string `yaml:"db_name"`
			DbType   string `yaml:"db_type"`
		} `yaml:"database"`
		Host      string `yaml:"host"`
		Port      int    `yaml:"port"`
		PushNotif struct {
			Url string `yaml:"url"`
		} `yaml:"push_notif"`
	}

	Sorting struct {
		Field string `json:"field" form:"field"`
		Order string `json:"order" form:"order"`
	}
	PagingData struct {
		Page   int32        `json:"page" query:"page" form:"page" validate:"required"`
		Rows   int32        `json:"rows" query:"rows" form:"rows" validate:"required"`
		Filter []FilterRule `json:"filter" query:"filter" form:"filter"`
		Sort   []Sorting    `json:"sort" form:"sort"`
	}

	FilterRule struct {
		Field string `json:"field" form:"field"`
		Op    string `json:"op" form:"op"`
		Value string `json:"value" form:"value"`
	}
)

func GenerateErrorResponse(destinationId, destinationType, text string) MessageModel.MessagePayload {
	msg := MessageModel.MessageBody{
		"content_type": "text",
		"content":      text,
		"sender_name":  "system",
	}

	var chatId string
	chatUUID, err := uuid.NewRandom()
	if err != nil {
		chatId = chatUUID.String()
	}
	msgPayload := MessageModel.MessagePayload{
		OwnerId:         destinationId,
		OwnerType:       destinationType,
		MessageType:     "error",
		SenderId:        "system",
		ChatId:          chatId,
		DestinationId:   destinationId,
		DestinationType: destinationType,
		Msg:             msg,
	}

	return msgPayload
}
