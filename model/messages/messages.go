package MessagesModel

import (
	"context"
	"fmt"
	"time"

	message "github.com/taufikardiyan28/chat/message"
	db "github.com/taufikardiyan28/chat/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	Message struct {
		*db.Conn
	}
)

func (m *Message) Get() (message.MessagePayload, error) {

	pool := m.GetPool()
	var result message.MessagePayload
	collection := pool.Database("chat").Collection("messages")
	msg := message.MessagePayload{
		ID:           primitive.NewObjectID(),
		Sender:       map[string]interface{}{"user_id": "asdf"},
		ReceiverId:   "11111",
		ReceiverType: "Private",
		Time:         time.Now(),
		MessageType:  "chat",
		ContentType:  "text",
		Content:      "tes message",
		Status: message.MessageStatus{
			Received:     true,
			ReceivedTime: time.Now(),
		},
	}

	insertResult, err := collection.InsertOne(context.Background(), msg)
	if err != nil {
		fmt.Println(err)
		return msg, err
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
	filter := bson.D{{"receiverid", "11111"}}

	err = collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		fmt.Println(err)
		return result, err
	}

	fmt.Printf("Found a single document: %+v\n", result)
	return result, err
}
