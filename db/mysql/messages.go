package MySqlDB

import (
	"encoding/json"

	"github.com/taufikardiyan28/chat/helper"
	MessageModel "github.com/taufikardiyan28/chat/model/messages"
)

func (c *Conn) GetChatList(userId string, paging helper.PagingData) ([]MessageModel.Chat, error) {
	return nil, nil
}

func (c *Conn) GetChatHistory(ownerId string, destId string, limit, offset int) ([]MessageModel.MessagePayload, error) {
	strSQL := `SELECT id, ownerId, ownerType, chatId, senderId, destinationId, destinationType, msg, createdAt FROM messages WHERE ownerId=? AND (senderId=? OR destinationId=?) ORDER BY id DESC`
	var res []MessageModel.MessagePayload
	err := c.Pool.Select(&res, strSQL, ownerId, destId, destId)
	return res, err
}

func (c *Conn) GetPendingMessage(ownerId string) ([]MessageModel.MessagePayload, error) {
	strSQL := `SELECT id, ownerId, ownerType, chatId, senderId, destinationId, destinationType, msg, createdAt FROM messages WHERE ownerId=? AND JSON_EXTRACT(msg, "$.status")=?`
	var res []MessageModel.MessagePayload
	err := c.Pool.Select(&res, strSQL, ownerId, "pending")
	return res, err
}

func (c *Conn) GetMessage(msg MessageModel.MessagePayload) ([]MessageModel.MessagePayload, error) {
	strSQL := `SELECT id, ownerId, ownerType, chatId, senderId, destinationId, destinationType, msg, createdAt FROM messages WHERE chatId=?`
	var res []MessageModel.MessagePayload
	err := c.Pool.Select(&res, strSQL, msg.ChatId)
	return res, err
}

func (c *Conn) InsertMessage(msg MessageModel.MessagePayload) error {
	var err error
	strSQL := `INSERT INTO messages (ownerId, ownerType, chatId, senderId, destinationId, destinationType, msg) VALUES (?, ?, ?, ?, ?, ?, ?)`
	msgString, err := json.Marshal(msg.Msg)
	if err != nil {
		return err
	}
	var args = []interface{}{msg.OwnerId, msg.OwnerType, msg.ChatId, msg.SenderId, msg.DestinationId, msg.DestinationType, msgString}
	_, err = c.Pool.Exec(strSQL, args...)
	return err
}

func (c *Conn) UpdateMessage(msg MessageModel.MessagePayload) error {
	var err error
	strSQL := `UPDATE messages SET msg=? WHERE ownerId=? AND chatId=?`
	msgString, err := json.Marshal(msg.Msg)
	if err != nil {
		return err
	}
	var args = []interface{}{msgString, msg.OwnerId, msg.ChatId}
	_, err = c.Pool.Exec(strSQL, args...)
	return err
}
