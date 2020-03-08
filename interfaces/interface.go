package interfaces

import (
	MessageModel "github.com/taufikardiyan28/chat/model/messages"
	UserModel "github.com/taufikardiyan28/chat/model/user"
)

type (
	IClient interface {
		GetID() string
		Start()
		Send(msg MessageModel.MessagePayload)
		GetPrivateChannel() chan MessageModel.MessagePayload
	}

	IDatabase interface {
		Ping() error
		Connect() error
		Exec(query string, args ...interface{}) (interface{}, error)
		Select(dest interface{}, query string, args ...interface{}) error
		Get(dest interface{}, query string, args ...interface{}) error
	}

	IUserRepo interface {
		GetUserInfo(id string) (UserModel.User, error)
		UpdateUser(id string, cols []string, val ...interface{}) error
	}

	IMessageRepo interface {
		GetChatList(userId string, limit, offset int) ([]MessageModel.Chat, error)
		GetChatHistory(ownerId string, destId string, limit, offset int) ([]MessageModel.MessagePayload, error)
		GetPendingMessage(ownerId string) ([]MessageModel.MessagePayload, error)
		GetMessage(msg MessageModel.MessagePayload) ([]MessageModel.MessagePayload, error)
		InsertMessage(msg MessageModel.MessagePayload) error
		UpdateMessage(msg MessageModel.MessagePayload) error
	}
)
