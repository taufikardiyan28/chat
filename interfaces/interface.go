package interfaces

import (
	"github.com/taufikardiyan28/chat/helper"
	MessageModel "github.com/taufikardiyan28/chat/model/messages"
	UserModel "github.com/taufikardiyan28/chat/model/user"
)

type (
	Client interface {
		GetID() string
		Start()
		Send(msg MessageModel.MessagePayload)
		GetPrivateChannel() chan MessageModel.MessagePayload
	}

	Database interface {
		Ping() error
		GetUserInfo(id string) (UserModel.User, error)
		UpdateUser(id string) error
		Connect() error

		GetChatList(userId string, paging helper.PagingData) ([]MessageModel.Chat, error)
		GetChatHistory(ownerId string, destId string, limit, offset int) ([]MessageModel.MessagePayload, error)
		GetPendingMessage(ownerId string) ([]MessageModel.MessagePayload, error)
		GetMessage(msg MessageModel.MessagePayload) ([]MessageModel.MessagePayload, error)
		InsertMessage(msg MessageModel.MessagePayload) error
		UpdateMessage(msg MessageModel.MessagePayload) error
	}
)
