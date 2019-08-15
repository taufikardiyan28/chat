package interfaces

import (
	message "github.com/taufikardiyan28/chat/message"
)

type (
	Client interface {
		GetID() string
		Listen()
		Send(msg message.MessagePayload)
		GetPrivateChannel() chan message.ResponseMessage
	}
)
