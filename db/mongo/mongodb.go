package MongoDB

import (
	"context"
	"fmt"
	"time"

	"github.com/taufikardiyan28/chat/helper"
	h "github.com/taufikardiyan28/chat/helper"
	MessageModel "github.com/taufikardiyan28/chat/model/messages"
	UserModel "github.com/taufikardiyan28/chat/model/user"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	Conn struct {
		Config *helper.Configuration
		Pool   *mongo.Client
	}
)

// NOT IMPLEMENTED YET

func (c *Conn) Connect() error {
	uri := fmt.Sprintf("mongodb://%s:%d", c.Config.Database.Host, c.Config.Database.Port)
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		fmt.Println(err)
		return err
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		return err
	}
	c.Pool = client
	//mongoClient = client
	return err
}

func (c *Conn) Ping() error {
	return nil //return c.Pool.Ping()
}

func (c *Conn) GetUserInfo(id string) (UserModel.User, error) {
	return UserModel.User{}, nil
}

func (c *Conn) UpdateUser(id string) error {
	return nil
}

func (c *Conn) GetChatList(id string, paging h.PagingData) ([]MessageModel.Chat, error) {
	return nil, nil
}

func (c *Conn) GetChatHistory(ownerId string, destId string, limit, offset int) ([]MessageModel.MessagePayload, error) {
	return nil, nil
}

func (c *Conn) GetPendingMessage(ownerId string) ([]MessageModel.MessagePayload, error) {
	var res []MessageModel.MessagePayload
	return res, nil
}

func (c *Conn) GetMessage(msg MessageModel.MessagePayload) ([]MessageModel.MessagePayload, error) {
	res := []MessageModel.MessagePayload{}
	return res, nil
}

func (c *Conn) InsertMessage(msg MessageModel.MessagePayload) error {
	return nil
}

func (c *Conn) UpdateMessage(msg MessageModel.MessagePayload) error {
	return nil
}

func (c *Conn) GetPool() interface{} {
	return c.Pool
}
